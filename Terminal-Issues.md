# Terminal Input Issue Fix - Claude Context

## Problem Description
After running the Conway's Game of Life program once, subsequent runs would show control characters literally instead of processing them:
- Enter key would display as `^M` instead of processing newlines
- Ctrl-C would display as `^C` instead of sending interrupt signal
- Terminal remained in raw input mode after program exit

## Root Cause Analysis
The issue was in the terminal state management within the `atomicgo.dev/keyboard` library:

1. **Raw Mode Not Restored**: The keyboard library puts the terminal in raw mode for real-time input capture but sometimes fails to properly restore normal mode on program exit.

2. **Library Vendor Files**: Initial attempts to fix the issue by modifying vendor files were inappropriate since those changes get lost during library updates.

3. **Missing Application-Level Cleanup**: The application needed to ensure terminal state restoration independent of library behavior.

4. **Signal Handling**: The program didn't handle unexpected termination signals properly, which could leave the terminal in raw mode.

## Solution Implemented

### 1. Application-Level Terminal Reset (`ui/legacy.go`)
- Added `resetTerminal()` function that uses `stty sane` command to force terminal reset
- This provides a reliable fallback that works regardless of library behavior
- Uses `os/exec` to run the system command that restores terminal to known good state

### 2. Multiple Exit Path Coverage
- **Deferred cleanup**: Added `defer` function that runs on any program exit
- **Normal quit**: Terminal reset on 'q' key exit path  
- **Signal handlers**: Reset terminal on Ctrl-C, SIGTERM, SIGHUP signals
- **Multiple layers**: Ensures terminal gets reset even if program crashes

### 3. Enhanced Event Listener (`event/listener.go`)
- Added `Stop()` method to `Listener` interface for proper cleanup
- Implemented proper channel closing in `gameListener.Stop()`
- Called from deferred cleanup in UI layer

### 4. Signal Handling Enhancement
- Expanded signal handling to cover SIGINT, SIGTERM, and SIGHUP
- Each signal handler explicitly calls terminal reset before exit
- Ensures cleanup even during unexpected termination

## Files Modified
- `ui/legacy.go` - Added terminal reset function and comprehensive cleanup
- `event/listener.go` - Added Stop method and improved resource management

## Technical Implementation
```go
// resetTerminal forces a terminal reset using stty to restore normal input mode
func resetTerminal() {
    cmd := exec.Command("stty", "sane")
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run() // Ignore errors as this is best-effort cleanup
}
```

The `stty sane` command is a Unix standard that resets terminal settings to sensible defaults, effectively undoing any raw mode or special terminal states.

## Testing
The fix ensures that:
1. Terminal is always restored to normal mode after program exit
2. Control characters work properly in subsequent runs  
3. Program handles interruption signals gracefully
4. Resources are properly cleaned up in all exit scenarios
5. Solution persists through library updates (no vendor file modifications)

## Technical Notes
- The issue was platform-specific to Unix-like systems (macOS, Linux)
- Raw mode is required for real-time key capture but must be properly restored
- Using `stty sane` provides reliable restoration independent of library implementation
- Multiple cleanup layers ensure restoration even if one method fails
- Application-level solution survives library updates and dependency changes
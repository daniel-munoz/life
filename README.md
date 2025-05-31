# Game of Life

A Go implementation of Conway's Game of Life featuring a terminal-based user interface, configuration through sample files, and clear modular code structure.

## Overview

Conway's Game of Life is a cellular automaton simulation where cells evolve based on their neighbors. It's not just a simulation; it provides interesting emergent behaviors.

## Installation

Ensure that Go is installed:

```sh
go version
```

Download dependencies:

```sh
go mod download
```

## Usage

Run the application from the root directory:

```sh
go run main.go [sample_name]
```

If no sample name is provided, the program will present an interactive menu to choose from available samples. Available samples are listed in the `samples/` directory.

### Controls

Once the simulation is running, use the following keys:

- **Arrow Keys**: Move the viewport (Up/Down/Left/Right)
- **I/K/J/L**: Move the viewport by larger increments (10 spaces)
- **Space**: Pause/Resume the simulation
- **H**: Display help
- **Q** or **Ctrl-C**: Quit the program

## Recent Fixes

### Terminal Input Issue (Fixed)
Previous versions had an issue where after running the program once, subsequent runs would display control characters literally (e.g., `^M` for Enter, `^C` for Ctrl-C) instead of processing them normally. This has been resolved through improved terminal state management and cleanup procedures.

For technical details about this fix, see [CLAUDE.md](CLAUDE.md).

## Contribution

Contributions are welcome! Submit your suggestions through GitHub issues or pull requests.

## License

See the [LICENSE.md](LICENSE.md) file for details.

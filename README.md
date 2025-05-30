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

If no sample name is provided, a default one (`gliders`) is used. Available samples are listed in the `samples/` directory.

## Contribution

Contributions are welcome! Submit your suggestions through GitHub issues or pull requests.

## License

See the [LICENSE.md](LICENSE.md) file for details.

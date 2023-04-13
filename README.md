# Spaceship Game

This is an old-fashioned game written in Go. 
A user controls a spaceship flying through a meteor shower.
The work is in progress.

## Getting Started

### Prerequisites

Go v.1.20 should be already installed.

### Installing

Open a project root directory in a console and install project dependencies:
```shell
go mod tidy
```

### Quick Start

Run the game:
```shell
go run cmd/main.go
```

If everything works fine, a console will display a game screen. A user can move his spaceship using left and right arrow keys and shoot using a space key. If the spaceship collides with a meteorite, the game will be over.

![presentation](presentation.gif)

## Development

Currently, there are no tests in the project. You can run a debug mode using special command line arguments.
```shell
go run cmd/main.go -debug true -log_file=spaceship.log
```

## License

This project is licensed under the GNU GPLv3  License - see the [LICENSE.md](LICENSE.md) file for details

# Spaceship Game

This is an old-fashioned game written in Go. 
A user controls a spaceship flying through a meteor shower.
The user has to destroy meteorites and avoid collisions.
The work is in progress.

## Getting Started

### Installation

#### Linux
- Download a zip file `spaceship-linux-{version_number}.zip` from [here](https://github.com/AndreyAD1/spaceship/releases).
- Unpack this archive
- Open your terminal
- Run the game
```shell
$ ./spaceship.exe
```
If everything works fine, a console will display a game screen.
A user can move a spaceship using arrow keys and shoot using a space key.

![presentation](presentation.gif)

## Development

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
go run main.go
```

If everything works fine, a console will display a game screen.

### Tests

Run the project tests: 
```shell
go test ./...
```

### Additional Info

You can collect logs and cpu profile info in a file.
```shell
go run main.go --debug true --log_file=spaceship.log --cpuprofile=cpu.out
```

## Acknowledgements

The inspiration for this project comes from a Python developer course at
[dvmn.org](https://dvmn.org/modules/async-python/).


## License

This project is licensed under the GNU GPLv3  License - see the [LICENSE](LICENSE) file for details

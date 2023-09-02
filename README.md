# Spaceship Game

This is a classic arcade-style game created in Go. In this game, players take command of a spaceship as it flies through a meteor shower. The goal is to protect the spaceship by destroying incoming meteorites while navigating to avoid collisions.

## Getting Started

### Installation

#### Linux
- Download a zip file `spaceship-linux-latest.zip` from [a release page](https://github.com/AndreyAD1/spaceship/releases).
- Unpack the downloaded archive
- Open your terminal
- Run the game
    ```shell
    $ <path-to-unpacked-folder/spaceship.exe
    ```
If everything works fine, a console will display a game screen.
Use the arrow keys to navigate and press the space key to shoot. 

![presentation](presentation.gif)

#### Windows
A Windows version struggles with perfomance issues that affect the gaming experience. However, you can still launch the Windows version by following these steps:
- Download a zip file `spaceship-windows-latest.zip` from [a release page](https://github.com/AndreyAD1/spaceship/releases).
- Unpack the downloaded archive
- Run cmd.exe or Windows PowerShell
- Run the game
  ```shell
  <path-to-unpacked-folder>\spaceship.exe
  ```
  For example, if the game was unpacked to `C:\Games\spacehip-windows-latest`,
  enter `C:\Games\spacehip-windows-latest\spaceship.exe` and press `Enter`.
  To exit, press `Escape`.


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

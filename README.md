# Spaceship Game

This is an old-fashioned game written in Go. 
A user controls a spaceship flying through a meteor shower.
The user has to destroy meteorites and avoid collisions.

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
A user can move a spaceship using arrow keys and shoot using a space key.

![presentation](presentation.gif)

#### Windows
A Windows version struggles with perfomance issues and 
no enjoyable game experience is possible. Anyway, you can launch the Windows version
following these steps:
- Download a zip file `spaceship-windows-latest.zip` from [a release page](https://github.com/AndreyAD1/spaceship/releases).
- Unpack the downloaded archive
- Run cmd.exe or Windows PowerShell
- Run the game
  ```shell
  <path-to-unpacked-folder>\spaceship.exe
  ```
  For example, if the game was unpacked to `C:\Games\spacehip-windows-latest`,
  enter `C:\Games\spacehip-windows-latest\spaceship.exe` and press `Enter`.
  To exit the game press `Escape`.


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

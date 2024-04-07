
# La Vie En Go

"La Vie En Go" is a concise and efficient implementation of [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) in Go (Golang). Started as a weekend project; aimed to provide a simple (and fun!) platform for exploring and visualising this cellular automaton in a terminal-based environment.

## Features

- **Terminal-Based Visualization**: Watch the evolution of cells directly in your terminal.
- **Custom Initial Configurations**: Input your starting configurations to explore different patterns and behaviors.
- **Stop Condition**: The game halts automatically when it detects no further changes or no life within the visible area, allowing for finite observation of the game.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (Golang) installed on your system ([installation guide](https://go.dev/doc/install))

### Installing

1. Clone the repository:

```bash
git clone https://github.com/ahmad-alkadri/LaVieEnGo.git
```

2. Navigate to the cloned directory:

```bash
cd LaVieEnGo
```

3. Run the game:

For interactive mode, where you can input your initial cells:

```bash
go run main.go
```

To run with predefined cell coordinates:

```bash
go run main.go -c "x1 y1, x2 y2, ..."
```

## Usage

After starting the game, you'll be prompted to input the initial live cells' coordinates, or you can use the `-c` flag to specify them as arguments. The game progresses automatically, showing each new generation in your terminal.

### Example Patterns

- **Block**: `-c "1 1, 1 2, 2 1, 2 2"`
- **Glider**: `-c "1 2, 2 3, 3 1, 3 2, 3 3"`
- **Blinker**: `-c "2 1, 2 2, 2 3"`

There are also some patterns in the `examples` folder.
Example on how to run them:

```bash
go run main.go < examples/heart.txt
```

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Acknowledgments

- John Conway, for the original concept of the Game of Life.


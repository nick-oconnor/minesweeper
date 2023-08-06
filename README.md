# Minesweeper

![Gameplay](gameplay.png "Gameplay")

A minesweeper simulator and solver written in Go. The solver has three modes, executed in the following order of
preference:

### Single Space Deduction

* If a revealed space results in a 100% probability of an unknown space containing a mine, the space is flagged
* If a revealed space results in a 0% probability of an unknown space containing a mine, the space is revealed

The calculation is re-run until all possible moves are made. When no more moves are found, the solver moves on to
multi-space deduction.

### Multi Space Deduction[^1]

Take the following layout:

```
  1 2 3
A| |1|-|
B|1|2|-|
C|-|-|-|
```

The layout has the following known spaces with unknown neighbors:

```
A2,B2,B1
```

And the following unknown spaces with known neighbors:

```
A3,B3,C1,C2,C3
```

An equation is generated for each of the known spaces which represents its affect on each of the unknown spaces. The
coefficients on the LHS indicate whether the unknown space is affected, the coefficients on the RHS are the number of
unknown mines adjacent to each of the revealed spaces:

```
1(A3) + 1(B3) + 0(C1) + 0(C2) + 0(C3) = 1(A2)
1(A3) + 1(B3) + 1(C1) + 1(C2) + 1(C3) = 2(B2)
0(A3) + 0(B3) + 1(C1) + 1(C2) + 0(C3) = 1(B1)
```

The resulting matrix of coefficients is then converted to reduced row echelon form:

```
1 1 0 0 0 = 1
0 0 1 1 0 = 1
0 0 0 0 1 = 0
```

The resulting coefficients are then used to either flag or reveal spaces. The generalized algorithm is:

* If the RHS coefficient equals the sum of all positive LHS coefficients
    * All non-zero LHS spaces with positive coefficients are mines
    * All non-zero LHS spaces with negative coefficients are not mines
* If the RHS coefficient equals the sum of all negative LHS coefficients
    * All non-zero LHS spaces with negative coefficients are mines
    * All non-zero LHS spaces with positive coefficients are not mines

In the above example, the sum of the negative LHS coefficients for row 3 is 0, which equals the RHS (also 0). In this
case all positive coefficients (namely `C3`) are not mines.

The calculation is re-run until all possible moves are made. When no more moves are found, the solver moves on to
guessing.

### Guessing

The solver makes the following decisions:

* If an unknown corner space exists, it reveals it
* If an unknown edge space exists, it reveals it
* If an unknown center space exists, it reveals it

A better algorithm would calculate all possible solutions in which each space does not contain a mine and use that to
reveal the space with the lowest probability of containing a mine[^2].

## Performance

The solver simulates games simultaneously using all available cores on the host. The following stats were collected
on an i9-11900H CPU:

For beginner fields:

```
$ time ./minesweeper -width 9 -height 9 -mines 10 -games 1000000
Games Simulated: 1000000, Won: 90.7%, Moves/Win: 30.2, Guesses/Win: 1.66, CPU Cycles/Game: 1.77e+05

real	0m22.617s
user	2m50.005s
sys	0m7.291s
```

For intermediate fields:

```
$ time ./minesweeper -width 16 -height 16 -mines 40 -games 1000000
Games Simulated: 1000000, Won: 73.8%, Moves/Win: 122.2, Guesses/Win: 2.14, CPU Cycles/Game: 6.31e+05

real	0m53.740s
user	10m12.862s
sys	0m18.146s
```

For expert fields:

```
$ time ./minesweeper -width 30 -height 16 -mines 99 -games 1000000
Games Simulated: 1000000, Won: 26.0%, Moves/Win: 306.8, Guesses/Win: 4.01, CPU Cycles/Game: 1.92e+06

real	2m17.572s
user	31m15.589s
sys	0m42.684s
```

The code is structured to isolate the simulation code from the solver such that illegal access to the field or spaces
will trigger a panic. This ensures that development of the solver did not rely on internal field state.

## Building

You can either compile and execute a binary with:

```
go build -v
./minesweeper --help
```

Or build a container and execute it with:

```
docker build -t minesweeper .
docker run -it --rm minesweeper --help
```

The available flags are:

```
  -duration duration
        visualize move duration (default 500ms)
  -games int
        number of games (default 1000)
  -height int
        height of the field (default 16)
  -mines int
        number of mines (default 99)
  -visualize
        visualize gameplay
  -width int
        width of the field (default 30)
```

If `visualize` is specified, only a single game is played.

The `solver` unit tests display the raw and reduced matrices used for multi space deduction:

```
cd solver/
go test
```

```
| - | - | - | - | - |
| - | - | - | - | - |
| 1 | 1 | 1 | 1 | 1 |
|   |   |   |   |   |

  1,2  2,2  3,2  4,2  5,2 
  1.0  1.0  0.0  0.0  0.0 |  1.0 1,3
  1.0  1.0  1.0  0.0  0.0 |  1.0 2,3
  0.0  1.0  1.0  1.0  0.0 |  1.0 3,3
  0.0  0.0  1.0  1.0  1.0 |  1.0 4,3
  0.0  0.0  0.0  1.0  1.0 |  1.0 5,3

  1,2  2,2  3,2  4,2  5,2 
  1.0  0.0  0.0  0.0  1.0 |  1.0 1,3
  0.0  1.0  0.0  0.0 -1.0 |  0.0 3,3
  0.0  0.0  1.0  0.0  0.0 |  0.0 2,3
  0.0  0.0  0.0  1.0  1.0 |  1.0 4,3
  0.0  0.0  0.0  0.0  0.0 |  0.0 5,3

| - | - | - | - | - |
| - | - | 1 | - | - |
| 1 | 1 | 1 | 1 | 1 |
|   |   |   |   |   |
3,2 revealed: success
```

[^1]: https://massaioli.wordpress.com/2013/01/12/solving-minesweeper-with-matricies/

[^2]: https://dash.harvard.edu/bitstream/handle/1/14398552/BECERRA-SENIORTHESIS-2015.pdf
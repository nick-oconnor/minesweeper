mod field;
mod matrix;
mod solver;

use clap::Parser;
use field::Field;
use solver::Solver;
use std::io::{self, Write};
use std::sync::{mpsc, Arc, Mutex};
use std::thread;

#[derive(Parser, Debug)]
#[command(name = "minesweeper")]
#[command(about = "A minesweeper simulator and solver", long_about = None)]
struct Args {
    #[arg(short, long, default_value_t = 30)]
    width: usize,

    #[arg(short = 'H', long, default_value_t = 16)]
    height: usize,

    #[arg(short, long, default_value_t = 99)]
    mines: usize,

    #[arg(short, long, default_value_t = 1000)]
    games: usize,

    #[arg(short, long)]
    visualize: bool,

    #[arg(short, long)]
    progress: bool,
}

fn simulate(width: usize, height: usize, mine_count: usize, game_count: usize, progress: bool) {
    let (game_tx, game_rx) = mpsc::channel();
    let (result_tx, result_rx) = mpsc::channel();
    let game_rx = Arc::new(Mutex::new(game_rx));

    let num_workers = thread::available_parallelism()
        .map(|n| n.get())
        .unwrap_or(1);

    // Spawn worker threads
    let mut workers = Vec::new();
    for _ in 0..num_workers {
        let game_rx = Arc::clone(&game_rx);
        let result_tx = result_tx.clone();

        let handle = thread::spawn(move || {
            loop {
                let msg = game_rx.lock().unwrap().recv();
                if msg.is_err() {
                    break;
                }
                let field = Field::new(width, height, mine_count).expect("Failed to create field");
                let mut solver = Solver::new(field, false);
                let result = solver.solve();
                let _ = result_tx.send(result);
            }
        });
        workers.push(handle);
    }

    drop(result_tx);

    // Send games to workers
    thread::spawn(move || {
        for _ in 0..game_count {
            let _ = game_tx.send(());
        }
    });

    // Collect results
    let mut won_count = 0;
    let mut move_count = 0;
    let mut guess_count = 0;
    let mut games_simulated = 0;

    for result in result_rx {
        games_simulated += 1;
        if result.won {
            won_count += 1;
            move_count += result.move_count;
            guess_count += result.guess_count;
        }

        let avg_moves = if won_count > 0 {
            move_count as f64 / won_count as f64
        } else {
            0.0
        };
        let avg_guesses = if won_count > 0 {
            guess_count as f64 / won_count as f64
        } else {
            0.0
        };
        if progress {
            print!(
                "Games Simulated: {}, Win Ratio: {:.1}%, Moves/Win: {:.1}, Guesses/Win: {:.2}\r",
                games_simulated,
                won_count as f64 / games_simulated as f64 * 100.0,
                avg_moves,
                avg_guesses
            );
            io::stdout().flush().unwrap();
        }
    }

    let avg_moves = if won_count > 0 {
        move_count as f64 / won_count as f64
    } else {
        0.0
    };
    let avg_guesses = if won_count > 0 {
        guess_count as f64 / won_count as f64
    } else {
        0.0
    };

    if progress {
        print!("\r");
    }
    println!(
        "Games Simulated: {}, Win Ratio: {:.1}%, Moves/Win: {:.1}, Guesses/Win: {:.2}",
        game_count,
        won_count as f64 / game_count as f64 * 100.0,
        avg_moves,
        avg_guesses
    );

    // Join all worker threads to ensure cleanup
    for worker in workers {
        let _ = worker.join();
    }
}

fn main() {
    let args = Args::parse();

    if args.visualize {
        let field = Field::new(args.width, args.height, args.mines).expect("Failed to create field");
        let mut solver = Solver::new(field, true);
        solver.solve();
        return;
    }

    simulate(args.width, args.height, args.mines, args.games, args.progress);
}

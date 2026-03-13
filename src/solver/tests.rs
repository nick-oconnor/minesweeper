use crate::field::{Field, State};
use crate::solver::Solver;
use std::collections::HashMap;

#[test]
fn test_readme() {
    let mut f = new_test_field(vec![
        vec!["R", " ", " "],
        vec![" ", " ", "M"],
        vec![" ", "M", " "],
    ]);
    solve(&mut f);
    assert_field(
        &f,
        vec![
            vec!["R", "R", "R"],
            vec!["R", "R", "F"],
            vec!["R", "F", "R"],
        ],
    );
}

#[test]
fn test_single_basic() {
    let mut f = new_test_field(vec![
        vec!["M", "M", "M"],
        vec![" ", " ", " "],
        vec!["R", " ", " "],
    ]);
    solve(&mut f);
    assert_field(
        &f,
        vec![
            vec!["F", "F", "F"],
            vec!["R", "R", "R"],
            vec!["R", "R", "R"],
        ],
    );
}

#[test]
fn test_multi_one_one() {
    let mut f = new_test_field(vec![
        vec![" ", " ", " "],
        vec![" ", "M", " "],
        vec!["R", "R", "R"],
    ]);
    solve(&mut f);
    assert_field(
        &f,
        vec![
            vec!["R", "R", "R"],
            vec!["R", "F", "R"],
            vec!["R", "R", "R"],
        ],
    );
}

#[test]
fn test_multi_one_two() {
    let mut f = new_test_field(vec![
        vec![" ", " ", " "],
        vec!["M", " ", "M"],
        vec!["R", "R", "R"],
    ]);
    solve(&mut f);
    assert_field(
        &f,
        vec![
            vec!["R", "R", "R"],
            vec!["F", "R", "F"],
            vec!["R", "R", "R"],
        ],
    );
}

#[test]
fn test_multi_one_two_corner() {
    let mut f = new_test_field(vec![
        vec![" ", "M", "M", " ", " ", "M"],
        vec![" ", " ", " ", " ", " ", "M"],
        vec!["R", " ", " ", " ", " ", " "],
    ]);
    solve(&mut f);
    assert_field(
        &f,
        vec![
            vec!["R", "F", "F", "R", "R", "F"],
            vec!["R", "R", "R", "R", "R", "F"],
            vec!["R", "R", "R", "R", "R", "R"],
        ],
    );
}

#[test]
fn test_two_two_corner() {
    let mut f = new_test_field(vec![
        vec!["M", "M", " ", "M", " "],
        vec![" ", " ", " ", " ", "M"],
        vec![" ", " ", " ", " ", " "],
        vec![" ", " ", " ", " ", "M"],
        vec!["R", " ", " ", " ", "M"],
    ]);
    solve(&mut f);
    assert_field(
        &f,
        vec![
            vec!["F", "F", "R", "F", "R"],
            vec!["R", "R", "R", "R", "F"],
            vec!["R", "R", "R", "R", "R"],
            vec!["R", "R", "R", "R", "F"],
            vec!["R", "R", "R", "R", "F"],
        ],
    );
}

#[test]
fn test_decoupled() {
    let mut f = new_test_field(vec![
        vec!["R", "R", " ", " ", " "],
        vec![" ", "M", " ", " ", " "],
        vec![" ", " ", " ", " ", " "],
        vec![" ", " ", " ", "M", " "],
        vec![" ", " ", " ", "R", "R"],
    ]);
    solve(&mut f);
    assert_field(
        &f,
        vec![
            vec!["R", "R", "R", "R", "R"],
            vec!["R", "F", "R", "R", "R"],
            vec!["R", "R", "R", "R", "R"],
            vec!["R", "R", "R", "F", "R"],
            vec!["R", "R", "R", "R", "R"],
        ],
    );
}

#[test]
fn test_probabilities1() {
    let f = new_test_field(vec![
        vec![" ", " ", "M", " ", "M"],
        vec!["F", "F", " ", " ", " "],
        vec![" ", " ", "R", " ", " "],
        vec![" ", " ", "M", " ", " "],
        vec!["R", " ", " ", " ", "R"],
    ]);
    let solver = Solver::new(f, true);
    let mut want = HashMap::new();
    want.insert(2, 2.0 / 4.0);
    want.insert(3, 2.0 / 4.0);
    want.insert(4, 2.0 / 4.0);
    want.insert(7, 2.0 / 4.0);
    want.insert(17, 2.0 / 4.0);
    want.insert(22, 2.0 / 4.0);
    assert_probabilities(&solver, want);
}

#[test]
fn test_probabilities2() {
    let f = new_test_field(vec![
        vec![" ", "M", "R", " ", "M"],
        vec![" ", " ", " ", "R", "M"],
        vec![" ", " ", " ", "F", "F"],
        vec![" ", " ", " ", " ", " "],
        vec!["R", " ", " ", " ", " "],
    ]);
    let solver = Solver::new(f, true);
    let mut want = HashMap::new();
    want.insert(0, 1.0 / 3.0);
    want.insert(1, 2.0 / 3.0);
    want.insert(3, 1.0 / 3.0);
    want.insert(4, 1.0 / 3.0);
    want.insert(9, 1.0 / 3.0);
    assert_probabilities(&solver, want);
}

#[test]
fn test_probabilities3() {
    let f = new_test_field(vec![
        vec!["M", "M", " ", " ", " "],
        vec!["M", " ", " ", " ", " "],
        vec![" ", " ", "R", "M", " "],
        vec![" ", " ", " ", " ", " "],
        vec![" ", " ", "M", " ", "R"],
    ]);
    let solver = Solver::new(f, true);
    let mut want = HashMap::new();
    want.insert(6, 6.0 / 7.0);
    want.insert(7, 6.0 / 7.0);
    want.insert(8, 6.0 / 7.0);
    want.insert(11, 6.0 / 7.0);
    want.insert(13, 6.0 / 7.0);
    want.insert(14, 1.0 / 7.0);
    want.insert(16, 6.0 / 7.0);
    want.insert(17, 6.0 / 7.0);
    want.insert(22, 1.0 / 7.0);
    assert_probabilities(&solver, want);
}

#[test]
fn test_probabilities4() {
    let f = new_test_field(vec![
        vec!["R", " ", " ", " ", "R"],
        vec!["R", "M", " ", " ", " "],
        vec!["R", "R", " ", " ", " "],
        vec![" ", "M", "R", " ", "M"],
        vec!["R", "M", "R", "F", "R"],
    ]);
    let solver = Solver::new(f, true);
    let mut want = HashMap::new();
    want.insert(1, 2.0 / 3.0);
    want.insert(6, 1.0 / 3.0);
    want.insert(15, 1.0 / 3.0);
    want.insert(16, 1.0 / 3.0);
    want.insert(18, 1.0 / 3.0);
    want.insert(19, 2.0 / 3.0);
    want.insert(21, 1.0 / 3.0);
    assert_probabilities(&solver, want);
}

#[test]
fn test_probabilities_split() {
    let f = new_test_field(vec![
        vec!["R", " ", " ", " ", " "],
        vec![" ", "M", " ", " ", " "],
        vec![" ", " ", " ", " ", " "],
        vec![" ", " ", " ", "M", " "],
        vec![" ", " ", " ", " ", "R"],
    ]);
    let solver = Solver::new(f, true);
    let mut want = HashMap::new();
    want.insert(1, 2.0 / 3.0);
    want.insert(5, 2.0 / 3.0);
    want.insert(6, 2.0 / 3.0);
    want.insert(18, 2.0 / 3.0);
    want.insert(19, 2.0 / 3.0);
    want.insert(23, 2.0 / 3.0);
    assert_probabilities(&solver, want);
}

// Helper functions

fn new_test_field(layout: Vec<Vec<&str>>) -> Field {
    use crate::field::SpaceIdx;

    let height = layout.len();
    let width = layout[0].len();
    let mut f = Field::new(width, height, 0).expect("Failed to create field");

    // Add mines
    for (row_index, row) in layout.iter().enumerate() {
        for (col_index, &char) in row.iter().enumerate() {
            if char == "M" || char == "F" {
                if let Err(e) = f.add_mine(SpaceIdx::new(row_index * width + col_index)) {
                    panic!("{}", e);
                }
            }
        }
    }

    // Reveal or flag spaces
    for (row_index, row) in layout.iter().enumerate() {
        for (col_index, &char) in row.iter().enumerate() {
            let space_idx = SpaceIdx::new(row_index * width + col_index);
            match char {
                "F" => {
                    if let Err(e) = f.flag(space_idx) {
                        panic!("{}", e);
                    }
                }
                "R" => {
                    if let Err(e) = f.reveal(space_idx) {
                        panic!("{}", e);
                    }
                }
                _ => {}
            }
        }
    }

    f
}

fn solve(f: &mut Field) {
    // Take ownership temporarily by swapping with a dummy field
    let temp_field = std::mem::replace(f, Field::new(1, 1, 0).expect("Failed to create field"));
    let mut solver = Solver::new(temp_field, true);
    let result = solver.solve();
    if !result.won {
        panic!("game lost");
    }
    // Put the field back
    *f = std::mem::replace(solver.field_mut(), Field::new(1, 1, 0).expect("Failed to create field"));
}

fn assert_field(f: &Field, layout: Vec<Vec<&str>>) {
    use crate::field::SpaceIdx;

    let height = layout.len();
    let width = layout[0].len();

    assert_eq!(f.width(), width, "field width mismatch");
    assert_eq!(f.height(), height, "field height mismatch");

    for (row_index, row) in layout.iter().enumerate() {
        for (col_index, &char) in row.iter().enumerate() {
            let space_idx = SpaceIdx::new(row_index * width + col_index);
            let space = &f.spaces()[space_idx.get()];
            let expected_state = match char {
                " " => State::Unknown,
                "F" => State::Flagged,
                "R" => State::Revealed,
                "M" => State::Unknown,
                _ => panic!("unknown assertion character: {}", char),
            };

            let actual_state = space.state();
            assert_eq!(
                actual_state, expected_state,
                "space {} state: want {:?}, found {:?}",
                space_idx.get(), expected_state, actual_state
            );
        }
    }
}

fn assert_probabilities(solver: &Solver, want: HashMap<usize, f32>) {
    use crate::field::SpaceIdx;

    let probabilities = solver.probability_per_space_test();

    for (idx, expected) in want {
        let actual = probabilities.get(&SpaceIdx::new(idx)).copied().unwrap_or(0.0);
        assert!(
            (actual - expected).abs() < 1e-10,
            "probability for {}: want {:.2}, found {:.2}",
            idx,
            expected,
            actual
        );
    }
}

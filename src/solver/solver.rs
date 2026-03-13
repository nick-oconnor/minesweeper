use super::move_type::{MoveInfo, MoveType};
use crate::field::{constants::*, Field, SpaceIdx, State};
use crate::matrix::Matrix;
use std::collections::HashMap;

pub struct GameResult {
    pub won: bool,
    pub move_count: usize,
    pub guess_count: usize,
}

pub struct Solver {
    field: Field,
    visualize: bool,
    move_queue: HashMap<SpaceIdx, MoveInfo>,
}

impl Solver {
    pub fn new(field: Field, visualize: bool) -> Self {
        Self {
            field,
            visualize,
            move_queue: HashMap::new(),
        }
    }

    #[cfg(test)]
    pub fn field_mut(&mut self) -> &mut Field {
        &mut self.field
    }

    #[cfg(test)]
    pub fn probability_per_space_test(&self) -> HashMap<SpaceIdx, f32> {
        let base_matrix = Matrix::new(&self.field);
        self.probability_per_space(&base_matrix)
    }

    pub fn solve(&mut self) -> GameResult {
        let mut game_result = GameResult {
            won: false,
            move_count: 0,
            guess_count: 0,
        };

        loop {
            let (space_idx, move_type, err) = self.find_and_execute_move();

            if space_idx.is_some() {
                game_result.move_count += 1;
                if Self::is_guess(move_type) {
                    game_result.guess_count += 1;
                }
            }

            if err.is_some() || space_idx.is_none() || self.field.unknown_spaces().is_empty() {
                game_result.won = self.field.unknown_spaces().is_empty() && err.is_none();
                return game_result;
            }
        }
    }

    fn find_and_execute_move(&mut self) -> (Option<SpaceIdx>, Option<MoveType>, Option<String>) {
        loop {
            let (space_idx, move_type, err) = self.next_move();
            if err.is_some() || space_idx.is_some() {
                return (space_idx, move_type, err);
            }

            let base_matrix = Matrix::new(&self.field);
            self.add_constrained_moves(&base_matrix);

            if !self.move_queue.is_empty() {
                continue;
            }

            if self.visualize {
                println!("no moves found by resolving constraints\n");
            }

            self.add_enumerated_moves(&base_matrix);

            if !self.move_queue.is_empty() {
                continue;
            }

            if self.visualize {
                println!("no moves found by enumeration\n");
            }

            self.add_unconstrained_move(&base_matrix);

            if self.move_queue.is_empty() {
                panic!("no moves added");
            }
        }
    }

    fn next_move(&mut self) -> (Option<SpaceIdx>, Option<MoveType>, Option<String>) {
        let keys: Vec<_> = self.move_queue.keys().cloned().collect();

        for space_idx in keys {
            let move_info = self.move_queue.remove(&space_idx).unwrap();
            let current_state = self.field.spaces()[space_idx.get()].state();

            if current_state == move_info.operation {
                continue;
            }

            let err = match move_info.operation {
                State::Flagged => self.field.flag(space_idx).err().map(|e| e.to_string()),
                State::Revealed => self.field.reveal(space_idx).err().map(|e| e.to_string()),
                State::Unknown => panic!("invalid move Unknown for space {}", space_idx.get()),
            };

            if err.is_some() && !Self::is_guess(Some(move_info.move_type)) {
                panic!("{}", err.as_ref().unwrap());
            }

            if self.visualize {
                println!(
                    "space {} {} by {}\n",
                    space_idx.get(), move_info.operation, move_info.move_type
                );
                print!("{}", self.field);
                if let Some(e) = &err {
                    println!("{}", e);
                }
            }

            return (Some(space_idx), Some(move_info.move_type), err);
        }

        (None, None, None)
    }

    fn add_constrained_moves(&mut self, base_matrix: &Matrix) {
        if base_matrix.rows.is_empty() {
            return;
        }

        let resolved_matrix = match base_matrix.resolve() {
            Ok(m) => m,
            Err(e) => panic!("{}", e),
        };

        if self.visualize {
            println!("resolving constraints\n");
            resolved_matrix.print();
        }

        if resolved_matrix.rows.is_empty() {
            return;
        }

        for cell in resolved_matrix.rows[0].lhs() {
            match cell.state {
                State::Flagged => {
                    self.move_queue
                        .insert(cell.space_idx, MoveInfo::new(State::Flagged, MoveType::Constrained));
                }
                State::Revealed => {
                    self.move_queue
                        .insert(cell.space_idx, MoveInfo::new(State::Revealed, MoveType::Constrained));
                }
                State::Unknown => {}
            }
        }
    }

    fn add_enumerated_moves(&mut self, base_matrix: &Matrix) {
        let (best_space, best_probability) = self.find_best_move(base_matrix);

        if !self.move_queue.is_empty() {
            return;
        }

        if let Some(space_idx) = best_space {
            self.add_enumeration_guess(space_idx, best_probability);
        }
    }

    fn find_best_move(&mut self, base_matrix: &Matrix) -> (Option<SpaceIdx>, f32) {
        let mut best_space = None;
        let mut best_probability = 0.0;

        for (space_idx, probability) in self.probability_per_space(base_matrix) {
            if eq(probability, 0.0) {
                self.move_queue
                    .insert(space_idx, MoveInfo::new(State::Flagged, MoveType::Enumeration));
            } else if eq(probability, 1.0) {
                self.move_queue.insert(
                    space_idx,
                    MoveInfo::new(State::Revealed, MoveType::Enumeration),
                );
            }

            if probability > best_probability {
                best_space = Some(space_idx);
                best_probability = probability;
            }
        }

        (best_space, best_probability)
    }

    fn add_enumeration_guess(&mut self, space_idx: SpaceIdx, probability: f32) {
        let mines_remaining = self.field.mine_count() - self.field.flagged_count();
        let field_probability =
            1.0 - (mines_remaining as f32 / self.field.unknown_spaces().len() as f32);

        if self.visualize {
            println!("unconstrained mine-free probability\n\n {:.2}\n", field_probability);
        }

        if probability > field_probability {
            self.move_queue.insert(
                space_idx,
                MoveInfo::new(State::Revealed, MoveType::EnumerationGuess),
            );
        }
    }

    fn add_unconstrained_move(&mut self, base_matrix: &Matrix) {
        let mut corner_space = None;
        let mut edge_space = None;
        let mut center_space = None;

        let constrained_spaces = base_matrix.constrained_spaces();
        let unconstrained_spaces_exist =
            self.field.unknown_spaces().len() > constrained_spaces.len();

        for &space_idx in self.field.unknown_spaces() {
            if unconstrained_spaces_exist && constrained_spaces.contains(&space_idx) {
                continue;
            }

            let neighbor_count = self.field.spaces()[space_idx.get()].neighbors().len();

            if neighbor_count == CORNER_NEIGHBORS {
                corner_space = Some(space_idx);
            }
            if neighbor_count == EDGE_NEIGHBORS {
                edge_space = Some(space_idx);
            }
            center_space = Some(space_idx);
        }

        if let Some(idx) = corner_space {
            self.move_queue.insert(
                idx,
                MoveInfo::new(State::Revealed, MoveType::UnconstrainedGuess),
            );
            return;
        }

        if let Some(idx) = edge_space {
            self.move_queue.insert(
                idx,
                MoveInfo::new(State::Revealed, MoveType::UnconstrainedGuess),
            );
            return;
        }

        if let Some(idx) = center_space {
            self.move_queue.insert(
                idx,
                MoveInfo::new(State::Revealed, MoveType::UnconstrainedGuess),
            );
        }
    }

    fn probability_per_space(&self, base_matrix: &Matrix) -> HashMap<SpaceIdx, f32> {
        let mut solutions_by_part = Vec::new();

        for (part_index, part) in base_matrix.split_coupled().iter().enumerate() {
            if self.visualize {
                println!("possible solutions for constraint group {}\n", part_index);
            }
            solutions_by_part.push(part.solve(self.visualize));
        }

        let mut min_flagged_sum = 0;
        let mut min_flagged_by_part = HashMap::new();

        for (part_index, solutions) in solutions_by_part.iter().enumerate() {
            let mut min_flagged = usize::MAX;
            for solution in solutions {
                if solution.flagged.len() < min_flagged {
                    min_flagged = solution.flagged.len();
                }
            }
            min_flagged_sum += min_flagged;
            min_flagged_by_part.insert(part_index, min_flagged);
        }

        let mut probabilities = HashMap::new();
        let mines_remaining = self.field.mine_count() - self.field.flagged_count();

        for (part_index, solution_part) in solutions_by_part.iter().enumerate() {
            let min_flagged_sum_other_parts = min_flagged_sum - min_flagged_by_part[&part_index];
            let mut solution_count = 0;
            let mut revealed_space_counts: HashMap<SpaceIdx, usize> = HashMap::new();

            for solution in solution_part {
                if min_flagged_sum_other_parts + solution.flagged.len() > mines_remaining {
                    continue;
                }

                for &space_idx in &solution.revealed {
                    *revealed_space_counts.entry(space_idx).or_insert(0) += 1;
                }

                solution_count += 1;
            }

            if self.visualize {
                println!(
                    "constrained mine-free probabilities for constraint group {}\n",
                    part_index
                );
            }

            for (space_idx, revealed_count) in revealed_space_counts {
                let probability = revealed_count as f32 / solution_count as f32;
                probabilities.insert(space_idx, probability);

                if self.visualize {
                    println!("{:4} {:.2}", space_idx.get(), probability);
                }
            }

            if self.visualize {
                println!();
            }
        }

        probabilities
    }

    fn is_guess(move_type: Option<MoveType>) -> bool {
        matches!(
            move_type,
            Some(MoveType::EnumerationGuess) | Some(MoveType::UnconstrainedGuess)
        )
    }
}

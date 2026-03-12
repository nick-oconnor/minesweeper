use super::cell::Cell;
use super::row::Row;
use crate::field::{constants::EPSILON, Field, SpaceIdx, State};
use std::collections::{HashMap, HashSet};

pub struct Matrix {
    pub rows: Vec<Row>,
}

pub struct Solution {
    pub flagged: Vec<SpaceIdx>,
    pub revealed: Vec<SpaceIdx>,
}

impl Matrix {
    pub fn new(field: &Field) -> Self {
        let revealed_edge_spaces: Vec<_> = field.revealed_edge_spaces().collect();
        let unknown_edge_spaces: Vec<_> = field.unknown_edge_spaces().collect();

        let mut rows = Vec::new();

        for &revealed_idx in &revealed_edge_spaces {
            let mut cells = Vec::new();
            let unknown_neighbors: HashSet<_> = field.spaces()[revealed_idx.get()]
                .unknown_neighbors()
                .iter()
                .copied()
                .collect();

            for &unknown_idx in &unknown_edge_spaces {
                let value = if unknown_neighbors.contains(&unknown_idx) {
                    1.0
                } else {
                    0.0
                };
                cells.push(Cell::new(value, unknown_idx, State::Unknown));
            }

            let space = &field.spaces()[revealed_idx.get()];
            let rhs = (space.mine_neighbor_count().unwrap_or(0) as i16
                - space.flagged_neighbor_count() as i16) as f64;
            cells.push(Cell::new(rhs, revealed_idx, State::Unknown));

            rows.push(Row::new(cells));
        }

        let mut matrix = Self { rows };
        matrix.reduce();
        matrix = matrix.remove_unconstrained();
        matrix
    }

    pub fn split_coupled(&self) -> Vec<Matrix> {
        let mut matrices = Vec::new();
        let mut src = self.copy();

        while !src.rows.is_empty() {
            let first_row = src.rows.remove(0);
            let mut dst = Matrix {
                rows: vec![first_row],
            };

            Self::coupled_rows(&mut src, &mut dst);
            dst = dst.remove_unconstrained();
            matrices.push(dst);
        }

        matrices
    }

    pub fn resolve(&self) -> Result<Matrix, String> {
        let mut matrix = self.copy();

        // Collect operations to perform
        let mut operations = Vec::new();

        for row in &matrix.rows {
            let lhs = row.constrained_lhs();
            let mut lhs_negative_sum = 0.0;
            let mut lhs_positive_sum = 0.0;

            for cell in &lhs {
                if cell.value < 0.0 {
                    lhs_negative_sum += cell.value;
                } else {
                    lhs_positive_sum += cell.value;
                }
            }

            let rhs = row.rhs().value;

            if (rhs - lhs_negative_sum).abs() < EPSILON {
                for cell in row.lhs() {
                    if cell.value < 0.0 {
                        operations.push((cell.space_idx, State::Flagged));
                    } else if cell.value != 0.0 {
                        operations.push((cell.space_idx, State::Revealed));
                    }
                }
            } else if (rhs - lhs_positive_sum).abs() < EPSILON {
                for cell in row.lhs() {
                    if cell.value < 0.0 {
                        operations.push((cell.space_idx, State::Revealed));
                    } else if cell.value != 0.0 {
                        operations.push((cell.space_idx, State::Flagged));
                    }
                }
            }
        }

        // Apply operations
        for (space_idx, state) in operations {
            match state {
                State::Flagged => matrix.flag(space_idx)?,
                State::Revealed => matrix.reveal(space_idx)?,
                State::Unknown => {}
            }
        }

        Ok(matrix)
    }

    pub fn solve(&self, visualize: bool) -> Vec<Solution> {
        let space_idx = self.most_constrained_space();

        if space_idx.is_none() {
            let mut solution = Solution {
                flagged: Vec::new(),
                revealed: Vec::new(),
            };

            if !self.rows.is_empty() {
                for cell in self.rows[0].lhs() {
                    match cell.state {
                        State::Flagged => solution.flagged.push(cell.space_idx),
                        State::Revealed => solution.revealed.push(cell.space_idx),
                        State::Unknown => {}
                    }
                }
            }

            if visualize {
                self.print();
            }

            return vec![solution];
        }

        let mut solutions = Vec::new();
        let space_idx = space_idx.unwrap();

        // Try revealing
        let mut reveal_matrix = self.copy();
        if reveal_matrix.reveal(space_idx).is_ok() {
            if let Ok(resolved) = reveal_matrix.resolve() {
                solutions.extend(resolved.solve(visualize));
            }
        }

        // Try flagging
        let mut flag_matrix = self.copy();
        if flag_matrix.flag(space_idx).is_ok() {
            if let Ok(resolved) = flag_matrix.resolve() {
                solutions.extend(resolved.solve(visualize));
            }
        }

        solutions
    }

    pub fn constrained_spaces(&self) -> HashSet<SpaceIdx> {
        let mut spaces = HashSet::new();
        for row in &self.rows {
            for cell in row.constrained_lhs() {
                spaces.insert(cell.space_idx);
            }
        }
        spaces
    }

    pub fn print(&self) {
        if self.rows.is_empty() {
            println!("(empty matrix)");
            return;
        }

        // Print header
        print!(" ");
        for cell in self.rows[0].lhs() {
            print!("{:4} ", cell.space_idx.get());
        }
        println!();

        for cell in self.rows[0].lhs() {
            let state_char = match cell.state {
                State::Unknown => "U",
                State::Flagged => "F",
                State::Revealed => "R",
            };
            print!("    {}", state_char);
        }
        println!();

        // Print rows
        for row in &self.rows {
            for (i, cell) in row.cells.iter().enumerate() {
                if i == row.cells.len() - 1 {
                    print!(
                        " | {:4} {}",
                        format!("{:.1}", cell.value),
                        cell.space_idx.get()
                    );
                } else {
                    print!(" {:4}", format!("{:.1}", cell.value));
                }
            }
            println!();
        }
        println!();
    }

    fn most_constrained_space(&self) -> Option<SpaceIdx> {
        let mut most_constrained_space = None;
        let mut most_constrained_count = 0;
        let mut constraint_counts: HashMap<SpaceIdx, usize> = HashMap::new();

        for row in &self.rows {
            for cell in row.constrained_lhs() {
                let idx = cell.space_idx;
                *constraint_counts.entry(idx).or_insert(0) += 1;
                if constraint_counts[&idx] > most_constrained_count {
                    most_constrained_space = Some(idx);
                    most_constrained_count = constraint_counts[&idx];
                }
            }
        }

        most_constrained_space
    }

    fn unconstrained_spaces(&self) -> HashSet<SpaceIdx> {
        if self.rows.is_empty() {
            return HashSet::new();
        }

        let mut spaces = HashSet::new();
        for cell in self.rows[0].lhs() {
            spaces.insert(cell.space_idx);
        }

        for row in &self.rows {
            for cell in row.constrained_lhs() {
                spaces.remove(&cell.space_idx);
            }
        }

        spaces
    }

    fn zero_rows(&self) -> HashSet<SpaceIdx> {
        let mut zero_rows = HashSet::new();
        for row in &self.rows {
            if row.constrained_lhs().is_empty() {
                zero_rows.insert(row.rhs().space_idx);
            }
        }
        zero_rows
    }

    fn coupled_rows(src: &mut Matrix, dst: &mut Matrix) {
        if src.rows.is_empty() {
            return;
        }

        let dst_constrained = dst.constrained_spaces();

        let mut to_move = Vec::new();
        for (i, row) in src.rows.iter().enumerate() {
            for cell in row.constrained_lhs() {
                if dst_constrained.contains(&cell.space_idx) {
                    to_move.push(i);
                    break;
                }
            }
        }

        if !to_move.is_empty() {
            for &i in to_move.iter().rev() {
                dst.rows.push(src.rows.remove(i));
            }
            Self::coupled_rows(src, dst);
        }
    }

    fn reduce(&mut self) {
        if self.rows.is_empty() {
            return;
        }

        let mut lead = 0;
        let row_count = self.rows.len();
        let col_count = self.rows[0].cells.len();

        for r in 0..row_count {
            if lead >= col_count {
                return;
            }

            let mut i = r;
            while self.rows[i].cells[lead].value == 0.0 {
                i += 1;
                if i == row_count {
                    i = r;
                    lead += 1;
                    if lead == col_count {
                        return;
                    }
                }
            }

            self.rows.swap(i, r);

            let div = self.rows[r].cells[lead].value;
            for cell in &mut self.rows[r].cells {
                cell.value /= div;
            }

            for j in 0..row_count {
                if j != r {
                    let sub = self.rows[j].cells[lead].value;
                    for k in 0..col_count {
                        let val = self.rows[r].cells[k].value * sub;
                        self.rows[j].cells[k].value -= val;
                    }
                }
            }

            lead += 1;
        }
    }

    fn reveal(&mut self, space_idx: SpaceIdx) -> Result<(), String> {
        for row in &mut self.rows {
            let len = row.cells.len();
            for cell in &mut row.cells[..len - 1] {
                if cell.space_idx == space_idx {
                    cell.value = 0.0;
                    cell.state = State::Revealed;
                }
            }
            row.validate()?;
        }
        Ok(())
    }

    fn flag(&mut self, space_idx: SpaceIdx) -> Result<(), String> {
        for row in &mut self.rows {
            let mut rhs_adjustment = 0.0;
            let len = row.cells.len();
            for cell in &mut row.cells[..len - 1] {
                if cell.space_idx == space_idx {
                    rhs_adjustment = cell.value;
                    cell.value = 0.0;
                    cell.state = State::Flagged;
                }
            }
            row.rhs_mut().value -= rhs_adjustment;
            row.validate()?;
        }
        Ok(())
    }

    fn remove_unconstrained(self) -> Matrix {
        if self.rows.is_empty() {
            return self;
        }
        let zero_rows = self.zero_rows();
        let unconstrained_spaces = self.unconstrained_spaces();
        self.remove_rows(&zero_rows)
            .remove_cols(&unconstrained_spaces)
    }

    fn copy(&self) -> Matrix {
        Matrix {
            rows: self
                .rows
                .iter()
                .map(|row| Row::new(row.cells.clone()))
                .collect(),
        }
    }

    fn remove_rows(self, rows: &HashSet<SpaceIdx>) -> Matrix {
        Matrix {
            rows: self
                .rows
                .into_iter()
                .filter(|row| !rows.contains(&row.rhs().space_idx))
                .collect(),
        }
    }

    fn remove_cols(self, cols: &HashSet<SpaceIdx>) -> Matrix {
        Matrix {
            rows: self
                .rows
                .into_iter()
                .map(|row| {
                    Row::new(
                        row.cells
                            .into_iter()
                            .filter(|cell| !cols.contains(&cell.space_idx))
                            .collect(),
                    )
                })
                .collect(),
        }
    }
}

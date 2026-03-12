use super::cell::Cell;

pub struct Row {
    pub cells: Vec<Cell>,
}

impl Row {
    pub fn new(cells: Vec<Cell>) -> Self {
        Self { cells }
    }

    pub fn lhs(&self) -> &[Cell] {
        &self.cells[..self.cells.len() - 1]
    }

    pub fn rhs(&self) -> &Cell {
        &self.cells[self.cells.len() - 1]
    }

    pub fn rhs_mut(&mut self) -> &mut Cell {
        let len = self.cells.len();
        &mut self.cells[len - 1]
    }

    pub fn constrained_lhs(&self) -> Vec<&Cell> {
        self.lhs()
            .iter()
            .filter(|cell| cell.value != 0.0)
            .collect()
    }

    pub fn validate(&self) -> Result<(), String> {
        let mut lhs_negative_sum = 0.0;
        let mut lhs_positive_sum = 0.0;

        for cell in self.constrained_lhs() {
            if cell.value < 0.0 {
                lhs_negative_sum += cell.value;
            } else {
                lhs_positive_sum += cell.value;
            }
        }

        let rhs_value = self.rhs().value;
        if rhs_value > lhs_positive_sum || rhs_value < lhs_negative_sum {
            return Err(format!(
                "invalid constraint for row {}",
                self.rhs().space_idx.get()
            ));
        }

        Ok(())
    }
}

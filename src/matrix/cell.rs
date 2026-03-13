use crate::field::{constants::EPSILON, SpaceIdx, State};

#[derive(Clone)]
pub struct Cell {
    pub value: f32,
    pub space_idx: SpaceIdx,
    pub state: State,
}

impl Cell {
    pub fn new(value: f32, space_idx: SpaceIdx, state: State) -> Self {
        Self {
            value,
            space_idx,
            state,
        }
    }

    pub fn eq(&self, other: f32) -> bool {
        (self.value - other).abs() < EPSILON
    }

    pub fn is_zero(&self) -> bool {
        self.eq(0.0)
    }

    pub fn is_non_zero(&self) -> bool {
        !self.is_zero()
    }
}

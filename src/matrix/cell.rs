use crate::field::{SpaceIdx, State};

#[derive(Clone)]
pub struct Cell {
    pub value: f64,
    pub space_idx: SpaceIdx,
    pub state: State,
}

impl Cell {
    pub fn new(value: f64, space_idx: SpaceIdx, state: State) -> Self {
        Self {
            value,
            space_idx,
            state,
        }
    }
}

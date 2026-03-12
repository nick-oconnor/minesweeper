use super::error::FieldError;
use super::space_idx::{constants::MAX_NEIGHBORS, SpaceIdx};
use super::state::State;
use smallvec::SmallVec;

pub struct Space {
    index: SpaceIdx,
    state: State,
    has_mine: bool,
    mine_neighbor_count: u8,
    flagged_neighbor_count: u8,
    revealed_neighbor_count: u8,
    neighbors: SmallVec<[SpaceIdx; MAX_NEIGHBORS]>,
    unknown_neighbors: SmallVec<[SpaceIdx; MAX_NEIGHBORS]>,
}

impl Space {
    pub fn new(index: SpaceIdx) -> Self {
        Self {
            index,
            state: State::Unknown,
            has_mine: false,
            mine_neighbor_count: 0,
            flagged_neighbor_count: 0,
            revealed_neighbor_count: 0,
            neighbors: SmallVec::new(),
            unknown_neighbors: SmallVec::new(),
        }
    }

    pub fn state(&self) -> State {
        self.state
    }

    pub fn neighbors(&self) -> &[SpaceIdx] {
        &self.neighbors
    }

    pub fn unknown_neighbors(&self) -> &[SpaceIdx] {
        &self.unknown_neighbors
    }

    pub fn mine_neighbor_count(&self) -> Result<u8, FieldError> {
        if self.state != State::Revealed {
            return Err(FieldError::UnknownSpaceQuery(self.index.get()));
        }
        Ok(self.mine_neighbor_count)
    }

    pub fn flagged_neighbor_count(&self) -> u8 {
        self.flagged_neighbor_count
    }

    pub fn revealed_neighbor_count(&self) -> u8 {
        self.revealed_neighbor_count
    }

    pub fn has_mine(&self) -> bool {
        self.has_mine
    }

    pub fn add_mine(&mut self) -> Result<(), FieldError> {
        if self.has_mine {
            return Err(FieldError::DuplicateMine(self.index.get()));
        }
        self.has_mine = true;
        Ok(())
    }

    pub fn remove_mine(&mut self) -> Result<(), FieldError> {
        if !self.has_mine {
            return Err(FieldError::MineNotPresent(self.index.get()));
        }
        self.has_mine = false;
        Ok(())
    }

    pub fn reveal(&mut self) -> Result<(), FieldError> {
        if self.state != State::Unknown {
            return Err(FieldError::InvalidRevealState(
                self.index.get(),
                self.state.to_string(),
            ));
        }
        self.state = State::Revealed;
        if self.has_mine {
            return Err(FieldError::RevealedMine(self.index.get()));
        }
        Ok(())
    }

    pub fn flag(&mut self) -> Result<(), FieldError> {
        if self.state != State::Unknown {
            return Err(FieldError::InvalidFlagState(
                self.index.get(),
                self.state.to_string(),
            ));
        }
        self.state = State::Flagged;
        if !self.has_mine {
            return Err(FieldError::FlaggedNonMine(self.index.get()));
        }
        Ok(())
    }

    pub fn add_neighbor(&mut self, neighbor_idx: SpaceIdx) {
        self.neighbors.push(neighbor_idx);
        self.unknown_neighbors.push(neighbor_idx);
    }

    pub fn remove_unknown_neighbor(&mut self, neighbor_idx: SpaceIdx) {
        self.unknown_neighbors.retain(|idx| *idx != neighbor_idx);
    }

    pub fn increment_mine_neighbors(&mut self) {
        self.mine_neighbor_count += 1;
    }

    pub fn decrement_mine_neighbors(&mut self) {
        self.mine_neighbor_count -= 1;
    }

    pub fn increment_revealed_neighbors(&mut self) {
        self.revealed_neighbor_count += 1;
    }

    pub fn increment_flagged_neighbors(&mut self) {
        self.flagged_neighbor_count += 1;
    }
}

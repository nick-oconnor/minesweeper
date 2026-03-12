use super::error::FieldError;
use super::space::Space;
use super::space_idx::{Dimension, SpaceIdx};
use super::state::State;
use rand::seq::SliceRandom;
use smallvec::SmallVec;
use std::collections::HashSet;
use std::num::NonZeroUsize;

pub struct Field {
    width: Dimension,
    #[allow(dead_code)] // Used in Display impl and tests
    height: Dimension,
    flagged_count: usize,
    mine_count: usize,
    spaces: Vec<Space>,
    revealed_spaces: Vec<SpaceIdx>,
    unknown_spaces: HashSet<SpaceIdx>,
    first_move: bool,
}

impl Field {
    pub fn new(width: usize, height: usize, mine_count: usize) -> Result<Self, FieldError> {
        let width = NonZeroUsize::new(width).ok_or(FieldError::InvalidDimension)?;
        let height = NonZeroUsize::new(height).ok_or(FieldError::InvalidDimension)?;

        let space_count = width.get() * height.get();
        let mut spaces = Vec::with_capacity(space_count);
        let mut unknown_spaces = HashSet::with_capacity(space_count);

        // Create all spaces
        for i in 0..space_count {
            spaces.push(Space::new(SpaceIdx::new(i)));
            unknown_spaces.insert(SpaceIdx::new(i));
        }

        // Set up neighbor references
        for i in 0..space_count {
            let row = i / width.get();
            let col = i % width.get();
            let min_row = row.saturating_sub(1);
            let max_row = (row + 1).min(height.get() - 1);
            let min_col = col.saturating_sub(1);
            let max_col = (col + 1).min(width.get() - 1);

            for neighbor_row in min_row..=max_row {
                for neighbor_col in min_col..=max_col {
                    let neighbor_idx = neighbor_row * width.get() + neighbor_col;
                    if neighbor_idx != i {
                        spaces[i].add_neighbor(SpaceIdx::new(neighbor_idx));
                    }
                }
            }
        }

        let mut field = Self {
            width,
            height,
            flagged_count: 0,
            mine_count: 0,
            spaces,
            revealed_spaces: Vec::with_capacity(space_count),
            unknown_spaces,
            first_move: true,
        };

        // Add mines
        for _ in 0..mine_count {
            field.add_random_mine(None)?;
        }

        Ok(field)
    }

    #[cfg(test)]
    pub fn width(&self) -> usize {
        self.width.get()
    }

    #[cfg(test)]
    pub fn height(&self) -> usize {
        self.height.get()
    }

    pub fn mine_count(&self) -> usize {
        self.mine_count
    }

    pub fn flagged_count(&self) -> usize {
        self.flagged_count
    }

    pub fn spaces(&self) -> &[Space] {
        &self.spaces
    }

    pub fn unknown_spaces(&self) -> &HashSet<SpaceIdx> {
        &self.unknown_spaces
    }

    pub fn reveal(&mut self, space_idx: SpaceIdx) -> Result<(), FieldError> {
        if self.first_move {
            self.first_move = false;
            if self.spaces[space_idx.get()].has_mine() {
                // Remove mine from this space
                self.spaces[space_idx.get()].remove_mine()?;
                // Collect neighbors to avoid borrow checker issues
                let neighbors: SmallVec<[SpaceIdx; 8]> =
                    self.spaces[space_idx.get()].neighbors().iter().copied().collect();
                for neighbor_idx in neighbors {
                    self.spaces[neighbor_idx.get()].decrement_mine_neighbors();
                }
                self.mine_count -= 1;
                // Add mine to random space
                self.add_random_mine(Some(space_idx))?;
            }
        }

        self.spaces[space_idx.get()].reveal()?;
        self.unknown_spaces.remove(&space_idx);
        self.revealed_spaces.push(space_idx);

        // Collect neighbors to avoid borrow checker issues
        let neighbors: SmallVec<[SpaceIdx; 8]> =
            self.spaces[space_idx.get()].neighbors().iter().copied().collect();
        for neighbor_idx in neighbors {
            self.spaces[neighbor_idx.get()].remove_unknown_neighbor(space_idx);
            self.spaces[neighbor_idx.get()].increment_revealed_neighbors();
        }

        self.recursive_reveal(space_idx)?;

        Ok(())
    }

    pub fn flag(&mut self, space_idx: SpaceIdx) -> Result<(), FieldError> {
        self.spaces[space_idx.get()].flag()?;
        self.unknown_spaces.remove(&space_idx);
        self.flagged_count += 1;

        // Collect neighbors to avoid borrow checker issues
        let neighbors: SmallVec<[SpaceIdx; 8]> =
            self.spaces[space_idx.get()].neighbors().iter().copied().collect();
        for neighbor_idx in neighbors {
            self.spaces[neighbor_idx.get()].remove_unknown_neighbor(space_idx);
            self.spaces[neighbor_idx.get()].increment_flagged_neighbors();
        }

        Ok(())
    }

    pub fn add_mine(&mut self, space_idx: SpaceIdx) -> Result<(), FieldError> {
        self.spaces[space_idx.get()].add_mine()?;
        self.mine_count += 1;

        // Collect neighbors to avoid borrow checker issues
        let neighbors: SmallVec<[SpaceIdx; 8]> =
            self.spaces[space_idx.get()].neighbors().iter().copied().collect();
        for neighbor_idx in neighbors {
            self.spaces[neighbor_idx.get()].increment_mine_neighbors();
        }

        Ok(())
    }

    pub fn revealed_edge_spaces(&self) -> impl Iterator<Item = SpaceIdx> + '_ {
        self.revealed_spaces.iter().copied().filter(|&idx| {
            let space = &self.spaces[idx.get()];
            !space.unknown_neighbors().is_empty() && space.mine_neighbor_count().unwrap_or(0) > 0
        })
    }

    pub fn unknown_edge_spaces(&self) -> impl Iterator<Item = SpaceIdx> + '_ {
        self.unknown_spaces
            .iter()
            .copied()
            .filter(|&idx| self.spaces[idx.get()].revealed_neighbor_count() > 0)
    }

    fn recursive_reveal(&mut self, space_idx: SpaceIdx) -> Result<(), FieldError> {
        let mine_count = self.spaces[space_idx.get()].mine_neighbor_count()?;
        if mine_count == 0 {
            // Use SmallVec to avoid allocation for small neighbor lists
            let unknown_neighbor_indices: SmallVec<[SpaceIdx; 8]> = self.spaces[space_idx.get()]
                .unknown_neighbors()
                .iter()
                .copied()
                .collect();

            for &neighbor_idx in &unknown_neighbor_indices {
                if self.unknown_spaces.contains(&neighbor_idx) {
                    self.spaces[neighbor_idx.get()].reveal()?;
                    self.unknown_spaces.remove(&neighbor_idx);
                    self.revealed_spaces.push(neighbor_idx);

                    // Update this neighbor's neighbors
                    let nn: SmallVec<[SpaceIdx; 8]> =
                        self.spaces[neighbor_idx.get()].neighbors().iter().copied().collect();
                    for n_idx in nn {
                        self.spaces[n_idx.get()].remove_unknown_neighbor(neighbor_idx);
                        self.spaces[n_idx.get()].increment_revealed_neighbors();
                    }
                }
            }

            for &neighbor_idx in &unknown_neighbor_indices {
                self.recursive_reveal(neighbor_idx)?;
            }
        }

        Ok(())
    }

    fn add_random_mine(&mut self, exclude: Option<SpaceIdx>) -> Result<(), FieldError> {
        let available_spaces: Vec<_> = self
            .spaces
            .iter()
            .enumerate()
            .filter(|(i, space)| Some(SpaceIdx::new(*i)) != exclude && !space.has_mine())
            .map(|(i, _)| SpaceIdx::new(i))
            .collect();

        if available_spaces.is_empty() {
            return Err(FieldError::NoAvailableSpaces);
        }

        let mut rng = rand::thread_rng();
        let chosen = *available_spaces
            .choose(&mut rng)
            .ok_or(FieldError::NoAvailableSpaces)?;
        self.add_mine(chosen)?;

        Ok(())
    }
}

impl std::fmt::Display for Field {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for (i, space) in self.spaces.iter().enumerate() {
            write!(f, "|")?;
            let content = match space.state() {
                State::Revealed => {
                    let count = space.mine_neighbor_count().unwrap_or(0);
                    if count > 0 {
                        format!(" {} ", count)
                    } else {
                        "   ".to_string()
                    }
                }
                State::Flagged => " * ".to_string(),
                State::Unknown => " - ".to_string(),
            };
            write!(f, "{}", content)?;
            if i % self.width.get() == self.width.get() - 1 {
                writeln!(f, "|")?;
            }
        }
        writeln!(f)
    }
}

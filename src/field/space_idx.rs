use std::num::NonZeroUsize;

/// Type-safe index for a space in the field
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub struct SpaceIdx(pub usize);

impl SpaceIdx {
    pub fn new(idx: usize) -> Self {
        Self(idx)
    }

    pub fn get(self) -> usize {
        self.0
    }
}

impl From<usize> for SpaceIdx {
    fn from(idx: usize) -> Self {
        Self(idx)
    }
}

impl From<SpaceIdx> for usize {
    fn from(idx: SpaceIdx) -> Self {
        idx.0
    }
}

/// Type-safe non-zero dimension
pub type Dimension = NonZeroUsize;

/// Constants for the game
pub mod constants {
    /// Epsilon for floating point comparisons
    pub const EPSILON: f64 = 1e-10;

    /// Maximum number of neighbors a space can have (corner = 3, edge = 5, center = 8)
    pub const MAX_NEIGHBORS: usize = 8;

    /// Corner space neighbor count
    pub const CORNER_NEIGHBORS: usize = 3;

    /// Edge space neighbor count
    pub const EDGE_NEIGHBORS: usize = 5;
}

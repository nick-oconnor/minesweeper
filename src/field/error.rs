use thiserror::Error;

#[derive(Error, Debug)]
pub enum FieldError {
    #[error("space {0} revealed which contains a mine")]
    RevealedMine(usize),

    #[error("mine added to space {0} which already contains a mine")]
    DuplicateMine(usize),

    #[error("mine removed from space {0} which does not contain a mine")]
    MineNotPresent(usize),

    #[error("space {0} revealed with invalid state {1}")]
    InvalidRevealState(usize, String),

    #[error("space {0} flagged with invalid state {1}")]
    InvalidFlagState(usize, String),

    #[error("space {0} flagged which does not contain a mine")]
    FlaggedNonMine(usize),

    #[error("retrieved mine neighbors for unknown space {0}")]
    UnknownSpaceQuery(usize),

    #[error("no available spaces for mine placement")]
    NoAvailableSpaces,

    #[error("invalid dimension: width and height must be greater than zero")]
    InvalidDimension,
}

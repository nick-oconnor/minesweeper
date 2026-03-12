use crate::field::State;
use std::fmt;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum MoveType {
    Constrained,
    Enumeration,
    EnumerationGuess,
    UnconstrainedGuess,
}

impl fmt::Display for MoveType {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            MoveType::Constrained => write!(f, "constraints"),
            MoveType::Enumeration => write!(f, "enumeration"),
            MoveType::EnumerationGuess => write!(f, "enumerated guess"),
            MoveType::UnconstrainedGuess => write!(f, "unconstrained guess"),
        }
    }
}

pub struct MoveInfo {
    pub operation: State,
    pub move_type: MoveType,
}

impl MoveInfo {
    pub fn new(operation: State, move_type: MoveType) -> Self {
        Self {
            operation,
            move_type,
        }
    }
}

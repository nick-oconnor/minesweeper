use std::fmt;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum State {
    Unknown,
    Flagged,
    Revealed,
}

impl fmt::Display for State {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            State::Unknown => write!(f, "unknown"),
            State::Flagged => write!(f, "flagged"),
            State::Revealed => write!(f, "revealed"),
        }
    }
}

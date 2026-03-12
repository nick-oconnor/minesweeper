FROM index.docker.io/library/rust:1.94.0-slim-trixie AS build

WORKDIR /src

# Install llvm for PGO
RUN apt update && \
    apt install -y llvm && \
    rm -rf /var/lib/apt/lists/*

# Copy dependency files
COPY Cargo.toml Cargo.lock ./

# Copy source code
COPY src/ ./src/

# Step 1: Build instrumented binary for PGO
RUN RUSTFLAGS="-Cprofile-generate=/tmp/pgo-data" \
    cargo build --release

# Step 2: Run instrumented binary to generate profile data (1000 expert games)
RUN /src/target/release/minesweeper --width 30 --height 16 --mines 99 --games 1000

# Step 3: Merge profile data
RUN llvm-profdata merge -o /tmp/pgo-data/merged.profdata /tmp/pgo-data

# Step 4: Build final optimized binary with PGO
RUN RUSTFLAGS="-Cprofile-use=/tmp/pgo-data/merged.profdata" \
    cargo build --release

FROM index.docker.io/library/debian:trixie-slim

COPY --from=build /src/target/release/minesweeper /

ENTRYPOINT ["/bin/bash", "-c", "time /minesweeper \"$@\"", "--"]

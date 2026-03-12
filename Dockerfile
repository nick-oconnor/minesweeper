FROM index.docker.io/library/golang:1.25.8-trixie AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY field/ ./field
COPY matrix/ ./matrix
COPY solver/ ./solver
COPY *.go ./

# Generate PGO profile with 1000 iterations
RUN go build -o /minesweeper-pgo && \
    /minesweeper-pgo -games 1000 -cpuprofile=/tmp/cpu.pprof && \
    go tool pprof -proto /tmp/cpu.pprof > default.pgo

# Build with PGO
RUN go build -v -o /minesweeper

FROM index.docker.io/library/debian:trixie-slim

COPY --from=build /minesweeper /

ENTRYPOINT ["/bin/bash", "-c", "time /minesweeper \"$@\"", "--"]

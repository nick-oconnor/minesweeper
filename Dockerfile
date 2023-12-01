FROM golang@sha256:4e09ca4c01cabba257942296eac2fdbe5c750c8e03b63a834d5c126f67ae488f AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY field/ ./field
COPY matrix/ ./matrix
COPY solver/ ./solver
COPY *.go ./
COPY *.pgo ./

RUN go build -v -o /minesweeper

FROM alpine@sha256:d695c3de6fcd8cfe3a6222b0358425d40adfd129a8a47c3416faff1a8aece389

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

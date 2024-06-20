FROM golang@sha256:c4fc712e0a823f4781266cccd1f2d0493b6259101224810bca2f2037602494c5 AS build

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

FROM alpine@sha256:dabf91b69c191a1a0a1628fd6bdd029c0c4018041c7f052870bb13c5a222ae76

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

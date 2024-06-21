FROM golang@sha256:7be705459d742b10b0cd7800c25a94f5d6e2895a330fe73fea20067e8664bddc AS build

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

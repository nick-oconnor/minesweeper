FROM golang@sha256:68dfce93aabedced2731fb2f799ab7c4b7191131e76317a6a0293eb8ffc861d2 AS build

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

FROM alpine@sha256:eafc1edb577d2e9b458664a15f23ea1c370214193226069eb22921169fc7e43f

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

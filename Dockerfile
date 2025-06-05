FROM golang@sha256:e5c2e59960f8636d02f77029c8f0a7a6b882f87fee8d2e4a9ce6c9ff112ed735 AS build

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

FROM alpine@sha256:08001109a7d679fe33b04fa51d681bd40b975d8f5cea8c3ef6c0eccb6a7338ce

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

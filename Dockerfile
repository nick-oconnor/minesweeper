FROM golang@sha256:26e5827d0bcee61db5b976014a99018c1ff8f8e6248de8f3ee85bf12229c4b34 AS build

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

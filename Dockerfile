FROM golang@sha256:be1cf73ca9fbe9c5108691405b627cf68b654fb6838a17bc1e95cc48593e70da AS build

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

FROM golang@sha256:44a637dca7f3632232d2a21e6d5091ba5e2bbdc4830be97f54eee76c204dee7f AS build

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

FROM index.docker.io/library/golang@sha256:c3dc5d5e8cf34ccb2172fb8d1aa399aa13cd8b60d27bba891d18e3b436a0c5f6 AS build

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

FROM index.docker.io/library/alpine@sha256:85f2b723e106c34644cd5851d7e81ee87da98ac54672b29947c052a45d31dc2f

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

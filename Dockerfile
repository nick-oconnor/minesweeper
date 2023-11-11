FROM golang@sha256:f475434ea2047a83e9ba02a1da8efc250fa6b2ed0e9e8e4eb8c5322ea6997795 AS build

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

FROM alpine@sha256:48d9183eb12a05c99bcc0bf44a003607b8e941e1d4f41f9ad12bdcc4b5672f86

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

FROM index.docker.io/library/golang@sha256:8aebdc2031bd790173104c00a69a5652658f9cfe14ae853a0264f2909d67fba7 AS build

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

FROM index.docker.io/library/alpine@sha256:eafc1edb577d2e9b458664a15f23ea1c370214193226069eb22921169fc7e43f

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

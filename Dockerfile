FROM golang@sha256:a8586fcec2ff60e2721a68a64a77caf6194dc1b48fc5e63384c362b4e393b37f AS build

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

FROM alpine@sha256:eddacbc7e24bf8799a4ed3cdcfa50d4b88a323695ad80f317b6629883b2c2a78

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

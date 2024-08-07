FROM golang@sha256:8ee9b9e11ef79e314a7584040451a6df8e72a66712e741bf75951e05e587404e AS build

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

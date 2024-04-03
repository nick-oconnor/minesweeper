FROM golang@sha256:84418d968b61fb3ade08bd7fdb446a2b1cd1769dd96049a8e9e56b8d9c3d5f22 AS build

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

FROM alpine@sha256:6457d53fb065d6f250e1504b9bc42d5b6c65941d57532c072d929dd0628977d0

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

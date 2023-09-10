FROM golang@sha256:0c860c7ceba62231d0f99fb92e9d7c1577f26fea794a12c75756a8f64b146e45 AS build

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

FROM alpine@sha256:c5c5fda71656f28e49ac9c5416b3643eaa6a108a8093151d6d1afc9463be8e33

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

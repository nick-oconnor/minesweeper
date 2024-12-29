FROM golang@sha256:052793ea3143a235a5b2d815ccead8910cfe547b36a1f4c8b070015b89da5eab AS build

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

FROM alpine@sha256:2c43f33bd1502ec7818bce9eea60e062d04eeadc4aa31cad9dabecb1e48b647b

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

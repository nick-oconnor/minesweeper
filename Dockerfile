FROM golang@sha256:c4fc712e0a823f4781266cccd1f2d0493b6259101224810bca2f2037602494c5 AS build

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

FROM alpine@sha256:216266c86fc4dcef5619930bd394245824c2af52fd21ba7c6fa0e618657d4c3b

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

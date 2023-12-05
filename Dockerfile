FROM golang@sha256:38d710f811d4b304c0db6e5af7f35d565c3555da77821a26f9a199d8fdf22ade AS build

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

FROM alpine@sha256:d695c3de6fcd8cfe3a6222b0358425d40adfd129a8a47c3416faff1a8aece389

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

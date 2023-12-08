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

FROM alpine@sha256:13b7e62e8df80264dbb747995705a986aa530415763a6c58f84a3ca8af9a5bcd

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

FROM golang@sha256:0f7e0b1413ed81da9abecb9a397c3cc12d7f078fb39086399cb0e324e0b86297 AS build

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

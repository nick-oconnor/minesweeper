FROM golang@sha256:55f716237933c85cee01748700755b4ac8736fb1ca974c9aed051691b68d6dc2 AS build

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

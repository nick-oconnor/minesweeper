FROM golang@sha256:d66e3181cf6c9883b2d98a1742a4fa1accca3c5629dbc4e65df5782861f743bb AS build

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

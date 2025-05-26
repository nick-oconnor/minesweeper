FROM golang@sha256:be1cf73ca9fbe9c5108691405b627cf68b654fb6838a17bc1e95cc48593e70da AS build

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

FROM alpine@sha256:1c4eef651f65e2f7daee7ee785882ac164b02b78fb74503052a26dc061c90474

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

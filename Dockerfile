FROM golang@sha256:b6da2ff7e4eb4c632f7f21532b775078f77a790b159c56a0a7963a1532364cf0 AS build

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

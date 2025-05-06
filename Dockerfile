FROM golang@sha256:671ee3eae824b94b8d6c138326f67ff927dd6f08ddf8a2e2c8785034039c5771 AS build

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

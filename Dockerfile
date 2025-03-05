FROM golang@sha256:44b3a1c80b6c3a3083df7893a996ecea1e45d4f4fd5682951b6a13569bdb64e3 AS build

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

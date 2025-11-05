FROM index.docker.io/library/golang@sha256:96f36e77302b6982abdd9849dff329feef03b0f2520c24dc2352fc4b33ed776d AS build

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

FROM index.docker.io/library/alpine@sha256:85f2b723e106c34644cd5851d7e81ee87da98ac54672b29947c052a45d31dc2f

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

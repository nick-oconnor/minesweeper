FROM golang@sha256:67c373c4369d0ce5e2b3785611890432cb0dc9292f76adb70233d79785726b5b AS build

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

FROM alpine@sha256:dabf91b69c191a1a0a1628fd6bdd029c0c4018041c7f052870bb13c5a222ae76

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

FROM golang@sha256:16dba6b538dbdf6bf75456f28224c5288a1c305b80f150dbcd7b40b6746a7a76 AS build

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

FROM alpine@sha256:6457d53fb065d6f250e1504b9bc42d5b6c65941d57532c072d929dd0628977d0

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

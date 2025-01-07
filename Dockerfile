FROM golang@sha256:fe60376abcd2ed517529818a5ebccc2b4493fa9c19d6f69c45deaddf3b738e72 AS build

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

FROM alpine@sha256:f3a728d5dcf0f45691478201526b30230de3a3e3b26ffe92462d0a98fcb8f4e5

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

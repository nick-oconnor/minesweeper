FROM golang@sha256:925ce614cda2bf90b311208946c371159ba162bebb3475817785a362c24fff4d AS build

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

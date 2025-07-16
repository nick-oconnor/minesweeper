FROM golang@sha256:9fadeb603e14f1f3e08bdbec6681fa14446053c498a554f3e57260bf892c487e AS build

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

FROM alpine@sha256:eafc1edb577d2e9b458664a15f23ea1c370214193226069eb22921169fc7e43f

COPY --from=build /minesweeper /

ENTRYPOINT ["time", "/minesweeper"]

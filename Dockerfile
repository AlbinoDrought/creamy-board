FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum /app/
RUN CGO_ENABLED=0 go mod download
COPY . /app/
RUN CGO_ENABLED=0 go test ./... && CGO_ENABLED=0 go build -o creamy-board

FROM gcr.io/distroless/base-debian10
COPY --from=0 /app/creamy-board /app/creamy-board
ENTRYPOINT ["/app/creamy-board"]

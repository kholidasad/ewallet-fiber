FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/ewallet ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /bin/ewallet /ewallet
EXPOSE 8080
ENV PORT=8080
CMD ["/ewallet"]

FROM golang:1.25 as builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o weather-app ./cmd/restapi

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/weather-app .
ENTRYPOINT ["./weather-app"]

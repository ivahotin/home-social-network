## Build
FROM golang:1.16-buster AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -v ./cmd/httpserver

## Deploy
FROM gcr.io/distroless/base-debian10 AS production
WORKDIR /
COPY --from=build /app/httpserver /httpserver
COPY --from=build /app/views      /views
EXPOSE 8000
USER nonroot:nonroot

CMD ["/httpserver"]
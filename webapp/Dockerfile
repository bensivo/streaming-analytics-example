FROM golang:1.22.4 as build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM debian:12
COPY --from=build /app/main /app/main
COPY --from=build /app/static /app/static
RUN chmod +x /app/main

EXPOSE 8080
CMD ["/app/main"]

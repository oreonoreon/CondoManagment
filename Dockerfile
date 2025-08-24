FROM golang:1.24.4-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/main/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/server .
COPY --from=build /app/etc ./etc
COPY --from=build /app/html ./html
CMD ["./server"]
EXPOSE 8080
FROM golang:1.20-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org
COPY . .
RUN go build -o /app ./


FROM alpine:3.18
COPY --from=build /app /app
EXPOSE 8080
CMD ["/app"]
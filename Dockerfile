FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org
COPY . .
RUN go build -o /app ./internal


FROM alpine:3.18
COPY --from=build /app /app
RUN apk add --no-cache wget
EXPOSE 8080
CMD ["/app"]
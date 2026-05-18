FROM golang:1.25.0-alpine3.22 AS build-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.22
RUN apk --no-cache add tzdata ca-certificates
WORKDIR /app
COPY --from=build-stage /server ./
EXPOSE 8080

CMD ["/server"]

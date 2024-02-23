FROM golang:1.21.4 AS build

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o run main.go

FROM alpine:latest

RUN apk --no-cache add tzdata ca-certificates

WORKDIR /app
COPY --from=build /app/run .
COPY --from=build /app/firebase-adminsdk.json .

EXPOSE 5000
CMD ["./run"]

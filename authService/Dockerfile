
FROM golang:1.19-alpine3.15 As builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/auth-svc

# Runtime Steps
FROM alpine

COPY --from=builder /app/auth-svc .

EXPOSE 8080

CMD ["./auth-svc"]
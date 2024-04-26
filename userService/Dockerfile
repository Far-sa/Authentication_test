FROM golang:1-alpine3.15 As builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/user-svc

# Runtime Steps
FROM alpine3.15

COPY --from=builder /app/user-svc .

EXPOSE 8080

CMD ["./user-svc"]

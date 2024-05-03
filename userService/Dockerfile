FROM golang:1.19-alpine3.15 As builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o /app/user-svc

# Runtime Steps
FROM alpine


COPY --from=builder /app/user-svc .

EXPOSE 5000

ENV RABBITMQ_HOST rabbitmq
ENV RABBITMQ_USER guest
ENV RABBITMQ_PASSWORD guest
ENV RABBITMQ_PORT 5672

CMD ["./user-svc"]

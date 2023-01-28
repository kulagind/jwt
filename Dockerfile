FROM golang:latest as builder

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/auth_app ./cmd/app/app.go && \
    chmod +x /app/bin/auth_app

RUN mkdir -p /app/rsa && \
    openssl genrsa -out /app/rsa/access_token_private_key.pem 2048 && \
    openssl rsa -in /app/rsa/access_token_private_key.pem -outform PEM -pubout -out /app/rsa/access_token_public_key.pem && \
    openssl genrsa -out /app/rsa/refresh_token_private_key.pem 2048 && \
    openssl rsa -in /app/rsa/refresh_token_private_key.pem -outform PEM -pubout -out /app/rsa/refresh_token_public_key.pem


FROM scratch

COPY --from=builder /app/bin/auth_app /auth_app
COPY --from=builder /app/rsa /rsa
EXPOSE 8080
CMD [ "/auth_app" ]

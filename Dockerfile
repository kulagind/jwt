FROM golang:latest as builder

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/auth_app ./cmd/app/app.go
RUN chmod +x /app/bin/auth_app

FROM scratch

COPY --from=builder /app/bin/auth_app /auth_app
EXPOSE 8080
CMD [ "/auth_app" ]

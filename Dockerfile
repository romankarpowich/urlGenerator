FROM golang:latest
WORKDIR /opt/app
COPY . .
RUN mv /opt/app/.env.example /opt/app/.env
RUN go mod tidy && go build -o ./main
CMD ["./main"]
EXPOSE 8080
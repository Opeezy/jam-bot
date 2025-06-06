FROM golang:1.24.3-bullseye

WORKDIR /app

COPY . .

RUN go mod download

RUN mkdir -p /bin && go build -o /bin/server .

CMD ["/bin/server"]
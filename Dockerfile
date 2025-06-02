FROM golang:1.24.3-bullseye

WORKDIR /bot

COPY . .

RUN go mod download

RUN mkdir -p /bin && go build -o /bin/bot .

CMD ["/bin/bot"]
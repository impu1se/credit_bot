FROM golang:1.13

RUN mkdir /credit_bot
ADD . /credit_bot/
WORKDIR /credit_bot

RUN go mod download
RUN go build -o credit_bot .

CMD ["/credit_bot/credit_bot"]

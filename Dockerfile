FROM golang

EXPOSE 80

ADD ./app ./src/github.com/ihsw/go-home/app
RUN go get ./... \
  && go get -t github.com/ihsw/go-home/app \
  && go install github.com/ihsw/go-home/app

CMD ["./bin/app"]

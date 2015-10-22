FROM golang

EXPOSE 80

ADD ./app ./src/github.com/ihsw/go-home/app
RUN go get ./... \
  && go install github.com/ihsw/go-home/app

CMD ["./bin/app"]

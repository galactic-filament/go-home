FROM golang

EXPOSE 80

ENV APP_PATH github.com/ihsw/go-home/app
ADD ./app ./src/$APP_PATH
RUN go get ./src/$APP_PATH/... \
  && go get -t $APP_PATH \
  && go install $APP_PATH

CMD ["./bin/app"]

FROM golang

EXPOSE 80

RUN apt-get update -q \
  && apt-get install -yq netcat

ENV APP_PATH github.com/ihsw/go-home/app
ADD ./app ./src/$APP_PATH
RUN go get ./src/$APP_PATH/... \
  && go get -t $APP_PATH \
  && go install $APP_PATH
ADD ./app/bin/run-app ./bin/run-app

CMD ["./bin/run-app"]

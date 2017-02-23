FROM golang

EXPOSE 80
ENV APP_PORT 80

# add app dir
ENV APP_DIR ./src/go-home
COPY . $APP_DIR

# add log dir
ENV APP_LOG_DIR /var/log/app
VOLUME $APP_LOG_DIR

# build app
ENV APP_PROJECT go-home/app
ENV VALIDATE_ENVIRONMENT_PROJECT go-home/validateEnvironment
RUN go get ./src/$APP_PROJECT/... \
  && go get -t ./src/$APP_PROJECT \
  && go get ./src/$VALIDATE_ENVIRONMENT_PROJECT/... \
  && go install $APP_PROJECT \
  && go install $VALIDATE_ENVIRONMENT_PROJECT

CMD ["$APP_DIR/bin/run-app"]

FROM golang AS builder

WORKDIR $GOPATH/src/github.com/galactic-filament/go-home/app
COPY ./app .
RUN go get github.com/galactic-filament/go-home/app/...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .


FROM scratch

EXPOSE 80
ENV APP_PORT 80

ENV APP_LOG_DIR /var/log/app
VOLUME $APP_LOG_DIR

COPY --from=builder /app ./
ENTRYPOINT ["./app"]

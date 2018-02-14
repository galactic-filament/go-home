FROM golang:1.9 AS builder

# installing dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# copying in source and building without linking
WORKDIR $GOPATH/src/github.com/galactic-filament/go-home/app
COPY ./app .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .


FROM scratch

EXPOSE 80
ENV APP_PORT 80

ENV APP_LOG_DIR /var/log/app
VOLUME $APP_LOG_DIR

COPY --from=builder /app ./
ENTRYPOINT ["./app"]

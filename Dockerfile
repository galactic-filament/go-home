FROM golang

COPY ./app /srv/app
WORKDIR /srv/app

RUN go build -v

CMD ["./app"]

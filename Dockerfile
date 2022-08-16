FROM golang:alpine
RUN mkdir /app
RUN apk --update add libc6-compat
COPY ./http2https* /app
WORKDIR /app
CMD ["/app/http2https"]

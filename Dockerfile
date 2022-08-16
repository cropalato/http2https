FROM golang:alpine
RUN mkdir /app
COPY ./http2https* /app
WORKDIR /app
CMD ["/app/http2https"]

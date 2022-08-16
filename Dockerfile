FROM golang:alpine AS builder
WORKDIR /go/src/http2https
COPY ./* .
RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o http2https *.go


FROM golang:alpine
RUN mkdir /app
RUN apk --update add libc6-compat
COPY --from=builder /go/src/http2https/http2https /app
WORKDIR /app
CMD ["/app/http2https"]

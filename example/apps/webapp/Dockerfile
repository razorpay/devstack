FROM golang:1.16.8-alpine3.14 as builder
RUN mkdir /build
ADD ./example/apps/webapp /build/
WORKDIR /build
RUN apk add build-base && go build -a -o crud-webapp .


# generate clean, final image for end users
FROM alpine:3.14
COPY --from=builder /build/crud-webapp .
COPY --from=builder /build/config/conf.default.toml /config/conf.default.toml

# executable
EXPOSE 9090
CMD [ "./crud-webapp" ]

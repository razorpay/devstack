FROM golang:1.16.8-alpine3.14 as builder
RUN mkdir /build
ADD ./example/apps/webapp /build/
WORKDIR /build
RUN go mod download
RUN apk add build-base && go build -a -o sqs-publisher-consumer .

# generate clean, final image for end users
FROM alpine:3.14
ENV SQS_URL http://host.docker.internal:4566
ENV AWS_REGION ap-south-1
ENV AWS_PROFILE localstack

COPY --from=builder /build/sqs-publisher-consumer .

# executable
CMD [ "./sqs-publisher-consumer" ]

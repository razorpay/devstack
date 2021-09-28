FROM golang:1.15-alpine
# enable Go modules support
ENV GO111MODULE=on
RUN mkdir /src
WORKDIR /src
COPY ./hooks/sqs_configurator/go.mod .
COPY ./hooks/sqs_configurator/go.sum .
RUN go mod download
ADD ./hooks/sqs_configurator/ /src/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin/sqsc main.go
RUN addgroup -S sqsc
RUN adduser -S sqsc -G sqsc
RUN chown -R sqsc:sqsc /src
RUN chmod +x ./entrypoint.sh
USER sqsc
ENTRYPOINT ["./entrypoint.sh"]

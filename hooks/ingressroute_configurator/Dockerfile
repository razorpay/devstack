FROM golang:1.15-alpine
# enable Go modules support
ENV GO111MODULE=on
RUN mkdir /src
WORKDIR /src
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
ADD . /src/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin/irc main.go
RUN chmod +x ./entrypoint.sh
RUN addgroup -S irc
RUN adduser -S irc -G irc
USER irc
ENTRYPOINT ["./entrypoint.sh"]

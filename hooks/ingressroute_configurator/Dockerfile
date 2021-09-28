FROM golang:1.15-alpine
# enable Go modules support
ENV GO111MODULE=on
RUN mkdir /src
WORKDIR /src
COPY ./hooks/ingressroute_configurator/go.mod .
COPY ./hooks/ingressroute_configurator/go.sum .
RUN go mod download
ADD ./hooks/ingressroute_configurator/ /src/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin/irc main.go
RUN addgroup -S irc
RUN adduser -S irc -G irc
RUN chown -R irc:irc /src
RUN chmod +x ./entrypoint.sh
USER irc
ENTRYPOINT ["./entrypoint.sh"]

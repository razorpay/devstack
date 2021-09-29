FROM golang:1.15-alpine
# enable Go modules support
ENV GO111MODULE=on
RUN mkdir /src
WORKDIR /src
COPY ./hooks/secret_cloner/go.mod .
COPY ./hooks/secret_cloner/go.sum .
RUN go mod download
ADD ./hooks/secret_cloner/ /src/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin/sec main.go
RUN chmod +x ./entrypoint.sh
RUN addgroup -S secc
RUN adduser -S secc -G secc
USER secc
ENTRYPOINT ["./entrypoint.sh"]

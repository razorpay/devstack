FROM golang:1.16-alpine

WORKDIR /app

COPY example/gitpod/go.mod ./
COPY example/gitpod/go.sum ./

RUN go mod download

COPY . ./
RUN go build -o /example example/gitpod/main.go

EXPOSE 8080

RUN chmod +x /example

CMD [ "/example" ]
FROM c.rzp.io/razorpay/onggi:golang-1.13-alpine3.10
# enable Go modules support
ENV GO111MODULE=on
RUN mkdir /src
WORKDIR /src
COPY ./tools/hooks/secret_cloner/go.mod .
COPY ./tools/hooks/secret_cloner/go.sum .
RUN go mod download
ADD ./tools/hooks/secret_cloner/ /src/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin/sec main.go
RUN chmod +x ./entrypoint.sh
RUN addgroup -S secc
RUN adduser -S secc -G secc
USER secc
ENTRYPOINT ["./entrypoint.sh"]

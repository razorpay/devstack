## Sample Golang SQS Producer Consumer CLI

## Pre Requisites
1. For local testing, we can setup [localstack](https://github.com/localstack/localstack)
```shell
pip install "localstack[full]"
localstack start
```
### CLI
This is a sample app that produces and consumes messages from AWS SQS
1. Run the CLI: `go run main.go`
3. go run main.go

### Documentation
Postman collection is available in `docs/postman_collection.json`

### Docker
```shell
docker build -t crud-webapp -f Dockerfile .
docker run -d -p 9090:9090 crud-webapp
```
In docker mode, it requires docker verion 18.03 and above.(Note: In `Dockerfile` we are using `host.docker.internal` to connect from the container to the localstack running on the host)


## Sample App / How to use the example

### CLI
The sample app here is a simple CRUD app which exposes a people model (refer `examples/gin/models/person.go`). It has simple operations for getting the list of people, adding, modifying and deleting these
- First run the application: `go run main.go` (Should start at port `:9090`)
- Post some data to the app:
  ```shell script
    curl -i -X POST http://localhost:9090/people -d '{ "FirstName": "Sachin", "LastName": "Tendulkar"}'
    curl  http://localhost:9090/people/ 
  ```
  
### Documentation
Postman collection is available in `docs/postman_collection.json`

### Docker
```shell
docker build -t crud-webapp -f Dockerfile .
docker run -d -p 9090:9090 crud-webapp
```
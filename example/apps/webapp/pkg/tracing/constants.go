package tracing

const (
	//DefaultServiceVersion for dev environment where deployment markers aren't available
	DefaultServiceVersion = "debug"
	//ServiceVersionKey tag used in the spans for deployment markers
	ServiceVersionKey = "service.version"
	//Environment cloud environment (stage, prod, func, canary etc)
	Environment = "environment"

	//Error
	Error        = "error"
	ErrorType    = "error.type"
	ErrorStack   = "error.stack"
	ErrorDetails = "error.details"

	//HTTP Request Specific
	HttpClientHost              = "request.client_host"
	HttpRequestMethod           = "http.method"
	HttpRequestHeaderPrefix     = "request.header"
	HttpRequestQueryParamPrefix = "request.query.param"
	HttpStatusCode              = "http.status_code"
	GinErrors                   = "gin.errors"

	//Database Specific
	DBTable  = "db.table"
	DbMethod = "db.method"
	DbError  = "db.err"
	DbCount  = "db.count"

	//Http Response Specific
	HttpResponseStatus = "response.status_code"
	HttpPath           = "http.path"
)

package api

type BaseApi struct {
	method	 	string
	endpoint 	string
	requestObject 	interface{}
	responseObject 	interface{}

	statusCode	int
	rawResponse	[]byte
	err 		error
}

func NewBaseApi(method string, endpoint string, requestObject interface{}, responseObject interface{}) *BaseApi {
	return &BaseApi{method, endpoint, requestObject, responseObject, 0, nil, nil}
}

func (this *BaseApi) RequestObject() interface{} {
	return this.requestObject
}

func (this *BaseApi) ResponseObject() interface{} {
	return this.responseObject
}

func (this *BaseApi) Method() string {
	return this.method
}

func (this *BaseApi) Endpoint() string {
	return this.endpoint
}

func (this *BaseApi) StatusCode() int {
	return this.statusCode
}

func (this *BaseApi) RawResponse() []byte {
	return this.rawResponse
}

func (this *BaseApi) Error() error {
	return this.err
}

func (this *BaseApi) SetStatusCode(statusCode int) {
	this.statusCode = statusCode
}

func (this *BaseApi) SetRawResponse(rawResponse []byte) {
	this.rawResponse = rawResponse
}

func (this *BaseApi) SetError(err error) {
	this.err = err
}

func (this *BaseApi) SetResponseObject(res interface{}) {
	this.responseObject = res
}

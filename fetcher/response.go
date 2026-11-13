package fetcher

type Response struct {
	statusCode int
	status     string
	headers    *Headers
	body       []byte
}

func (response *Response) StatusCode() int {
	return response.statusCode
}

func (response *Response) Status() string {
	return response.status
}

func (response *Response) Headers() *Headers {
	return response.headers
}

func (response *Response) Body() []byte {
	return response.body
}

func (response *Response) BodyString() string {
	return string(response.body)
}

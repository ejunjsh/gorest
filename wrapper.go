package gorest

import "net/http"

type HttpRequest struct {
	*http.Request
	PathParams map[string] string
}

func newHttpRequest(r *http.Request) *HttpRequest{
	return &HttpRequest{r,make(map[string]string)}
}

type HttpResponse struct {
	http.ResponseWriter
}

func newHttpResponse(r http.ResponseWriter) HttpResponse{
	return HttpResponse{r}
}

func (response *HttpResponse) WriteJson(json interface{}) {

}

func (response *HttpResponse) WriteXml(xml interface{}) {

}

func (response *HttpResponse) WriteFile(file interface{}) {

}
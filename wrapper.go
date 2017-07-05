package gorest

import "net/http"
import (
	"encoding/json"
	"log"
	"encoding/xml"
	"io"
	"os"
	"bufio"
)

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

func (response *HttpResponse) WriteJson(jsonObj interface{}) {
	b,err:=json.Marshal(jsonObj)
	if err!=nil {
		log.Println(err)
	}
	response.Write(b)
}

func (response *HttpResponse) WriteXml(xmlObj interface{}) {
	b,err:=xml.Marshal(xmlObj)
	if err!=nil {
		log.Println(err)
	}
	response.Write(b)
}

func (response *HttpResponse) WriteFile(filepath string) {
    f,err:= os.Open(filepath)
	if err!=nil {
		log.Println(err)
	}
	defer f.Close()
	r:= bufio.NewReader(f)
	io.Copy(response,r)
}
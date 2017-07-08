package gorest

import "net/http"
import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"bufio"
	"html/template"
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

func (response *HttpResponse) WriteString(str string) error {
	_,err:=response.Write([]byte(str))
	return err
}

func (response *HttpResponse) WriteJson(jsonObj interface{}) error {
	b,err:=json.Marshal(jsonObj)
	if err!=nil {
		return err
	}
	response.Write(b)
	return nil
}

func (response *HttpResponse) WriteXml(xmlObj interface{}) error {
	b,err:=xml.Marshal(xmlObj)
	if err!=nil {
		return err
	}
	response.Write(b)
	return nil
}

func (response *HttpResponse) WriteFile(filepath string) error {
    f,err:= os.Open(filepath)
	if err!=nil {
		return err
	}
	defer f.Close()
	r:= bufio.NewReader(f)
	io.Copy(response,r)
	return nil
}

func (response *HttpResponse) WriteTemplates(data interface{},tplPath ...string) error  {
	t, err := template.ParseFiles(tplPath...)
	if err != nil {
		return err
	}

	err = t.Execute(response, data)
	if err != nil {
		return err
	}
	return nil
}
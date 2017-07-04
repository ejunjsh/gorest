# gorest
[![Build Status](https://travis-ci.org/ejunjsh/gorest.svg?branch=master)](https://travis-ci.org/ejunjsh/gorest)

a restful go framework
## usage
````
go get github.com/ejunjsh/gorest
````
````
import "github.com/ejunjsh/gorest"
...
...
func main(){
	app:=gorest.NewApp()
	app.Get("/file", func(r *gorest.HttpRequest, w gorest.HttpResponse) {
		w.WriteFile("/Users/zhouff/abc")
	})
	app.Run(":8081")
}

...
...
func main(){
	app:=gorest.NewApp()
	app.Get("/json", func(r *gorest.HttpRequest, w gorest.HttpResponse) {
		a:= struct {
			Abc string `json:"abc"`
			Cba string `json:"cba"`
		}{"123","321"}
		w.WriteJson(a)
	})
	app.Run(":8081")
}
...
...
w.WriteXml(xmlobj)
````
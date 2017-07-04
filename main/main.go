package main

import "github.com/ejunjsh/gorest"

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

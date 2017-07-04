package gorest

import (
	"testing"

	"io/ioutil"
	"net/http"
	"fmt"
)

func TestApp_Get(t *testing.T) {
	expect:="hello world"
	a:=NewApp()
	a.Get("/helloworld/:abc/op/:bca", func(r *HttpRequest, w HttpResponse) {
		w.Write([]byte(expect+r.PathParams["abc"]+r.PathParams["bca"]))
	})
	go a.Run(":8080")

	res, err := http.Get("http://127.0.0.1:8080/helloworld/123/op/kkk/")
	if err != nil {
		t.Fatal(err)
	}


	got, gerr := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if gerr!=nil{
		t.Fatal(gerr)
	}

	fmt.Println(string(got))
	if string(got)!=expect+"123"+"kkk"{
		t.Fatal("no meet expectation.")
	}

}




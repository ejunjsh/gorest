package gorest

import (
	"net/http"
	"regexp"
	"strings"
	"fmt"
)



type App struct {
	handlers map[string]func(r *HttpRequest,w HttpResponse)
	patterns []string
	methods map[string]string
	regexps map[string]*regexp.Regexp
	pathparamanmes map[string][]string
}

type hodler struct {
	app *App
}

func NewApp() *App {
	return &App{
		handlers: make(map[string]func(r *HttpRequest,w HttpResponse)),
		patterns:make([]string,0),
		methods:make(map[string]string),
		regexps:make(map[string]*regexp.Regexp),
		pathparamanmes:make(map[string][]string),
	}
}

func(a *App) handle(method string,pattern string, handler func(r *HttpRequest,w HttpResponse)){
	a.handlers[pattern]=handler
	a.methods[pattern]=method
	a.regexps[pattern],a.pathparamanmes[pattern]=convertPatterntoRegex(pattern)
	for _,s:=range a.patterns{
		if s==pattern{
			return
		}
	}
	a.patterns=append(a.patterns,pattern)
}

func (a *App) Get(pattern string, handler func(r *HttpRequest,w HttpResponse))  {
	a.handle("GET",pattern,handler)
}

func (a *App) Post(pattern string, handler func(r *HttpRequest,w HttpResponse))  {
	a.handle("POST",pattern,handler)
}

func (a *App) Delete(pattern string, handler func(r *HttpRequest,w HttpResponse))  {
	a.handle("DELETE",pattern,handler)
}

func (a *App) Put(pattern string, handler func(r *HttpRequest,w HttpResponse))  {
	a.handle("PUT",pattern,handler)
}

func(a *App) Run(address string) error{
	fmt.Printf("Server listens on %s",address)
	err:=http.ListenAndServe(address,&hodler{app:a})
	if err!=nil{
		return err
	}
	return nil
}

func(a *App) RunTls(address string,cert string,key string) error{
	err:=http.ListenAndServeTLS(address,cert,key,&hodler{app:a})
	if err!=nil{
		return err
	}
	return nil
}

func (h *hodler) ServeHTTP(w http.ResponseWriter, r *http.Request){
   for _,p:=range h.app.patterns{
       if reg,ok:= h.app.regexps[p];ok{
		   if method,ok:=h.app.methods[p];ok&&r.Method==method{
			   if reg.Match([]byte(r.URL.Path)) {
				   matchers:=reg.FindSubmatch([]byte(r.URL.Path))
				   pathParamMap:=make(map[string]string)
				   if len(matchers)>1{
                       if pathParamNames,ok:=h.app.pathparamanmes[p];ok{
						   for i:=1;i<len(matchers);i++{
							   pathParamMap[pathParamNames[i]]=string(matchers[i])
						   }
					   }
				   }

				   request:= newHttpRequest(r)
				   request.PathParams=pathParamMap
				   response:=newHttpResponse(w)
				   if handler,ok:=h.app.handlers[p];ok{
					   handler(request,response)
				   }
			   }
		   }
	   }
   }
}

func convertPatterntoRegex(pattern string) (*regexp.Regexp,[]string) {
	b:=regexp.MustCompile(`:[a-zA-Z1-9]+`).ReplaceAll([]byte(pattern),[]byte(`([a-zA-Z1-9]+)`))
	if strings.LastIndex(string(b),"/")>0{
		b=append(b,byte('?'))
	}else {
		b=append(b,byte('/'))
		b=append(b,byte('?'))
	}
	reg:= regexp.MustCompile(string(b))
	b1:=regexp.MustCompile(`:[a-zA-Z1-9]+`).ReplaceAll([]byte(pattern),[]byte(`:([a-zA-Z1-9]+)`))
	reg1:= regexp.MustCompile(string(b1))
	matchers:=reg1.FindSubmatch([]byte(pattern))
	pathparamnames:=make([]string,0)
	if len(matchers)>0{
		for i:=0;i<len(matchers);i++{
			pathparamnames=append(pathparamnames,string(matchers[i]))
		}
	}

	return reg,pathparamnames
}


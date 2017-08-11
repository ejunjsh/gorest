package gorest

import (
	"net/http"
	"regexp"
	"strings"
	"fmt"
)



type App struct {
	handlers map[string]func(r *HttpRequest,w HttpResponse)error
	patterns []string
	methods map[string]map[string]string
	regexps map[string]*regexp.Regexp
	pathparamanmes map[string][]string
	errHandler func( err error, r *HttpRequest,w HttpResponse)
}

type hodler struct {
	app *App
}

func NewApp() *App {
	return &App{
		handlers: make(map[string]func(r *HttpRequest,w HttpResponse)error),
		patterns:make([]string,0),
		methods:make(map[string]map[string]string),
		regexps:make(map[string]*regexp.Regexp),
		pathparamanmes:make(map[string][]string),
		errHandler: func(err error, r *HttpRequest, w HttpResponse) {
			w.Write( []byte(err.Error()))
		},
	}
}

func(a *App) handle(method string,pattern string, handler func(r *HttpRequest,w HttpResponse) error){
	a.handlers[pattern]=handler
	if c, exist := a.methods[pattern]; exist {
		c[method] = method
	} else {
		c := make(map[string]string)
		c[method] = method
		a.methods[pattern] = c
	}
	a.regexps[pattern],a.pathparamanmes[pattern]=convertPatterntoRegex(pattern)
	for _,s:=range a.patterns{
		if s==pattern{
			return
		}
	}
	a.patterns=append(a.patterns,pattern)
}

func (a *App) Get(pattern string, handler func(r *HttpRequest,w HttpResponse)error)  {
	a.handle("GET",pattern,handler)
}

func (a *App) Post(pattern string, handler func(r *HttpRequest,w HttpResponse)error)  {
	a.handle("POST",pattern,handler)
}

func (a *App) Delete(pattern string, handler func(r *HttpRequest,w HttpResponse)error)  {
	a.handle("DELETE",pattern,handler)
}

func (a *App) Put(pattern string, handler func(r *HttpRequest,w HttpResponse) error)  {
	a.handle("PUT",pattern,handler)
}

func (a *App) Error(handler func(err error,r *HttpRequest,w HttpResponse))  {
	a.errHandler=handler
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
	fmt.Printf("Server listens on %s",address)
	err:=http.ListenAndServeTLS(address,cert,key,&hodler{app:a})
	if err!=nil{
		return err
	}
	return nil
}

func (h *hodler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	request:= newHttpRequest(r)
	response:=newHttpResponse(w)
	defer func() {
		if err:=recover();err!=nil{
			if e,ok:=err.(error);ok{
				h.app.errHandler(InternalError{e,""},request,response)
			}
			if e,ok:=err.(string);ok{
				h.app.errHandler(InternalError{nil,e},request,response)
			}
		}
	}()
   for _,p:=range h.app.patterns{
       if reg,ok:= h.app.regexps[p];ok{
		   if method,ok:=h.app.methods[p];ok&&r.Method==method[r.Method]{
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
				   request.PathParams=pathParamMap
				   if handler,ok:=h.app.handlers[p];ok{
					   err:=handler(request,response)
					   if err!=nil{
						   h.app.errHandler(err,request,response)
					   }
					   return
				   }
			   }
		   }
	   }
   }
	h.app.errHandler(NoFoundError{},request,response)
}




func convertPatterntoRegex(pattern string) (*regexp.Regexp,[]string) {
	b:=regexp.MustCompile(`:[a-zA-Z1-9]+`).ReplaceAll([]byte(pattern),[]byte(`([a-zA-Z1-9]+)`))
	if strings.HasSuffix(string(b),"/"){
		b=append(b,byte('?'))
	}else {
		b=append(b,byte('/'))
		b=append(b,byte('?'))
	}
	reg:= regexp.MustCompile("^"+string(b)+"$")
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


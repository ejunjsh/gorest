package rest

import (
	"net/http"
	"regexp"
	"fmt"
	"github.com/ejunjsh/gopervisor/node"
)

type Server struct {
	Address string
    Node *node.Node
}

func (s *Server) Run(){
	http.ListenAndServe(s.Address,s)
}


var pattern=regexp.MustCompile(`/gopervisor/(\w+)/(\w+)`)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request){

	if pattern.MatchString(r.URL.Path){
		matches:=pattern.FindSubmatch([]byte(r.URL.Path))
		if len(matches)==3 {
			process := string(matches[1])
			operation := string(matches[2])
			w.Write([]byte(fmt.Sprintf("process:%s \n operation:%s \n",process, operation)))
			return
		}
	}

	w.WriteHeader(404)
	w.Write([]byte("404,no found!"))

}
package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var routes = map[string]map[string]http.HandlerFunc{}

func Status(w http.ResponseWriter, code int, messages ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, code, http.StatusText(code))

	if len(messages) > 0 {
		fmt.Println(strings.Join(messages, ", "))
	}
}

func JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	bs, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	fmt.Fprint(w, string(bs))
}

func Handle(method string, pattern string, handler http.HandlerFunc) {
	_, ok := routes[pattern]
	if !ok {
		routes[pattern] = map[string]http.HandlerFunc{}
		http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			handle, ok := routes[pattern][r.Method]
			if ok {
				handle(w, r)
			} else {
				Status(w, http.StatusMethodNotAllowed)
			}
		})
	}
	routes[pattern][method] = handler
}

func Get(pattern string, handler http.HandlerFunc) {
	Handle(http.MethodGet, pattern, handler)
}

func Post(pattern string, handler http.HandlerFunc) {
	Handle(http.MethodPost, pattern, handler)
}

func Put(pattern string, handler http.HandlerFunc) {
	Handle(http.MethodPut, pattern, handler)
}

func Patch(pattern string, handler http.HandlerFunc) {
	Handle(http.MethodPatch, pattern, handler)
}

func Delete(pattern string, handler http.HandlerFunc) {
	Handle(http.MethodDelete, pattern, handler)
}

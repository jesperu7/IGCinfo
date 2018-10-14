package main

import (
	"github.com/Jesperu7/IGCinfo/handler"
	"github.com/Jesperu7/IGCinfo/struct"
	"net/http"
)


func main(){
	_struct.Db.Init()
	http.HandleFunc("/igcinfo/api/", handler.HandlerApi)
	http.HandleFunc("/igcinfo/api/igc/", handler.HandlerIgc)
	//http.HandleFunc("/igcinfo/api/igc/", handlerIdAndField)

	http.ListenAndServe("127.0.0.1:8080", nil)
}
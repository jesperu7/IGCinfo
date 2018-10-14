package main

import (
	"IGCinfo/handler"
	"net/http"
)


func main(){
	http.HandleFunc("/igcinfo/api/", handler.HandlerApi)
	http.HandleFunc("/igcinfo/api/igc/", handler.HandlerIgc)
	//http.HandleFunc("/igcinfo/api/igc/", handlerIdAndField)

	http.ListenAndServe("127.0.0.1:8080", nil)
}
package main

import (
	"net/http"
	//"github.com/marni/goigc"
)


func main(){
	http.HandleFunc("/igcinfo/api/", handlerApi)
	http.HandleFunc("/igcinfo/api/igc/", handlerIgc)
	//http.HandleFunc("/igcinfo/api/igc/", handlerIdAndField)

	http.ListenAndServe("127.0.0.1:8080", nil)
}
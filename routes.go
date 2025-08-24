package main



import (

    "github.com/gorilla/mux"
	"GDS/httpHandle"
	"net/http"
	"fmt"

)



func RegisterRoutes (r *mux.Router) {



	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	})

	
    r.HandleFunc("/download/{game}/{version}", httphandle.Download).Methods("GET")
	r.HandleFunc("/download/{game}", httphandle.Download).Methods("GET")
	r.HandleFunc("/info/{game}", httphandle.Download).Methods("GET")






	r.HandleFunc("/info/{game}", httphandle.Download).Methods("GET")

	r.HandleFunc("/list",httphandle.Download).Methods("GET")



}
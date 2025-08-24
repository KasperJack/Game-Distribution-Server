package httphandle

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)









func Info(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)      
	game := vars["game"]

	
	if game == "" {
		http.NotFound(w, r)
		return
	}


	// example response
	fmt.Fprintf(w, "requested info: %s\n", game)
}


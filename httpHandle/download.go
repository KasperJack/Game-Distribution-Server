package httphandle

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)









func Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)      
	game := vars["game"]
	version :=vars["version"] 
	
	
	if game == "" {
		http.NotFound(w, r)
		return
	}

	
	if version == "" {
		version = "latest"
	}


	// example response
	fmt.Fprintf(w, "requested download: %s v:%s\n", game,version)
}



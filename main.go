package main

import (
	"GDS/FileTree"
	"GDS/TCPServer"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)






func main() {

    r := mux.NewRouter()
	RegisterRoutes(r)

    

    t,err := FileTree.Parse("./myproject.protocol")

    if err != nil {
        log.Fatal(err)
    }
    
	files, err := t.FilesFrom("tree.py")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range files {
		fmt.Println(v)
	}


    go tcpserver.StartTCPServer()

	log.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}













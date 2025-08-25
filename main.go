package main

import (
	//"GDS/FileTree"
	"GDS/TCPServer"
	//"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)






func main() {

    r := mux.NewRouter()
	RegisterRoutes(r)

    

 /*
    t,err := FileTree.Parse("./game-repo/Cursed Mansion/conf/myproject.tree")

    if err != nil {
        log.Fatal(err)
    }
    
	f,_ := t.FilesFrom("System/KGL2.klib")
	
	for _,v := range f {
		fmt.Println(v)
	}

 

*/




    go tcpserver.StartTCPServer()

	log.Println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}













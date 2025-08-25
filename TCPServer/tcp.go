package tcpserver

import (
	"GDS/FileTree"
	"GDS/config"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
    "strings"
)








func handleClientDownload(conn net.Conn) {
    defer conn.Close()
    log.Printf("Client connected from %s", conn.RemoteAddr())




    // GET GAME NAME 
    game,_ := bufio.NewReader(conn).ReadString('\n')
    game = strings.TrimSuffix(game, "\n")
    fmt.Println(game)

    // SEND TREE 
    gametree := filepath.Join(config.GamesRepo,game, config.GameTree)

    t,err := FileTree.Parse(gametree)

    if err != nil {
        log.Println(err)
        return
    }





    

    // GAME INFO
    
    fileList, err := t.FileInfo()
    if err != nil {
        log.Println(err)
    }

            // Send the list of files as a JSON payload.
    encoder := json.NewEncoder(conn)
    if err := encoder.Encode(fileList); err != nil {
        log.Println("Error sending file list:", err)
        return
    }
    
    




  // 3️⃣ Stream each file sequentially
    for _, fileInfo := range fileList {
        filePath := filepath.Join("filesDir", fileInfo.Name)
        file, err := os.Open(filePath)
        if err != nil {
            log.Printf("Error opening file %s: %v", fileInfo.Name, err)
            continue
        }

        // Use a large buffer for higher throughput (e.g., 4 MB)
        buf := make([]byte, 4*1024*1024)
        if _, err := io.CopyBuffer(conn, file, buf); err != nil {
            log.Printf("Error streaming file %s: %v", fileInfo.Name, err)
            file.Close()
            break // Stop sending if network error occurs
        }

        file.Close()
        log.Printf("Finished sending %s (%d bytes)", fileInfo.Name, fileInfo.Size)
    }

    log.Printf("All files sent to %s", conn.RemoteAddr())
}








func StartTCPServer() {


    listener, err := net.Listen("tcp", ":5050")
    if err != nil {
        //log.Fatalf("Failed to start server: %v", err)
		panic(err)
    }

    defer listener.Close()
    log.Println("TCP server listening on :5050...")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Error accepting connection: %v", err)
            continue
        }
        go handleClientDownload(conn)
    }
}





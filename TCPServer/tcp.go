package tcpserver

import (
    "encoding/json"
    "log"
    "net"
    "os"
    "path/filepath"
    "io"
)




const filesDir = "./files_to_download"

// FileInfo holds metadata about a file.
type FileInfo struct {
    Name string
    Size int64
}






func handleClientDownload(conn net.Conn) {
    defer conn.Close()
    log.Printf("Client connected from %s", conn.RemoteAddr())

    // 1. Get the list of files and send metadata.
    files, err := os.ReadDir(filesDir)
    if err != nil {
        log.Println("Error reading files directory:", err)
        return
    }

    var fileList []FileInfo
    for _, file := range files {
        if !file.IsDir() {
            info, _ := file.Info()
            fileList = append(fileList, FileInfo{Name: file.Name(), Size: info.Size()})
        }
    }

    // Send the list of files as a JSON payload.
    encoder := json.NewEncoder(conn)
    if err := encoder.Encode(fileList); err != nil {
        log.Println("Error sending file list:", err)
        return
    }
    



  // 3️⃣ Stream each file sequentially
    for _, fileInfo := range fileList {
        filePath := filepath.Join(filesDir, fileInfo.Name)
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





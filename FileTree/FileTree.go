package FileTree

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)


type TreeManifest struct {
	protocolVersion string
	rootName        string
	rootDir         string
	blob            []byte
	generated       time.Time
	TotalDirs       int
	TotalFiles      int
	dirs            []string
	files           []string
	
}



func Parse(path string) (*TreeManifest, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil,err
	}
	m := &TreeManifest{}
	m.blob = data



	//f, _ := os.Open(path)
	//defer f.Close()

	parent := filepath.Dir(path)
	gameDir := filepath.Dir(parent)


	
	m.rootDir = gameDir
	scanner := bufio.NewScanner(bytes.NewReader(data))

	inDirs, inFiles := false, false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "PROTOCOL_VERSION:"):
			m.protocolVersion = strings.TrimPrefix(line, "PROTOCOL_VERSION:")
		case strings.HasPrefix(line, "ROOT_NAME:"):
			m.rootName = strings.TrimPrefix(line, "ROOT_NAME:")
		case strings.HasPrefix(line, "GENERATED:"):
			t, _ := time.Parse(time.RFC3339Nano, strings.TrimPrefix(line, "GENERATED:"))
			m.generated = t
		case strings.HasPrefix(line, "TOTAL_DIRS:"):
			fmt.Sscanf(line, "TOTAL_DIRS:%d", &m.TotalDirs)
		case strings.HasPrefix(line, "TOTAL_FILES:"):
			fmt.Sscanf(line, "TOTAL_FILES:%d", &m.TotalFiles)
		case line == "BEGIN_DIRECTORIES":
			inDirs = true
		case line == "BEGIN_FILES":
			inDirs = false
			inFiles = true
		case line == "END_MANIFEST":
			inFiles = false
		default:
			if inDirs && strings.HasPrefix(line, "DIR:") {
				parts := strings.SplitN(line, ":", 3)
				if len(parts) > 1 {
					m.dirs = append(m.dirs, parts[1])
				}
			}
			if inFiles && strings.HasPrefix(line, "FILE:") {
				m.files = append(m.files, strings.TrimPrefix(line, "FILE:"))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return m, nil
}
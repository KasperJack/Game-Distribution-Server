package FileTree



import (
	"os"
	"path/filepath"
	"fmt"
)




type FileInfo struct {
    Name string
    Size int64
}

type FileBlob struct {
	Blob []uint8 
}




func (t *TreeManifest) ManifestBlob() *FileBlob {
	
	b := &FileBlob{Blob: t.blob}
	return b
}







func (t *TreeManifest) Files() []string {
	
	return t.files //paths 
}




func (t *TreeManifest) FilesFrom(Filename string) ([]string, error) {
	for i, f := range t.files {
		if f == Filename {
			return t.files[i:],nil // inclusive //paths
		}
	}
	return []string{},fmt.Errorf("file does not exist")
}









func (t *TreeManifest) FileInfo() ([]FileInfo, error) {
	var fileList []FileInfo
	for _, f := range t.files {

		fullPath := filepath.Join(t.rootDir, f)
		info, err := os.Stat(fullPath)
		if err != nil {
			return nil, err 
		}

		fileList = append(fileList, FileInfo{Name: info.Name(), Size: info.Size()})
	}

	return fileList,nil
}





func (t *TreeManifest) FileInfoFrom(Filename string) ([]FileInfo, error) {

	files, err := t.FilesFrom(Filename) 
	
	if err != nil {
		return nil,err
	}

	var fileList []FileInfo
	for _, f := range files {

		fullPath := filepath.Join(t.rootDir, f)

		info, err := os.Stat(fullPath)
		if err != nil {
			return nil, err 
		}

		fileList = append(fileList, FileInfo{Name: info.Name(), Size: info.Size()})
	}

	return fileList,nil
}








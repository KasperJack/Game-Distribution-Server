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






func (m *TreeManifest) CreateDirs(target string) error {
	for _, dir := range m.dirs {
		path := filepath.Join(target, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}


func (m *TreeManifest) CreateDirsRoot(target string) error {
	rootPath := filepath.Join(target, m.RootName)
	for _, dir := range m.dirs {
		path := filepath.Join(rootPath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}





func (m *TreeManifest) FilesFrom(Filename string) ([]string, error) {
	for i, f := range m.files {
		if f == Filename {
			return m.files[i:],nil // inclusive
		}
	}
	return []string{},fmt.Errorf("file does not exist")
}





func (m *TreeManifest) FileInfo() ([]FileInfo, error) {
	var fileList []FileInfo
	for _, f := range m.files {


		info, err := os.Stat(f)
		if err != nil {
			return nil, err 
		}

		fileList = append(fileList, FileInfo{Name: info.Name(), Size: info.Size()})
	}

	return fileList,nil
}





func (m *TreeManifest) FileInfoFrom(Filename string) ([]FileInfo, error) {

	for i, f := range m.files {
		if f == Filename {
			return 
		}
	}
	return nil,nil
}







// FilePaths returns all file paths (relative to root)
func (m *TreeManifest) FilePaths() []string {
	return m.files
}



func (m *TreeManifest) ChangeRootName(NewName string) {
	m.RootName = NewName
}
package file

import (
	"io/ioutil"
	"os"
	"strings"
)

const BackupSuffix = ".ld"

func vendorDir(path string) bool {
	sps := strings.Split(path, `/`)
	return len(sps) > 1 && sps[1] == "vendor"
}

func hiddenDir(path string) bool { // just tmp impl
	sps := strings.Split(path, `/`)
	if len(sps) == 0 {
		return true
	}
	last := sps[len(sps)-1]

	if len(last) == 0 {
		return true
	}

	if len(last) == 1 {
		return false
	}

	return last[0] == '.' && last[1] != '.'
}

func ListGoFile(path string, jumpBacked bool) []string {
	return listSuffixFile(path, []string{".go"}, jumpBacked, "_test.go", "_lzd.go")
}

func isExcluded(fileName string, excludes []string) bool {
	for _, exclude := range excludes {
		if strings.HasSuffix(fileName, exclude) {
			return true
		}
	}
	return false
}

func listSuffixFile(path string, includes []string, jumpBacked bool, excludes ...string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	fileList := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if isExcluded(file.Name(), excludes) {
			continue
		}

		// backup file exists
		if jumpBacked {
			if _, err := os.Stat(path + "/" + file.Name() + BackupSuffix); !os.IsNotExist(err) {
				continue
			}
		}

		for _, include := range includes {
			if strings.HasSuffix(file.Name(), include) {
				fileList = append(fileList, path+"/"+file.Name())
			}
		}
	}
	return fileList
}

func ListGoFileByPaths(paths []string, jumpBacked bool) []string {
	ret := []string{}
	for _, path := range paths {
		fs := ListGoFile(path, jumpBacked)
		ret = append(ret, fs...)
	}
	return ret
}

func TreeDir(path string, deepth int) []string {
	if hiddenDir(path) {
		return []string{}
	}
	if vendorDir(path) {
		return []string{}
	}
	paths := []string{path}
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, fi := range fis {
		if fi.IsDir() && deepth != 0 {
			nextDepth := -1
			if deepth != -1 {
				nextDepth = deepth - 1
			}
			paths = append(paths, TreeDir(path+"/"+fi.Name(), nextDepth)...)
		}
	}

	return paths
}

func backupFileName(fileName string) string {
	sps := strings.Split(fileName, `/`)
	backupFileName := sps[len(sps)-1] + BackupSuffix
	sps[len(sps)-1] = backupFileName
	return strings.Join(sps, "/")
}

func restoreFileName(fileName string) string {
	sps := strings.Split(fileName, `/`)
	if !strings.HasSuffix(fileName, BackupSuffix) {
		return ""
	}
	restoreFileName := strings.Replace(sps[len(sps)-1], BackupSuffix, "", 1)
	sps[len(sps)-1] = restoreFileName
	return strings.Join(sps, "/")
}

func copyFile(src string, dst string) error {

	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, data, 0644)
}

type Jumper struct{}

func (j *Jumper) BackupPath(path string) error {
	files := ListGoFile(path, true)
	for _, fn := range files {
		if err := copyFile(fn, backupFileName(fn)); err != nil {
			return err
		}
	}
	return nil
}

func (j *Jumper) RestorePath(path string) error {
	files := listSuffixFile(path, []string{".go" + BackupSuffix}, false)
	for _, fn := range files {
		if err := copyFile(fn, restoreFileName(fn)); err != nil {
			return err
		}

		if err := os.Remove(fn); err != nil {
			return err
		}
	}
	return nil
}

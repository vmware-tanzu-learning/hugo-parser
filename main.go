package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	shutil "github.com/termie/go-shutil"
	yaml "gopkg.in/yaml.v2"
)

type Mapping struct {
	Name      string
	Exercises []string
}

type Structure struct {
	Mappings []Mapping
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func cleanDirectory(dir string) {
	files, _ := filepath.Glob(dir + "/*")
	for _, file := range files {
		os.RemoveAll(file)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalln("Usage: hugo-parser path/to/structure.yml path/to/source/dir path/to/target/dir ")
	}
	filename := os.Args[1]
	sourceDir := os.Args[2]
	targetDir := os.Args[3]

	mappingContents, err := ioutil.ReadFile(filename)
	check(err)
	fmt.Print(string(mappingContents))

	structure := Structure{}
	err = yaml.Unmarshal(mappingContents, &structure)
	if err != nil {
		log.Fatal("Unable to unmarshal " + filename)
	}
	learningPath := fmt.Sprintf("%s/learning-path", sourceDir)
	exercisePath := fmt.Sprintf("%s/exercises", sourceDir)

	cleanDirectory(targetDir)

	fileInfo, err := os.Lstat(sourceDir)
	mode := fileInfo.Mode()

	if !exists(targetDir) {
		os.MkdirAll(targetDir, mode)
	}

	for i, mapping := range structure.Mappings {
		index := fmt.Sprintf("%02d", i+1)
		fullSrc := learningPath + "/" + mapping.Name + "/index.md"
		fullTgtDir := targetDir + "/" + index + "-" + mapping.Name
		fullTgt := fullTgtDir + "/index.md"
		if !exists(fullTgtDir) {
			os.MkdirAll(fullTgtDir, mode)
		}

		fmt.Println(fullSrc + " to " + fullTgt)

		err = shutil.CopyFile(fullSrc, fullTgt, false)
		check(err)
		for _, exercise := range mapping.Exercises {
			fullSrc := exercisePath + "/" + exercise + "/README.md"
			fullTgt := targetDir + "/" + index + "-" + mapping.Name + "/" + exercise + ".md"
			fmt.Println(fullSrc + " to " + fullTgt)
			err = shutil.CopyFile(fullSrc, fullTgt, false)
			check(err)
		}
	}

}

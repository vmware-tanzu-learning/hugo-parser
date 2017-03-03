package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
	oldTarget := targetDir + "_old"
	if exists(oldTarget) {
		os.RemoveAll(oldTarget)
	}
	if exists(targetDir) {
		fmt.Println("Found " + targetDir + ", renaming...")
		err = os.Rename(targetDir, oldTarget)
		check(err)
	}

	err = shutil.CopyTree(learningPath, targetDir, nil)
	check(err)
	exercisePath := fmt.Sprintf("%s/exercises", sourceDir)

	for _, mapping := range structure.Mappings {
		for _, exercise := range mapping.Exercises {
			fullSrc := exercisePath + "/" + exercise + "/README.md"
			fullTgt := targetDir + "/" + mapping.Name + "/" + exercise + ".md"
			fmt.Println(fullSrc + " to " + fullTgt)
			err = shutil.CopyFile(fullSrc, fullTgt, false)
			check(err)
		}
	}

}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	mappingContents, err := os.ReadFile(filename)
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
		fullTgt := fullTgtDir + "/_index.md"
		if !exists(fullTgtDir) {
			os.MkdirAll(fullTgtDir, mode)
		}

		fmt.Println(fullSrc + " to " + fullTgt)

		err = addWeightToHeader(fullSrc, fullTgt, i+1)
		check(err)
		for j, exercise := range mapping.Exercises {
			subIndex := fmt.Sprintf("%02d", j+1)
			fullSrc := exercisePath + "/" + exercise + "/README.md"
			fullTgt := targetDir + "/" + index + "-" + mapping.Name + "/" + subIndex + "-" + exercise + ".md"
			fmt.Println(fullSrc + " to " + fullTgt)
			err = shutil.CopyFile(fullSrc, fullTgt, false)
			check(err)
			fullSrcImages := exercisePath + "/" + exercise + "/images"
			if _, err := os.Stat(fullSrcImages); err == nil {
				fullTgtImages := targetDir + "/" + index + "-" + mapping.Name + "/" + subIndex + "-" + exercise + "/images"
				err = shutil.CopyTree(fullSrcImages, fullTgtImages, nil)
				check(err)
			}
		}
	}

}

func linesFromFile(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return linesFromReader(f)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func addWeightToHeader(source, target string, weight int) error {
	lines, err := linesFromFile(source)
	if err != nil {
		return err
	}

	seenMarker := false
	fileContent := ""
	for _, line := range lines {
		fileContent += line
		fileContent += "\n"
		if strings.HasPrefix(line, "+++") && !seenMarker {
			fileContent += fmt.Sprintf("weight = %d\n", weight)
			seenMarker = true
		}
	}

	return os.WriteFile(target, []byte(fileContent), 0644)
}

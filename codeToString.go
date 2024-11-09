package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func isValidExtension(ext string) bool {
	validExtension := []string{".go", ".js", ".ts", ".csharp", ".c", ".h", ".html", ".css", ".java"}
	return slices.Contains(validExtension, ext)
}

func gatherFiles(directoryPath string) (strings.Builder, error) {
	var filesContents strings.Builder

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path != directoryPath {
			if info.IsDir() {
				otherFilesContents, err := gatherFiles(path)
				if err != nil {
					return err
				}
				filesContents.WriteString(otherFilesContents.String())

			} else {
				ext := filepath.Ext(info.Name())

				if isValidExtension(ext) {
					bytes, err := os.ReadFile(path)
					if err != nil {
						return err
					}

					content := "\n//FILE " + info.Name() + "\n" + string(bytes)

					filesContents.WriteString(content)
				}
			}
		}

		return nil
	})

	return filesContents, err
}

func putFileAsInFile(code strings.Builder) error {
	content := code.String()
	err := os.WriteFile("output.txt", []byte(content), 0755)
	if err != nil {
		fmt.Println("could not write code in output file")
		return err
	}
	return nil
}

func main() {
	if len(os.Args) == 2 {
		projectPath := os.Args[1]
		files, err := gatherFiles(projectPath)
		if err != nil {
			panic(err)
		}
		fmt.Println(files.String())
		err = putFileAsInFile(files)
		if err != nil {
			panic(err)
		}
	} else {
		panic("usage: codeToString [projectPath]")
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func checkFile(path string) (bool, error) {
	fp, err := os.Open(path)
	if err != nil {
		return false, err
	}
	scn := bufio.NewScanner(fp)
	for scn.Scan() {
		line := scn.Text()
		isCopyRight := strings.Contains(line, "Copyright (c) ")
		if isCopyRight {
			return true, nil
		}
	}
	return false, nil
}

func filesWithNoCopyRight(root string, ignorePaths []string, mustHaveSuffixes []string) ([]string, error) {

	// read all go files in the path
	// check if they do not have a copy right notice add them to a list of files to update with a copyright notice
	files := make([]string, 0)
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// check if the path is ins the ignored set and skip it

		for _, ignorePath := range ignorePaths {
			if strings.Contains(path, ignorePath) {
				return nil
			}
		}

		// check if the file has a suffix we are interested in
		hasSuffix := false
		filename := filepath.Base(path)
		for _, suffix := range mustHaveSuffixes {
			if strings.HasSuffix(filename, suffix) {
				hasSuffix = true
				break
			}

		}
		if !hasSuffix {
			return nil
		}

		// check if the file has a copy right notice
		isCopyRight, err := checkFile(path)
		if err != nil {
			return err
		}
		if !isCopyRight {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return files, nil
}

var root string
var copyright string
var dry bool

func main() {

	// parse flags
	flag.StringVar(&root, "dir", "", "the root directory to search for files")
	flag.StringVar(&copyright, "copyright", "", "the copyright notice to add to files")
	flag.BoolVar(&dry, "dry", false, "dry run")
	flag.Parse()

	// cleanup the relative path to a full path
	root, err := filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	if dry {
		fmt.Println("Checking files in:", root)
	}

	ignorePaths := []string{
		"vendor",
		"node_modules",
		".git",
		".idea",
		".vscode",
		"build",
		"dist",
		"bin",
		"tmp",
		"tests",
	}
	// only does go files (for now!)
	mustHaveSuffixes := []string{
		".go",
	}

	files, err := filesWithNoCopyRight(root, ignorePaths, mustHaveSuffixes)
	if err != nil {
		panic(err)
	}
	if dry {
		if len(files) == 0 {
			println("dry run: no files to update")
			os.Exit(0)
		}
		for _, file := range files {
			println(file)
		}
		println("dry run: would have updated the listed files")
		os.Exit(0)
	}
	// update each file with the copyright notice
	for _, file := range files {
		_, err := prependWithCopyRight(file, copyright)
		if err != nil {
			panic(err)
		}
	}
}

func prependWithCopyRight(path string, copyright string) (bool, error) {
	tmpFile, err := os.CreateTemp("", "tmp")
	if err != nil {
		return false, err
	}
	defer os.Remove(tmpFile.Name())

	_, _ = tmpFile.Write([]byte(copyright))

	originalFile, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer originalFile.Close()

	// copy the original file to the temp file
	_, err = tmpFile.ReadFrom(originalFile)
	if err != nil {
		return false, err
	}

	// rename the tmpfile to the original file path
	err = os.Rename(tmpFile.Name(), path)
	if err != nil {
		return false, err
	}
	return true, nil
}

package main

import (
	"os"
	"strings"
	"testing"
)

func TestPrePendCopyRight(t *testing.T) {
	testFileContents := []byte(`package main
func main() {
	println("hello world!")
}`)
	testFile, err := os.CreateTemp("", "tmp-copyright-workspace")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(testFile.Name(), testFileContents, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile.Name())

	expectedCopyright := []byte(`/**
 * Copyright (c) 2021, The Authors
 **/`)
	testFilePath := testFile.Name()
	testFile.Close()

	ok, err := prependWithCopyRight(testFilePath, string(expectedCopyright))
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected to be ok")
	}
	hasCopyRight, err := checkFile(testFilePath)
	if err != nil {
		t.Fatal("expected not to err", err)
	}
	if !hasCopyRight {
		currentTestFileContents, _ := os.ReadFile(testFilePath)
		t.Fatal("expected to have copy right but test file has contents:", string(currentTestFileContents))
	}
	allUpdateTestContent, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(allUpdateTestContent), string(testFileContents)) {
		t.Fatal("expected to have all the original file contents but test file has contents:\n", string(allUpdateTestContent))
	}

}

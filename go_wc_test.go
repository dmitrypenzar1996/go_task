package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomLine(l int) string {
	var runes []rune
	for i := 0; i < l-1; i++ {
		runes = append(runes, rune(randInt(65, 9000)))
		//runes[i] = rune(randInt(65, 9000))
	}
	runes = append(runes, rune('\n'))
	return string(runes)
}

func makeTestFile(lineCount int, maxLineLength int) (filename string) {
	tmpfile, err := ioutil.TempFile(".", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer tmpfile.Close()
	filename = tmpfile.Name()

	for i := 0; i < lineCount; i++ {
		lineLength := randInt(0, maxLineLength)
		line := randomLine(lineLength)
		if _, err := tmpfile.WriteString(line); err != nil {
			log.Fatal(err)
		}
	}
	return
}

var filesParamsTest = []struct {
	lineCount   int
	lineMaxSize int
}{
	{10, 10},
	{100, 20},
	{300, 400},
	{100, 1000},
	{10000, 10000},
	{5, 5},
	{0, 0},
}

func TestCountLines(t *testing.T) {
	for ind, tt := range filesParamsTest {
		filename := makeTestFile(tt.lineCount, tt.lineMaxSize)
		defer os.Remove(filename)
		count, err := countLines(filename)
		if err != nil {
			t.Logf("%v\n", err)
			t.Fail()
		}
		if count != tt.lineCount {
			t.Logf("Wrong number of lines in test file, expected: %d, get: %d\n",
				tt.lineCount, count)
			t.Fail()
		} else {
			t.Logf("Passed %d of %d\n", ind+1, len(filesParamsTest))
		}
	}
}

var notExistFileName = "__NOT__EXIST__"

func TestCountLinesOnWrong(t *testing.T) {
	_, err := countLines(notExistFileName)
	if err == nil {
		t.Logf("Return no error on not-existing file")
		t.Fail()
	}
}

func TestProcessFiles(t *testing.T) {
	var testFilesNames []string
	for _, tt := range filesParamsTest {
		filename := makeTestFile(tt.lineCount, tt.lineMaxSize)
		defer os.Remove(filename)
		testFilesNames = append(testFilesNames, filename)
	}
	outInfo := processFiles(testFilesNames)

	if len(testFilesNames) != len(outInfo.data) {
		t.Logf("Wrong number of lookuped files")
		t.Fail()
	}

	for ind, tt := range filesParamsTest {
		count := outInfo.data[testFilesNames[ind]]
		if count != tt.lineCount {
			t.Logf("Wrong number of lines in test file, expected: %d, get: %d\n",
				tt.lineCount, count)
			t.Fail()
		} else {
			t.Logf("Passed %d of %d\n", ind+1, len(filesParamsTest))
		}
	}

	wrongFiles := []string{notExistFileName}
	outInfo = processFiles(wrongFiles)
	if len(outInfo.data) != 0 {
		t.Logf("Read not existing file")
		t.Fail()
	}

	mixedFiles := append(wrongFiles, testFilesNames...)
	outInfo = processFiles(mixedFiles)
	if len(outInfo.data) != len(testFilesNames) {
		t.Logf("Non-existing file affects result")
		t.Fail()
	}
}

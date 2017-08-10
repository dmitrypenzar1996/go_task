package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type safe_dict struct {
	sync.Mutex
	data map[string]int
}

func main() {
	outInfo := processFiles(os.Args[1:])

	for _, arg := range os.Args[1:] {
		if value, ok := outInfo.data[arg]; ok {
			fmt.Printf("    %d %s\n", value, arg)
		}
	}
	if len(os.Args) > 2 {
		totalCount := countTotal(outInfo)
		fmt.Printf("    %d total\n", totalCount)
	}
}

func processFiles(filename_arr []string) (outInfo *safe_dict) {
	var wg sync.WaitGroup
	outInfo = &safe_dict{
		data: make(map[string]int),
	}

	for _, name := range filename_arr {
		wg.Add(1)
		go func(filename string) {
			processFile(filename, outInfo)
			wg.Done()
		}(name)
	}

	wg.Wait()
	return
}

func processFile(filename string, outInfo *safe_dict) {
	count, err := countLines(filename)
	if err != nil {
		fmt.Printf("test_task: %v\n", err)
		return
	}
	outInfo.Lock()
	outInfo.data[filename] = count
	outInfo.Unlock()
}

func countTotal(outInfo *safe_dict) (totalCount int) {
	totalCount = 0
	for _, value := range outInfo.data {
		totalCount += value
	}
	return
}

func countLines(fileName string) (count int, err error) {
	count = 0
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		count++
	}
	return
}

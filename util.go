package main

import (
	"mime/multipart"
	"math/rand"
	"strconv"
	"bytes"
	"bufio"
	"time"
	"sync"
	"fmt"
	"os"
	"io"
)

func openFile(filename string) (data []string) {
	var file, _ = os.Open(filename)
	var scan *bufio.Scanner = bufio.NewScanner(file)
	for scan.Scan() {
		data = append(data, scan.Text())
	}
	file.Close()
	return data
}

func createFile(slice []string, filename string) {
	var file, _ = os.Create(filename)
	for i := range slice {
		fmt.Fprintln(file, slice[i])
	}
	file.Close()
}

func writeFile(str string, filename string) {
	var file, _ = os.OpenFile(filename, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	fmt.Fprintln(file, str)
	file.Close()
}

func createChannel(slice []string) chan string {
	var channel chan string = make(chan string, len(slice))
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	go func() {
		for i := range slice {
			channel <- slice[i]
		}
		waitgroup.Done()
	}()
	waitgroup.Wait()
	return channel
}

func containsString(slice []string, str string) bool {
	for i := range slice {
		if slice[i] == str {
			return true
		}
	}
	return false
}

func removeString(slice []string, str string) []string {
	for i := range slice {
		if slice[i] ==  str {
			return append(slice[:i], slice[i + 1:]...)
		}
	}
	return slice
}

func formatNumber(number int64) string {
	var in string = strconv.FormatInt(number, 10)
	var out []byte = make([]byte, len(in)+(len(in)-1)/3)
	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

func loadImage(filepath string) (string, *bytes.Buffer) {
	var buffer *bytes.Buffer = &bytes.Buffer{}
	var writer *multipart.Writer = multipart.NewWriter(buffer)

	var formFile, _ = writer.CreateFormFile("media", "image.png")
	var file, _ = os.Open(filepath)
	io.Copy(formFile, file)
	writer.Close()
	file.Close()

	return writer.FormDataContentType(), buffer
}

func generateRandomCSRF() string {
	var bytes []byte = make([]byte, 16)
	rand.Seed(time.Now().UnixNano())
    rand.Read(bytes)
    return fmt.Sprintf("%02x", bytes)
}
package main

import (
	"testing"

	"io/ioutil"

	"os"
)

func TestGetKeyFromFile(t *testing.T) {
	var data []byte 

	fileName := "./test/test_getKeyFromFile.txt"

	data = []byte("KEY: abc123")

	err := ioutil.WriteFile(fileName, data, 0644)

	if err != nil {
		panic(err)
	}

	file, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}
	
	key := getKeyFromFile(file)

	if key != "abc123" {
		t.Fatalf("failed asserting that expected: abc123 equals actual: %s ", key)
	}
}
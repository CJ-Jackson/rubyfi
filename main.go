package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Much have an argument")
		os.Exit(1)
	}

	if os.Args[1] == "-" {
		readFromStdIn()
		return
	}

	readFromFile()
}

func readFromStdIn() {
	cRead := csv.NewReader(os.Stdin)
	csvWriter(os.Stdout, cRead)
}

func readFromFile() {
	file, err := os.OpenFile(os.Args[1], os.O_RDWR, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	cRead := csv.NewReader(file)
	buf := &bytes.Buffer{}

	csvWriter(buf, cRead)

	file.Truncate(0)
	file.Seek(0, 0)

	_, err = io.Copy(file, buf)
	if err != nil {
		log.Fatalln(err)
	}
}

func csvWriter(buf io.Writer, cRead *csv.Reader) {
	cWrite := csv.NewWriter(buf)

	records, err := cRead.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	for _, record := range records {
		record[0] = rubyfi(html.EscapeString(record[0]))
		record[1] = html.EscapeString(record[1])
		cWrite.Write(record)
	}

	cWrite.Flush()
}

func rubyfi(str string) string {
	str = strings.ReplaceAll(str, "[", "<ruby>")
	str = strings.ReplaceAll(str, "]", "</ruby>")
	str = strings.ReplaceAll(str, "{", "<rt>")
	str = strings.ReplaceAll(str, "}", "</rt>")
	return str
}

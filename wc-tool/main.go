package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	var data []byte
	fi, _ := os.Stdin.Stat() // get the FileInfo struct describing the standard input.
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		data, _ = io.ReadAll(os.Stdin)
	}
	//As  ccwc will be exported, so commands will run as ccwc -c,-l,-m
	// Output the number of bytes in a file
	countOfBytesPtr := flag.Bool("c", false, "Output the number of bytes in a file")
	// Output the number of lines in a file
	countOfLinesPtr := flag.Bool("l", false, "Output the number of lines in a file")
	// Output the number of words in a file
	countOfWordsPtr := flag.Bool("w", false, "Output the number of words in a file")
	// Output the number of characters in a file
	countOfCharactersPtr := flag.Bool("m", false, "Output the number of characters in a file")
	flag.Parse()
	var fileName string
	if len(data) == 0 {
		fileName = flag.Args()[0]
		var err error
		//read from a file
		data, err = os.ReadFile(fileName)
		if err != nil {
			panic("no file present ")
		}
	}
	text := string(data)

	if *countOfBytesPtr {
		fmt.Printf("%v %s\n", countBytes(text), fileName)
	} else if *countOfLinesPtr {
		fmt.Printf("%v %s\n", countLines(text), fileName)
	} else if *countOfWordsPtr {
		fmt.Printf("%v %s\n", countWords(text), fileName)
	} else if *countOfCharactersPtr {
		fmt.Printf("%v %s\n", countChars(text), fileName)
	} else {
		fmt.Printf("%v %v %v %v\n", countLines(text), countWords(text), countBytes(text), fileName)
	}

}

func countBytes(text string) int {
	return len([]byte(strings.TrimSpace(text)))
}

func countWords(text string) int {
	return len(strings.Fields(text))
}

func countLines(text string) int {
	return len(strings.Split(text, "\n"))
}

func countChars(text string) int {
	return utf8.RuneCountInString(text)
}

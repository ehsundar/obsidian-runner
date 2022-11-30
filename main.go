package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	mdFile := flag.String("mdfile", "", "Path to markdown file")
	flag.Parse()

	//fmt.Println(*mdFile)
	f, err := os.Open(*mdFile)

	if err != nil {
		panic(err)
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	//fmt.Println(string(contents))

	expr := regexp.MustCompile("(?s)```go(.*?)```")

	for _, segment := range expr.FindAllSubmatchIndex(contents, 100) {
		codeSegment := contents[segment[2]:segment[3]]
		fmt.Printf("%q\n", codeSegment)
		results := CalcResults(codeSegment)
		fmtResults := []byte(fmt.Sprintf("\n```shell\n%s```", string(results)))

		newContents := make([]byte, len(contents)+len(fmtResults))
		copy(newContents, contents[:segment[1]])
		copy(newContents[segment[1]:segment[1]+len(fmtResults)], fmtResults)
		copy(newContents[segment[1]+len(fmtResults):], contents[segment[1]:])

		err = os.WriteFile(*mdFile, newContents, 0644)
		if err != nil {
			panic(err)
		}
		break
	}
}

func CalcResults(code []byte) []byte {
	err := os.WriteFile("prog.go", code, 0644)
	if err != nil {
		panic(err)
	}

	defer func() {
		os.Remove("prog.go")
	}()

	cmd := exec.Command("go", "run", "prog.go")
	result, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s\n", err)
		return ([]byte(err.Error()))
	}
	return (result)
}

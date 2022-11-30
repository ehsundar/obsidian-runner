package main

import (
	"flag"
	"fmt"
	"golang.org/x/exp/slices"
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

	expr1 := regexp.MustCompile("(?s)```result.*?```\n")
	contents = expr1.ReplaceAll(contents, []byte(""))

	fmt.Println(string(contents))

	expr := regexp.MustCompile("(?s)```go(.*?)```")

	for _, segment := range expr.FindAllSubmatchIndex(contents, 100) {
		codeSegment := contents[segment[2]:segment[3]]
		fmt.Printf("%q\n", codeSegment)
		results := CalcResults(codeSegment)
		fmtResults := []byte(fmt.Sprintf("\n```result\n%s```", string(results)))

		contents = slices.Insert(contents, segment[1], fmtResults...)

		err = os.WriteFile(*mdFile, contents, 0644)
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
		return []byte(err.Error())
	}
	return result
}

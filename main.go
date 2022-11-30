package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	oldResultReplaceExpr = regexp.MustCompile("(?s)```result\n.*?```\n")
	goCodeExtractorExpr  = regexp.MustCompile("(?s)```go\n(.*?)```")
)

func main() {
	mdFile := flag.String("mdfile", "", "Path to markdown file")
	flag.Parse()

	f, err := os.Open(*mdFile)

	if err != nil {
		panic(err)
	}

	contents, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	contents = oldResultReplaceExpr.ReplaceAll(contents, []byte(""))

	codeSegments := goCodeExtractorExpr.FindAllSubmatchIndex(contents, 100)

	for i := 0; i < len(codeSegments); i++ {
		segment := codeSegments[i]

		codeSegment := contents[segment[2]:segment[3]]

		isComplete := CheckIsCompleteCode(codeSegment)
		if !isComplete {
			codeSegment = GetCompleteCode(codeSegment)
		}

		newContents, results := CalcResults(codeSegment)
		resultsStr := EnsureNewLineEnding(string(results))
		fmtResults := []byte(fmt.Sprintf("\n```result\n%s```", resultsStr))

		contents = slices.Insert(contents, segment[1], fmtResults...)
		contents = slices.Replace(contents, segment[2], segment[3], newContents...)

		err = os.WriteFile(*mdFile, contents, 0644)
		if err != nil {
			panic(err)
		}
		contents, err = os.ReadFile(*mdFile)
		if err != nil {
			panic(err)
		}
		codeSegments = goCodeExtractorExpr.FindAllSubmatchIndex(contents, 100)
	}
}

func CalcResults(code []byte) ([]byte, []byte) {
	err := os.WriteFile("prog.go", code, 0644)
	if err != nil {
		panic(err)
	}

	defer func() {
		os.Remove("prog.go")
	}()

	exec.Command("gofmt", "-w", "prog.go").Run()

	newContents, err := os.ReadFile("prog.go")
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("go", "run", "prog.go")
	result, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s\n", err)
		return newContents, []byte(err.Error())
	}
	return newContents, result
}

func EnsureNewLineEnding(s string) string {
	if strings.HasSuffix(s, "\n") {
		return s
	}
	return s + "\n"
}

func CheckIsCompleteCode(code []byte) bool {
	c := string(code)
	c = strings.TrimSpace(c)
	return strings.HasPrefix(c, "package main")
}

func GetCompleteCode(baseCode []byte) []byte {
	base := `
package main
import (
"fmt"
"sort"
)

func main() {
	%s
}

`
	code := fmt.Sprintf(base, baseCode)
	return []byte(code)
}

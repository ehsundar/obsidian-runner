# Obsidian Runner

`obsidian-runner` CLI helps you format and run code blocks embedded inside any valid Markdown
document and store results below the code block.

### Supported Languages

- Go

### Showcase

`obsidian-runner` compiles and runs any go code. It also `gofmt`s your code block
in-place.

```go
package main

import "fmt"

func main() {
	fmt.Println("results right below")
}
```
```result
results right below
```

This small code block is also runnable but needs `goimports` tool installed to
remove unused imports.

```go
fmt.Println("results generated by obsidian-runner!")
```
```result
results generated by obsidian-runner!
```

## Usage

Standalone:

```shell
obsidian-runner -mdfile README.md
```

With Obsidian:

1. Install [obsidian-shellcommands](https://github.com/Taitava/obsidian-shellcommands) plugin
2. Install obsidian-runner:
    ```shell
    go install github.com/ehsundar/obsidian-runner@latest
    ```
3. Add this command `PATH='$PATH:/opt/homebrew/bin:/Users/ehsan/go/bin' /Users/ehsan/go/bin/obsidian-runner -mdfile {{file_path:absolute}}`.
    you need to replace your username and maybe change path based on your OS.
4. Use `cmd+p` to run obsidian commands. You may run your shell-command from there.

> [!warning]
> Do not use shell-command events to run `obsidian-runner` automatically. May
> cause inconsistency in the contents of the final result MD file.
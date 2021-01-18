package main

import (
	"os"
	"strings"
	"time"

	"github.com/zaydek/retro/color"
	"github.com/zaydek/retro/logger"
)

var spaces = strings.Repeat(" ", 2)

// TODO: Add support for duration.
var stdout = logger.New(os.Stdout, func(msg string) string {
	t := time.Now()

	arr := strings.Split(msg, "\n")
	for x := range arr {
		if arr[x] != "" {
			if x == 0 {
				arr[x] = spaces + color.Bold(arr[x])
				continue
			}
			arr[x] = spaces + arr[x]
		}
	}

	dur := time.Since(t) // FIXME
	transformed := "\n" +
		strings.Join(arr, "\n") +
		// "\n" + fmt.Sprintf("⚡️ %0.3fs", dur.Seconds()) +
		"\n"
	return transformed
})

// TODO: Add support for duration.
var stderr = logger.New(os.Stdout, func(msg string) string {
	t := time.Now()

	arr := strings.Split(msg, "\n")
	for x := range arr {
		if arr[x] != "" {
			if x == 0 {
				arr[x] = spaces + color.BoldRed("error: ") + color.Bold(arr[x])
				continue
			}
			arr[x] = spaces + strings.Repeat(" ", len("error: ")) + arr[x]
		}
	}

	dur := time.Since(t)
	transformed := "\n" +
		strings.Join(arr, "\n") +
		// "\n" + fmt.Sprintf("⚡️ %0.3fs", dur.Seconds()) +
		"\n"
	return transformed
})

func main() {
	// defer color.TerminateFormatting()
	// fmt.Println(color.BoldWhitef("Hello, world! %s", "haha"))

	// stdout.Println("Hello, world!")
	// stderr.Println("Hello, world!\nOops\n\nI’m not sure this a good idea")
	// stdout.Println("Hello, world!")
	stderr.Println("Hello, world!")

	// err := os.MkdirAll("a/hello.go", 0755)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("path=%s\n", path.Dir("a/b/c"))
	// fmt.Printf("path=%s\n", path.Base("a/b/c.go"))

	// err := ioutil.WriteFile("a/b/c", []byte(""), 0644)
	// if err != nil {
	// 	panic(err)
	// }
}

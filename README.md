# go 限时并发
``` bash
go get -u github.com/cat3306/goasynctask
```
## example
``` go
package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cat3306/goasynctask"
)

func netIOTask(uri string) ([]byte, error) {
	rsp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	raw, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
func main() {
	now := time.Now()
	task := goasynctask.New[[]byte]()
	task.AddWithKey("baidu", func() ([]byte, error) {
		return netIOTask("https://www.baidu.com")
	})
	task.AddWithKey("bilibili", func() ([]byte, error) {
		return netIOTask("https://www.bilibili.com")
	})
	task.AddWithKey("godev", func() ([]byte, error) {
		return netIOTask("https://go.dev")
	})
	task.AddWithKey("github", func() ([]byte, error) {
		return netIOTask("https://github.com")
	})
	task.AddWithKey("blog", func() ([]byte, error) {
		return netIOTask("https://cat3306.github.io")
	})
	err := task.Run(10 * time.Second)
	if err != nil {
		panic(err)
	}
	m := task.Result()
	for k, v := range m {
		fmt.Printf("%s,byte len %d,cost:%dms\n", k, len(v.Result), v.Cost)
	}
	fmt.Printf("total cost:%d\n", time.Since(now).Milliseconds())
}
```
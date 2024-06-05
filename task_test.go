package goasynctask

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func Test(t *testing.T) {
	task := New[[]byte]()
	task.AddWithKey("baidu", func() ([]byte, error) {
		rsp, err := http.Get("https://www.baidu.com")
		if err != nil {
			return nil, err
		}
		raw, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		return raw, nil
	})
	task.AddWithKey("bilibili", func() ([]byte, error) {
		rsp, err := http.Get("https://www.bilibili.com")
		if err != nil {
			return nil, err
		}
		raw, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		return raw, nil
	})
	task.AddWithKey("godev", func() ([]byte, error) {
		rsp, err := http.Get("https://go.dev")
		if err != nil {
			return nil, err
		}
		raw, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		return raw, nil
	})
	task.AddWithKey("github", func() ([]byte, error) {
		rsp, err := http.Get("https://github.com")
		if err != nil {
			return nil, err
		}
		raw, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		return raw, nil
	})
	task.AddWithKey("blog", func() ([]byte, error) {
		rsp, err := http.Get("https://cat3306.github.io")
		if err != nil {
			return nil, err
		}
		raw, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		return raw, nil
	})
	err := task.Run(time.Second)
	if err != nil {
		t.Fatal(err)

	}
	m := task.Result()
	for k, v := range m {
		t.Logf("%s,byte len %d,cost:%dms", k, len(v.Result), v.Cost)
	}
}

func TestTimeOut(t *testing.T) {
	task := New[any]()
	task.Add(func() (any, error) {
		return 1, nil
	})
	task.Add(func() (any, error) {
		return 2, nil
	})
	task.Add(func() (any, error) {
		time.Sleep(time.Second * 2)
		return 3, nil
	})
	err := task.Run(time.Second)
	if err != nil {
		t.Logf(err.Error())
	}
	time.Sleep(time.Second * 5)
}

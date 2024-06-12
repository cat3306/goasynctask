package goasynctask

import (
	"io"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/panjf2000/ants"
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
func Test(t *testing.T) {
	now := time.Now()
	task := New[[]byte]()
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
	err := task.Run(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	m := task.ResultSet()
	for k, v := range m {
		t.Logf("%s,byte len %d,cost:%dms", k, len(v.Result), v.Cost)
	}
	t.Logf("total cost:%d", time.Since(now).Milliseconds())
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

func TestWithGoPool(t *testing.T) {
	now := time.Now()
	p, err := ants.NewPool(runtime.NumCPU())
	if err != nil {
		t.Fatal(err)
	}
	task := NewWithPool[[]byte](p)
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
	err = task.Run(time.Second)
	if err != nil {
		t.Fatal(err)

	}
	m := task.ResultSet()
	for k, v := range m {
		t.Logf("%s,byte len %d,cost:%dms", k, len(v.Result), v.Cost)
	}
	t.Logf("total cost:%d", time.Since(now).Milliseconds())
}

func TestWithoutTask(t *testing.T) {
	now := time.Now()
	raw0, err := netIOTask("https://www.baidu.com")
	if err != nil {
		t.Fatal(err)
	}
	raw1, err := netIOTask("https://www.bilibili.com")
	if err != nil {
		t.Fatal(err)
	}
	raw2, err := netIOTask("https://go.dev")
	if err != nil {
		t.Fatal(err)
	}
	raw3, err := netIOTask("https://github.com")
	if err != nil {
		t.Fatal(err)
	}
	raw4, err := netIOTask("https://cat3306.github.io")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("baidu:%d,bilibili:%d,godev:%d,github.com:%d,blog:%d", len(raw0), len(raw1), len(raw2), len(raw3), len(raw4))
	t.Logf("total cost:%d", time.Since(now).Milliseconds())
}

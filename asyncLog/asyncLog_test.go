package asyncLog

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func Test_test(t *testing.T) {
	for i := 0; i < 10; i++ {
		CLILog("./cli.log", "hello log "+strconv.Itoa(i), Error)
	}
	time.Sleep(100 * time.Second)
}

func Test_http(t *testing.T) {
	fmt.Println("server start...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
		WebLog(r, "./web.log", "hello log", Info)
	})
	http.ListenAndServe(":8080", nil)
}

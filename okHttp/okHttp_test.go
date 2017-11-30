package okHttp

import (
	"fmt"
	"testing"
)

func Test_test(t *testing.T) {

	// var start = time.Now()
	// time.Sleep(time.Second * 2)
	// //time.Sleep(time.Millisecond * 250)
	// var stop = time.Now()
	// print(stop.Sub(start) / 1000000)
	//fmt.Println(stop.Sub(start))

	//rand.Seed(time.Now().UnixNano())
	//fmt.Println(rand.Intn(10))

	//urls := urlSlave("http://www.teelet.cn|www.teelet.com/aaa/bbb/")
	//fmt.Println(urls)

	// _, ext, err := Request(1, "https://www.teelet.cn", "", nil, 100)
	// fmt.Println(ext, err)

	res, ext, _ := Get("http://www.teelet.com", MaxTimeOut, MaxRetry, nil)
	fmt.Println(string(res), ext)

}

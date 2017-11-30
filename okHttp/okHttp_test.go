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

	//urls := urlSlave("http://www.360.cn|www.360.com/aaa/bbb/")
	//fmt.Println(urls)

	// _, ext, err := Request(1, "https://www.360.cn", "", nil, 100)
	// fmt.Println(ext, err)

	res, ext, _ := Get("http://u.api.look.360.cn/comment/lists?callback=jQuery183030720481380558473_1508990360678&type=1&num=10&sub_limit=5&start_date=&page=1&page_key=f820374d29edebdf&url=f820374d29edebdf&client_id=15&uid=91251416.3900915413553927000.1505385580562.1743&_=1508990361429", MaxTimeOut, MaxRetry, nil)
	fmt.Println(string(res), ext)

	//res1, ext, _ := Post("http://u.api.look.360.cn/comment/lists", "callback=jQuery183030720481380558473_1508990360678&type=1&num=10&sub_limit=5&start_date=&page=1&page_key=f820374d29edebdf&url=f820374d29edebdf&client_id=15&uid=91251416.3900915413553927000.1505385580562.1743&_=1508990361429", MaxTimeOut, MaxRetry, nil)
	//fmt.Println(string(res1), ext)
}

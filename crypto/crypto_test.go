package crypto

import (
	"fmt"
	"testing"
)

func Test_test(t *testing.T) {
	// addr, _ := net.Interfaces()
	// var mac bytes.Buffer
	// for _, v := range addr {
	// 	ar := fmt.Sprintf("%v", v.HardwareAddr)
	// 	mac.WriteString(ar)
	// }
	// fmt.Println(mac.String())

	res := EnidEncode("1.111", "2.222", "beijing")
	fmt.Println(res)

	res2 := EnidDecode(res)
	fmt.Println(res2)

	res3 := Sha1("xxx")
	fmt.Println(res3)
}

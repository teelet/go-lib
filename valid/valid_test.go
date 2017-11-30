package valid

import (
	"net/url"
	// "fmt"
	"testing"
)

func Test_CheckUid(t *testing.T) {
	err := CheckUid("")
	if err == nil {
		t.Error("empty test failed")
	}

	err = CheckUid("-asdkjh123WASDasdkjh.asddas___jdkh")
	if err != nil {
		t.Error("check good case failed")
	}

	err = CheckUid("asdf  ")
	if err == nil {
		t.Error("check empty char failed")
	}

	err = CheckUid("asdf$aaa")
	if err == nil {
		t.Error("check invalid char failed")
	}
}

func Test_CheckSign(t *testing.T) {
	err := CheckSign("")
	if err == nil {
		t.Error("empty test failed")
	}

	err = CheckSign("asdkjhasFFdkjh123asddas___jdkh")
	if err != nil {
		t.Error("check good case failed")
	}

	err = CheckSign("asdf  ")
	if err == nil {
		t.Error("check empty char failed")
	}

	err = CheckSign("asdf$aaa")
	if err == nil {
		t.Error("check invalid char failed")
	}
}

func Test_CheckUrls(t *testing.T) {
	err := CheckUrls([]string{})
	if err == nil {
		t.Error("empty test failed")
	}

	err = CheckUrls([]string{"http://360.cn", "http://taobao.com"})
	if err != nil {
		t.Error("check good case failed")
	}

	err = CheckUrls([]string{"tcp://360.cn", "http://taobao.com"})
	if err == nil {
		t.Error("check bad case failed")
	}
}

func Test_CheckChannel(t *testing.T) {
	err := CheckChannel("")
	if err == nil {
		t.Error("empty test failed")
	}

	err = CheckChannel("asdkjhasFFdkjh123asddas___jdkh")
	if err != nil {
		t.Error("check good case failed")
	}

	err = CheckChannel("asdf  ")
	if err == nil {
		t.Error("check empty char failed")
	}

	err = CheckChannel("asdf$aaa")
	if err == nil {
		t.Error("check invalid char failed")
	}
}

func Test_CheckNumber(t *testing.T) {
	err := CheckNumber("", "myname")
	if err == nil {
		t.Error("empty test failed")
	}

	err = CheckNumber("360", "myname")
	if err != nil {
		t.Error("check good case failed")
	}

	err = CheckNumber("3.14", "myname")
	if err == nil {
		t.Error("check dot failed")
	}

	err = CheckNumber("0", "myname")
	if err != nil {
		t.Error("check zero failed")
	}
}

func Test_CheckWord(t *testing.T) {
	err := CheckWord("", "myname")
	if err == nil {
		t.Error("empty test failed")
	}

	err = CheckWord("asdkjhasFFdkjh123asddas___jdkh", "myname")
	if err != nil {
		t.Error("check good case failed")
	}

	err = CheckWord("asdf  ", "myname")
	if err == nil {
		t.Error("check empty char failed")
	}

	err = CheckWord("asdf$aaa", "myname")
	if err == nil {
		t.Error("check invalid char failed")
	}
}

func Test_CheckEmpty(t *testing.T) {
	err := CheckEmpty("", "myname")
	if err == nil {
		t.Error("empty test failed")
	}

	err = CheckEmpty("asdkjhasFFdkjh123asddas___jdkh", "myname")
	if err != nil {
		t.Error("check good case failed")
	}

	err = CheckEmpty("  ", "myname")
	if err == nil {
		t.Error("check empty char failed")
	}
}

func Test_CheckInArray(t *testing.T) {
	err := CheckInArray("haha", []string{"haha", "heihei"}, "myname")
	if err != nil {
		t.Error("in array test failed")
	}

	err = CheckInArray("papa", []string{"haha", "heihei"}, "myname")
	if err == nil {
		t.Error("not in array failed")
	}
}

func Test_CheckRequest(t *testing.T) {
	err := CheckRequest(url.Values{"aa": {"yes"}, "_asdf": {"aloha"}})
	if err != nil {
		t.Error("check good case failed")
	}

	err = CheckRequest(url.Values{"some": {"yes"}, "_asdfG": {"kawai"}})
	if err == nil {
		t.Error("check bad case failed")
	}
}

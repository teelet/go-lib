package valid

import (
	"net/url"
	"regexp"
	"sdk.look.360.cn/application/library/badcase"
	"sdk.look.360.cn/application/library/errno"
	"sdk.look.360.cn/application/library/helper"
	"strings"
)

func CheckUid(str string) error {
	if !regexp.MustCompile(`^[-.\w]+$`).MatchString(str) {
		return badcase.New("check uid is error.uid="+helper.JsonEncode(str), errno.CHECK_UID)
	}
	return nil
}

func CheckSign(str string) error {
	if !regexp.MustCompile(`^\w+$`).MatchString(str) {
		return badcase.New("check sign is error.sign="+helper.JsonEncode(str), errno.CHECK_SIGN)
	}
	return nil
}

func CheckUrls(urls []string) error {
	if len(urls) == 0 {
		return badcase.New("check urls is empty", errno.CHECK_URLS)
	}

	reg := regexp.MustCompile(`^http://.+$`)
	for _, str := range urls {
		if !reg.MatchString(str) {
			return badcase.New("check url is error.url="+helper.JsonEncode(str), errno.CHECK_URL)
		}
	}
	return nil
}

func CheckChannel(str string) error {
	if str == "" || !regexp.MustCompile(`^\w+$`).MatchString(str) {
		return badcase.New("check channel is error.channel="+helper.JsonEncode(str), errno.CHECK_CHANNEL)
	}
	return nil
}

func CheckNumber(str string, name string) error {
	if !regexp.MustCompile(`^\d+$`).MatchString(str) {
		return badcase.New("check "+name+" is error.number="+helper.JsonEncode(str), errno.CHECK_NUMBER)
	}
	return nil
}

func CheckWord(str string, name string) error {
	if !regexp.MustCompile(`^\w+$`).MatchString(str) {
		return badcase.New("check "+name+" is error.word="+helper.JsonEncode(str), errno.CHECK_WORD)
	}
	return nil
}

func CheckCallback(str string) error {
	// if !regexp.MustCompile(`^[\w_]+$`).MatchString(str) {
	//     return badcase.New("check callback is error.uid=" + helper.JsonEncode(str), errno.CHECK_WORD)
	// }
	return nil
}

func CheckEmpty(str string, name string) error {
	if strings.TrimSpace(str) == "" {
		return badcase.New("check "+name+" is empty", errno.CHECK_EMPTY)
	}
	return nil
}

func CheckInArray(str string, arr []string, name string) error {
	for _, v := range arr {
		if v == str {
			return nil
		}
	}
	return badcase.New("check "+name+" is not in array", errno.CHECK_EMPTY)
}

func CheckRequest(urls url.Values) error {
	for key, _ := range urls {
		if key[0] == '_' && len(key) > 5 {
			return badcase.New("check request is error.key="+helper.JsonEncode(key), errno.CHECK_REQUEST)
		}
	}
	return nil
}

package alarm

import (
	"fmt"
	"testing"
)

func Test_test(t *testing.T) {
	res := DingDingRobot(
		"https://oapi.dingtalk.com/robot/send?access_token=...",
		"keyword",
		"content content",
	)

	fmt.Println(res)
}


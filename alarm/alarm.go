package alarm

import (
	"fmt"

	"gitlab.bitmain.com/bitdeer/go-lib/okHttp"
)

func DingDingRobot(webHookUrl, keyword, content string) string {
	header := map[string]string{"Content-Type": "application/json"}
	msgFmt := `{
			"msgtype": "text", 
        	"text": {
				"content": "[%s]: %s"
			}
		}`
	content = fmt.Sprintf(msgFmt, keyword, content)
	res, _, _ := okHttp.Post(webHookUrl, content, 500, 0, header)

	return string(res)
}


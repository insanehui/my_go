package ytx

import (
	"testing"
	"time"
)

// 消息推送使用示例
func TestPush0(t *testing.T) {

	data := make(map[string]interface{})
	data["from"] = "0"
	data["to"] = []string{"456", "888"}
	data["msg"] = "测试消息：" + time.Now().Format("15点4分5秒")
	// 以下是自定义字段
	data["whoami"] = "服务端推送"
	data["bilibala"] = "hahaha"

	Push(data)
}

// 测试发送自定义消息类型
func TestPushCustom(t *testing.T) {

	data := make(map[string]interface{})
	data["_type"] = "NewPic"
	data["to"] = "456"
	// 以下是自定义字段
	data["path_id"] = "这是路径id"
	data["pic_url"] = "这是图片的url"

	Push(data)
}

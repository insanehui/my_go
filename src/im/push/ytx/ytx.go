package ytx

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"utils"
	j "utils/json"

	"im/push"
)

// # 推送消息
// 需要的字段：
// _type: 消息类型。( = "ChatMsg" )
// from : 发送方。( = "0" )，id可为字符串，或者数字
// to : 接收方。可以是 "123" 或者 ["123", "456"]，这里的id必须为字符串（暂时还没做智能识别 - -）
// msg : 如果 _type 为 ChatMsg，则不能缺省。否则可以缺省。但 _type 需要指定为其他类型（约定只能使用文档中定义的类型）
// 其他字段：取决于_type对应的类型
// 示例：
/*
	data := make(map[string]interface{}) // 先构造一个map
	data["_type"] = "NewPic" // 必填！填上与客户端约定的消息类型
	data["to"] = "456" // 必填！填上接收方用户id
	// 以下是消息类型对应的字段
	data["path_id"] = "这是路径id"
	data["pic_url"] = "这是图片的url"
	Push(data) // 发出消息。完整形式是：ytx.Push(data)
*/
func Push(data map[string]interface{}) (bool, string) {

	ok := true
	msg := ""

	for {

		// 生成sig
		account_sid := "8a216da85600ef240156070dc9cd0710"
		auth_token := "d798935e20c547969d88545b26ec3d7b"
		app_id := "8a216da85600ef240156070dca340716"

		ts := time.Now().Format("20060102150405") // string类型，记忆法：1月2日下午3点4分5秒06年
		sig_raw := account_sid + auth_token + ts
		sig := utils.MD5(sig_raw) // 大写MD5
		log.Println("sig:", sig)

		// 生成Authorization
		auth := base64.URLEncoding.EncodeToString([]byte(account_sid + ":" + ts))
		log.Println("auth:", auth)

		// 发起http post请求
		msg_push_url := "https://sandboxapp.cloopen.com:8883/2013-12-26/Accounts/" + account_sid + "/IM/PushMsg?sig=" + sig
		// msg_push_url := "https://app.cloopen.com:8883/2013-12-26/Accounts/" + account_sid + "/IM/PushMsg?sig=" + sig
		client := &http.Client{}
		//# 构建json内容（包体）

		c := make(map[string]interface{})
		c["pushType"] = "1"
		c["appId"] = app_id
		if from, ok := data["from"]; ok {
			c["sender"] = from
		} else {
			c["sender"] = "0"
		}

		var to interface{}

		if _, ok := data["to"].(string); ok { // 如果是字符串
			// 转成数组
			to = []interface{}{data["to"]}
		} else {
			to = data["to"]
		}
		c["receiver"] = to
		c["msgType"] = "1"

		// 过滤出一个jmsg
		//		jmsg := utils.JSON(data).OmitNew("from to time file_type ses_id is_offline url url_small file_path file_name file_size voice_length _progress _return")
		jmsg := j.OmitNew(data, "from to time file_type ses_id is_offline url url_small file_path file_name file_size voice_length _progress _return")
		log.Printf("省略之后的jmsg: %+v", jmsg)

		// 对jmsg进行编码
		var bmsg []byte
		bmsg, _ = json.Marshal(jmsg)
		c["msgContent"] = msg_encode(string(bmsg))
		//	c["msgContent"] = data["msg"]

		c_buf, _ := json.MarshalIndent(c, "", "    ")
		c_str := string(c_buf)
		log.Println("content:", c_str)

		req, _ := http.NewRequest("POST", msg_push_url, strings.NewReader(c_str))
		req.Header.Set("Accept", "application/json") // 设置http头
		req.Header.Set("Content-Type", "application/json;charset=utf-8")
		req.Header.Set("Authorization", auth)
		res, _ := client.Do(req)            // 真正发起请求
		body, _ := ioutil.ReadAll(res.Body) // 取得返回结果

		// 返回结果形如：
		// {"statusMsg":"success","statusCode":"000000"}
		jret := make(Data)
		json.Unmarshal(body, &jret)

		log.Printf("YTX服务端返回结果：%+v", jret)
		if jret["statusCode"] != "000000" {
			ok = false
		}
		msg = string(body)

		break
	}

	return ok, msg
}

type Data push.Data

// 由于服务端推送对json格式的消息进行了限制，故进行“编码”
func msg_encode(str string) string {
	return "x" + str
}

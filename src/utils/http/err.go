package http

// 错误类型
// 使用方法：
// JsonDo函数中 定义ret结构的时候，可以用来组合
// var ret struct {
// 	Err // 这里套一个匿名的本 Err类型
// 	Data interface{} `json:"data"` // 这里存放真正业务返回数据
// }

type Err struct {
	Code int    `json:"code"` // 返回码
	Msg  string `json:"msg"`  // 错误消息
}

func (me *Err) Ok() bool {
	return me.Code == 0
}

func (me *Err) FromStr(s string) {
	me.Msg = s;
	me.Code = -1; // 缺省设为-1吧
}

func (me *Err) FromError(e error) {
	me.Msg = e.Error();
	me.Code = -1; // 缺省设为-1吧
}

func (me *Err) FromPanic(p interface{}) {
	// 先支持 error 格式
	if e, ok := p.(error); ok {
		me.FromError(e)
	} else if e, ok := p.(string); ok {
		me.FromStr(e)
	}
}


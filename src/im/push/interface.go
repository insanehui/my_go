package push

/*
参数：同ChatMsg的结构，to除了支持字符串之外，还可以支持数组
*/
type Data map[string]interface{} // 使用ytx的接口时，居然不能使用push.Data，而要使用ytx.Data...

type I interface {
	PushMsg(data Data) (bool, string) // 使用json作为参数，返回ok, msg
}

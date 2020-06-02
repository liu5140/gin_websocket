# gin_websocket
由于在项目过程中，连续做了好几个项目，都要做消息列表,原来一直都是基于一个拉模式，即前端定时轮询接口，查询后台数据，对于这种，
有点浪费资源，并且很多需要及时更新得弹出，他不一定能做到，所以心血来潮，想换个长链接玩一下

websocket 主要是实现三个方法，一个read一个write,还有一个close  前面两个非线程安全，其中close是一个可重入的，可以多次进行调用

场景：
该例子适用场景
1、单个用户的通知,例如该用户充值成功，需要弹出一个提示
2、多个用户的通知,例如系统维护公告

具体写入数据调用方法：
func WriteMessage(userid int64, data []byte) (err error) {
	if syncma, ok := concurrentHolder.Load(userid); ok {
		ws := syncma.(*WebSocketService)
		ws.outChan <- data
	} else {
		Log.Infoln("================没有建立连接=============")
	}
	return nil
}

如何保证read,write的线程安全

增加了两个chan chan是线程安全，所以我们重写read和write也是线程安全的


题外话：如果有人问你golang一分钟处理十万，或者百万访问量，这个时候人家想听到的，应该是你使用go的chan 一边往chan里面写东西，一边从chan里面读取东西，
做你接下来的业务


package wechat

const EVENT_SUBSCRIBE = "subscribe"     //订阅
const EVENT_UNSUBSCRIBE = "unsubscribe" //取消订阅
const EVENT_SCAN = "SCAN"               //扫描带参数二维码
const EVENT_LOCATION = "LOCATION"       //上报地理位置

// Message数据接口
// http://mp.weixin.qq.com/wiki/17/f298879f8fb29ab98b2f2971d42552fd.html
type MessageMsg struct {
	MsgId        int64
	Content      Cdata
	PicUrl       Cdata
	MediaId      Cdata
	Format       Cdata
	ThumbMediaId Cdata
	Location_X   float64
	Location_Y   float64
	Scale        float64
	Label        Cdata
	Title        Cdata
	Description  Cdata
	Url          Cdata

	Ticket    Cdata
	Latitude  float64
	Longitude float64
	Precision float64
}

// 文本消息
type MessageText struct {
	MsgId   int64
	Content Cdata
}

func (ctx *Context) GetMessageText() (res MessageText, ok bool) {
	if ctx.GetMsgType() != MSGTYPE_TEXT {
		return
	}
	ok = true
	res = MessageText{
		MsgId:   ctx.WXMsg.MsgId,
		Content: ctx.WXMsg.Content,
	}
	return
}

// 图片消息
type MessageImage struct {
	MsgId   int64
	PicUrl  Cdata
	MediaId Cdata
}

func (ctx *Context) GetMessageImage() (res MessageImage, ok bool) {
	if ctx.GetMsgType() != MSGTYPE_IMAGE {
		return
	}
	ok = true
	res = MessageImage{
		MsgId:   ctx.WXMsg.MsgId,
		PicUrl:  ctx.WXMsg.PicUrl,
		MediaId: ctx.WXMsg.MediaId,
	}
	return
}

// 语音消息
type MessageVoice struct {
	MsgId   int64
	Format  Cdata
	MediaId Cdata
}

func (ctx *Context) GetMessageVoice() (res MessageVoice, ok bool) {
	if ctx.GetMsgType() != MSGTYPE_VOICE {
		return
	}
	ok = true
	res = MessageVoice{
		MsgId:   ctx.WXMsg.MsgId,
		Format:  ctx.WXMsg.Format,
		MediaId: ctx.WXMsg.MediaId,
	}
	return
}

// 视频消息
type MessageVideo struct {
	MsgId        int64
	MediaId      Cdata
	ThumbMediaId Cdata
}

func (ctx *Context) GetMessageVideo() (res MessageVideo, ok bool) {
	if ctx.GetMsgType() != MSGTYPE_VIDEO {
		return
	}
	ok = true
	res = MessageVideo{
		MsgId:        ctx.WXMsg.MsgId,
		MediaId:      ctx.WXMsg.MediaId,
		ThumbMediaId: ctx.WXMsg.ThumbMediaId,
	}
	return
}

// 小视频消息
func (ctx *Context) GetMessageShortVideo() (res MessageVideo, ok bool) {
	if ctx.GetMsgType() != MSGTYPE_SHORTVIDEO {
		return
	}
	ok = true
	res = MessageVideo{
		MsgId:        ctx.WXMsg.MsgId,
		MediaId:      ctx.WXMsg.MediaId,
		ThumbMediaId: ctx.WXMsg.ThumbMediaId,
	}
	return
}

// 地理位置消息
type MessageLocation struct {
	MsgId      int64
	Location_X float64
	Location_Y float64
	Scale      float64
	Label      Cdata
}

func (ctx *Context) GetMessageLocation() (res MessageLocation, ok bool) {
	if ctx.GetMsgType() != MSGTYPE_LOCATION {
		return
	}
	ok = true
	res = MessageLocation{
		MsgId:      ctx.WXMsg.MsgId,
		Location_X: ctx.WXMsg.Location_X,
		Location_Y: ctx.WXMsg.Location_Y,
		Scale:      ctx.WXMsg.Scale,
		Label:      ctx.WXMsg.Label,
	}
	return
}

// 链接消息
type MessageLink struct {
	MsgId       int64
	Title       Cdata
	Description Cdata
	Url         Cdata
}

func (ctx *Context) GetMessageLink() (res MessageLink, ok bool) {
	if ctx.GetMsgType() != MSGTYPE_LINK {
		return
	}
	ok = true
	res = MessageLink{
		MsgId:       ctx.WXMsg.MsgId,
		Title:       ctx.WXMsg.Title,
		Description: ctx.WXMsg.Description,
		Url:         ctx.WXMsg.Url,
	}
	return
}

// 接收事件推送
// http://mp.weixin.qq.com/wiki/7/9f89d962eba4c5924ed95b513ba69d9b.html

// 关注/取消事件, 包括点击关注和扫描二维码(公众号二维码和公众号带参数二维码)关注
type MessageEvent struct {
	Event Cdata
	// 下面两个字段只有在扫描带参数二维码进行关注时才有值, 否则为空值
	EventKey Cdata // 事件KEY值, 格式为: qrscene_二维码的参数值
	Ticket   Cdata // 二维码的ticket, 可用来换取二维码图片
}

func (ctx *Context) GetMessageSubscribeEvent() (res MessageEvent, ok bool) {
	if ctx.GetMsgEvent() != EVENT_SUBSCRIBE && ctx.GetMsgEvent() != EVENT_UNSUBSCRIBE {
		return
	}
	ok = true
	res = MessageEvent{
		Event:    ctx.WXMsg.Event,
		EventKey: ctx.WXMsg.EventKey,
		Ticket:   ctx.WXMsg.Ticket,
	}
	return
}

func (ctx *Context) GetMessageScanEvent() (res MessageEvent, ok bool) {
	if ctx.GetMsgEvent() != EVENT_SCAN {
		return
	}
	ok = true
	res = MessageEvent{
		Event:    ctx.WXMsg.Event,
		EventKey: ctx.WXMsg.EventKey,
		Ticket:   ctx.WXMsg.Ticket,
	}
	return
}

// 上报地理位置事件
type MessageLocationEvent struct {
	Event     Cdata   // LOCATION
	Latitude  float64 // 地理位置纬度
	Longitude float64 // 地理位置经度
	Precision float64 // 地理位置精度(整数? 但是微信推送过来是浮点数形式)
}

func (ctx *Context) GetMessageLocationEvent() (res MessageLocationEvent, ok bool) {
	if ctx.GetMsgEvent() != EVENT_LOCATION {
		return
	}
	ok = true
	res = MessageLocationEvent{
		Event:     ctx.WXMsg.Event,
		Latitude:  ctx.WXMsg.Latitude,
		Longitude: ctx.WXMsg.Longitude,
		Precision: ctx.WXMsg.Precision,
	}
	return
}

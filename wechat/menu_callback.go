package wechat

const EVENT_MENU_VIEW = "VIEW"   //菜单 - 点击菜单跳转链接
const EVENT_MENU_CLICK = "CLICK" //菜单 - 点击菜单拉取消息
// 请注意, 下面的事件仅支持微信iPhone5.4.1以上版本, 和Android5.4以上版本的微信用户,
// 旧版本微信用户点击后将没有回应, 开发者也不能正常接收到事件推送.
const EVENT_MENU_SCAN_PUSH = "scancode_push"       //菜单 - 扫码推事件(客户端跳URL)
const EVENT_MENU_SCAN_WAITMSG = "scancode_waitmsg" //菜单 - 扫码推事件(客户端不跳URL)
const EVENT_MENU_PIC_SYS = "pic_sysphoto"          //菜单 - 弹出系统拍照发图
const EVENT_MENU_PIC_PHOTO = "pic_photo_or_album"  //菜单 - 弹出拍照或者相册发图
const EVENT_MENU_PIC_WEIXIN = "pic_weixin"         //菜单 - 弹出微信相册发图器
const EVENT_MENU_LOCATION = "location_select"      //菜单 - 弹出地理位置选择器

// menu数据接口
type MenuMsg struct {
	MenuId           int64            `xml:"MenuId" json:"MenuId"` // 菜单ID，如果是个性化菜单，则可以通过这个字段，知道是哪个规则的菜单被点击了
	ScanCodeInfo     scanCodeInfo     `xml:"ScanCodeInfo,omitempty" json:"ScanCodeInfo,omitempty"`
	SendPicsInfo     sendPicsInfo     `xml:"SendPicsInfo,omitempty" json:"SendPicsInfo,omitempty"`
	SendLocationInfo sendLocationInfo `xml:"SendLocationInfo,omitempty" json:"SendLocationInfo,omitempty"`
}

type scanCodeInfo struct {
	ScanType   string `xml:"ScanType"   json:"ScanType"`   // 扫描类型, 一般是qrcode
	ScanResult string `xml:"ScanResult" json:"ScanResult"` // 扫描结果, 即二维码对应的字符串信息
}

type sendPicsInfo struct {
	Count   int `xml:"Count" json:"Count"`
	PicList []struct {
		PicMd5Sum string `xml:"PicMd5Sum" json:"PicMd5Sum"`
	} `xml:"PicList>item,omitempty" json:"PicList,omitempty"`
}

type sendLocationInfo struct {
	LocationX float64 `xml:"Location_X" json:"Location_X"` // 地理位置纬度
	LocationY float64 `xml:"Location_Y" json:"Location_Y"` // 地理位置经度
	Scale     int     `xml:"Scale"      json:"Scale"`      // 精度, 可理解为精度或者比例尺, 越精细的话 scale越高
	Label     string  `xml:"Label"      json:"Label"`      // 地理位置的字符串信息
	PoiName   string  `xml:"Poiname"    json:"Poiname"`    // 朋友圈POI的名字, 可能为空
}

func (ctx *Context) checkEvent(event string) bool {
	if ctx.GetMsgType() == MSGTYPE_EVENT && ctx.GetMsgEvent() == event {
		return true
	}
	return false
}

// CLICK: 点击菜单拉取消息时的事件推送
type ClickEvent struct {
	Event    Cdata // 事件类型, CLICK
	EventKey Cdata // 事件KEY值, 与自定义菜单接口中KEY值对应
}

func (ctx *Context) GetClickEvent() (ClickEvent, bool) {
	if !ctx.checkEvent(EVENT_MENU_CLICK) {
		return ClickEvent{}, false
	}
	return ClickEvent{
		Event:    ctx.WXMsg.Event,
		EventKey: ctx.WXMsg.EventKey,
	}, true
}

// VIEW: 点击菜单跳转链接时的事件推送
type ViewEvent struct {
	Event    Cdata
	EventKey Cdata
	MenuId   int64 // 菜单ID，如果是个性化菜单，则可以通过这个字段，知道是哪个规则的菜单被点击了
}

func (ctx *Context) GetViewEvent() (ViewEvent, bool) {
	if !ctx.checkEvent(EVENT_MENU_VIEW) {
		return ViewEvent{}, false
	}
	return ViewEvent{
		Event:    ctx.WXMsg.Event,
		EventKey: ctx.WXMsg.EventKey,
		MenuId:   ctx.WXMsg.MenuId,
	}, true
}

// scancode_push: 扫码推事件的事件推送
type ScanCodePushEvent struct {
	Event        Cdata
	EventKey     Cdata
	ScanCodeInfo scanCodeInfo
}

func (ctx *Context) GetScanCodePushEvent() (ScanCodePushEvent, bool) {
	if !ctx.checkEvent(EVENT_MENU_SCAN_PUSH) {
		return ScanCodePushEvent{}, false
	}
	return ScanCodePushEvent{
		Event:        ctx.WXMsg.Event,
		EventKey:     ctx.WXMsg.EventKey,
		ScanCodeInfo: ctx.WXMsg.ScanCodeInfo,
	}, true
}

// scancode_waitmsg: 扫码推事件且弹出"消息接收中"提示框的事件推送
func (ctx *Context) GetScanCodeWaitMsgEvent() (ScanCodePushEvent, bool) {
	if !ctx.checkEvent(EVENT_MENU_SCAN_WAITMSG) {
		return ScanCodePushEvent{}, false
	}
	return ScanCodePushEvent{
		Event:        ctx.WXMsg.Event,
		EventKey:     ctx.WXMsg.EventKey,
		ScanCodeInfo: ctx.WXMsg.ScanCodeInfo,
	}, true
}

// pic_sysphoto: 弹出系统拍照发图的事件推送
type PicSysPhotoEvent struct {
	Event        Cdata
	EventKey     Cdata
	SendPicsInfo sendPicsInfo
}

func (ctx *Context) GetPicSysPhotoEvent() (PicSysPhotoEvent, bool) {
	if !ctx.checkEvent(EVENT_MENU_PIC_SYS) {
		return PicSysPhotoEvent{}, false
	}
	return PicSysPhotoEvent{
		Event:        ctx.WXMsg.Event,
		EventKey:     ctx.WXMsg.EventKey,
		SendPicsInfo: ctx.WXMsg.SendPicsInfo,
	}, true
}

// pic_photo_or_album: 弹出拍照或者相册发图的事件推送
func (ctx *Context) GetPicPhotoOrAlbumEvent() (PicSysPhotoEvent, bool) {
	if !ctx.checkEvent(EVENT_MENU_PIC_PHOTO) {
		return PicSysPhotoEvent{}, false
	}
	return PicSysPhotoEvent{
		Event:        ctx.WXMsg.Event,
		EventKey:     ctx.WXMsg.EventKey,
		SendPicsInfo: ctx.WXMsg.SendPicsInfo,
	}, true
}

// pic_weixin: 弹出微信相册发图器的事件推送
func (ctx *Context) GetPicWeixinEvent() (PicSysPhotoEvent, bool) {
	if !ctx.checkEvent(EVENT_MENU_PIC_WEIXIN) {
		return PicSysPhotoEvent{}, false
	}
	return PicSysPhotoEvent{
		Event:        ctx.WXMsg.Event,
		EventKey:     ctx.WXMsg.EventKey,
		SendPicsInfo: ctx.WXMsg.SendPicsInfo,
	}, true
}

// location_select: 弹出地理位置选择器的事件推送
type LocationSelectEvent struct {
	Event            Cdata
	EventKey         Cdata
	SendLocationInfo sendLocationInfo
}

func (ctx *Context) GetLocationSelectEvent() (LocationSelectEvent, bool) {
	if !ctx.checkEvent(EVENT_MENU_LOCATION) {
		return LocationSelectEvent{}, false
	}
	return LocationSelectEvent{
		Event:            ctx.WXMsg.Event,
		EventKey:         ctx.WXMsg.EventKey,
		SendLocationInfo: ctx.WXMsg.SendLocationInfo,
	}, true
}

package main

import (
	"fmt"
	"github.com/bjdgyc/weixin/wechat"
	"net/http"
)

func main() {
	wx := wechat.NewWechat("bjdgyc", "wx90be1dc3ca5b7e40", "c1da30ad71195b812e0facbe640951a3")
	//ip,_ := wx.GetCallbackIP()
	//fmt.Println(ip)
	//bu := wechat.Button{}
	//bu.SetAsClickButton("士大夫","ddddd")
	//menu := wechat.Menu{
	//	Buttons:[]wechat.Button{bu},
	//}

	//err := wx.CreateMenu(menu)
	//fmt.Println(err)
	//wx.GetMenu()

	//wx.SetEncodingAESKey("jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C")
	//ctx := wx.NewContext(nil,nil)
	//msgEncrypt := "RypEvHKD8QQKFhvQ6QleEB4J58tiPdvo+rtK1I9qca6aM/wvqnLSV5zEPeusUiX5L5X/0lWfrf0QADHHhGd3QczcdCUpj911L3vg3W/sYYvuJTs3TUUkSUXxaccAS0qhxchrRYt66wiSpGLYL42aM6A8dTT+6k4aSknmPj48kzJs8qLjvd4Xgpue06DOdnLxAUHzM6+kDZ+HMZfJYuR+LtwGc2hgf5gsijff0ekUNXZiqATP7PF5mZxZ3Izoun1s4zG4LUMnvw2r+KqCKIw+3IQH03v+BCA9nMELNqbSf6tiWSrXJB3LAVGUcallcrw8V2t9EL4EhzJWrQUax5wLVMNS0+rUPA3k22Ncx4XXZS9o0MBH27Bo6BpNelZpS+/uh9KsNlY6bHCmJU9p8g7m3fVKn28H3KDYA5Pl/T8Z1ptDAVe0lXdQ2YoyyH2uyPIGHBZZIs2pDBS8R07+qN+E7Q=="

	//a,err := ctx.Decrypt(msgEncrypt)
	//fmt.Println(a , err)

	//ctx.ResponseText("ddd")
	//ctx.ResponseTransferCustomerService("")

	//w := []wechat.RArticle{
	//	wechat.SetRArticle("sss","dd","wefsd","dsf"),
	//	wechat.SetRArticle("s2ss","2dd","2wefsd","2dsf"),
	//}
	//
	//ctx.ResponseNews(w)

	http.HandleFunc("/", wx.CreateHandler(weixinHandler))
	http.ListenAndServe(":8090", nil)

}

func weixinHandler(ctx *wechat.Context, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}

	switch ctx.GetMsgType() {
	case wechat.MSGTYPE_TEXT:
		fmt.Println(ctx.WXMsg.Content)
		ctx.ResponseText("回复: " + ctx.WXMsg.Content.String())
	case wechat.MSGTYPE_EVENT:
		fmt.Println(ctx.GetMsgEvent())
		if ctx.GetMsgEvent() == wechat.EVENT_MENU_CLICK {
			e, _ := ctx.GetClickEvent()
			fmt.Println("click事件", e)
			ctx.ResponseText("haha你吗松岛枫算法大的")

		}
	default:
		fmt.Println("default")
		a, _ := ctx.GetMessageVoice()
		fmt.Println(a)
	}

}

# weixin

## Introduction
weixin-golang-sdk 微信golang工具包


## Installation
`go get github.com/bjdgyc/weixin`


## example

```
package main

import (
	"fmt"
	"github.com/bjdgyc/weixin/wechat"
	"net/http"
)



func main() {
	wx := wechat.NewWechat("bjdgyc", "wx90be1dc3ca5b7e40", "c1da30ad71195b812e0facbe640951a3")

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
```
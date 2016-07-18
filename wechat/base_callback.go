package wechat

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const MSGTYPE_TEXT = "text"
const MSGTYPE_IMAGE = "image"
const MSGTYPE_LOCATION = "location"
const MSGTYPE_LINK = "link"
const MSGTYPE_EVENT = "event"
const MSGTYPE_MUSIC = "music"
const MSGTYPE_NEWS = "news"
const MSGTYPE_VOICE = "voice"
const MSGTYPE_VIDEO = "video"
const MSGTYPE_SHORTVIDEO = "shortvideo"
const MSGTYPE_DUOKEFU = "transfer_customer_service"


const EVENT_SEND_MASS = "MASSSENDJOBFINISH"         //发送结果 - 高级群发完成
const EVENT_SEND_TEMPLATE = "TEMPLATESENDJOBFINISH" //发送结果 - 模板消息发送结果
const EVENT_KF_SEESION_CREATE = "kfcreatesession"   //多客服 - 接入会话
const EVENT_KF_SEESION_CLOSE = "kfclosesession"     //多客服 - 关闭会话
const EVENT_KF_SEESION_SWITCH = "kfswitchsession"   //多客服 - 转接会话
const EVENT_CARD_PASS = "card_pass_check"           //卡券 - 审核通过
const EVENT_CARD_NOTPASS = "card_not_pass_check"    //卡券 - 审核未通过
const EVENT_CARD_USER_GET = "user_get_card"         //卡券 - 用户领取卡券
const EVENT_CARD_USER_DEL = "user_del_card"         //卡券 - 用户删除卡券
const EVENT_MERCHANT_ORDER = "merchant_order"       //微信小店 - 订单付款通知

type Cdata struct {
	Cdata string `xml:",cdata"`
}

func (c *Cdata) String() string {
	return c.Cdata
}

// 微信服务器推送过来的消息(事件)的通用消息头.
type MsgHeader struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	ToUserName   Cdata
	FromUserName Cdata
	CreateTime   int64
	MsgType      Cdata
}

// 微信服务器推送过来的消息(事件)的合集.
type WXMsg struct {
	MsgHeader
	Event    Cdata // 事件类型, CLICK
	EventKey Cdata // 事件KEY值, 与自定义菜单接口中KEY值对应

	// Message
	MessageMsg

	// Menu
	MenuMsg
}

type EncryptMsg struct {
	XMLName    struct{} `xml:"xml" json:"-"`
	ToUserName Cdata
	Encrypt    Cdata
}


// Context 是 Handler 处理消息(事件)的上下文环境. 非并发安全!
type Context struct {
	*Wechat
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	Values url.Values // 回调请求 URL 的查询参数集合
	// 回调请求 URL 的加密方式参数: encrypt_type
	// 回调请求 URL 的消息体签名参数: msg_signature
	// 回调请求 URL 的签名参数: signature
	// 回调请求 URL 的时间戳参数: timestamp
	// 回调请求 URL 的随机数参数: nonce

	MsgSigned      bool
	MsgCiphertext []byte // 消息的密文文本
	MsgPlaintext  []byte // 消息的明文文本, xml格式
	WXMsg         *WXMsg // 消息

	AESKey []byte // 当前消息加密所用的 aes-key, 接受返回数据使用同一个!!!

	kvs map[string]interface{}
}

func (wx *Wechat) CreateHandler(f func(*Context, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx *Context
			err error
		)
		ctx = wx.NewContext(w, r)
		ctx.Values = r.URL.Query()

		// 初始化验证
		echostr := ctx.Values.Get("echostr")
		if echostr != "" {
			err = ctx.checkValid(true)
			if err == nil {
				fmt.Fprint(ctx.ResponseWriter, echostr)
				return
			}

		} else if r.Method == "POST" {
			// 加密验证消息体，并解密
			err = ctx.checkValid(false)
		} else {
			err = fmt.Errorf("无效的请求")
		}

		// TODO 调试代码
		fmt.Println(ctx.Values.Encode())
		fmt.Printf("%#v \n", ctx.WXMsg)

		f(ctx, err)
	}
}

func (wx *Wechat) NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		Wechat:         wx,
		ResponseWriter: w,
		Request:        r,
	}
	ctx.WXMsg = &WXMsg{}
	return ctx
}

// 验证回调URL是否有效
func (ctx *Context) checkValid(init bool) (err error) {
	values := ctx.Values

	signature := values.Get("signature")
	if signature == "" {
		err = errors.New("not found signature query parameter")
		return
	}
	timestamp := values.Get("timestamp")
	if timestamp == "" {
		err = errors.New("not found timestamp query parameter")
		return
	}
	nonce := values.Get("nonce")
	if nonce == "" {
		err = errors.New("not found nonce query parameter")
		return
	}

	if ctx.token == "" {
		err = errors.New("token was not set, see NewWechat function or SetToken method")
		return
	}

	// 初始化验证
	if init {
		wantSignature := Sha1Sign(ctx.token, timestamp, nonce)
		if signature != wantSignature {
			err = fmt.Errorf("check signature failed, have: %s, want: %s", signature, wantSignature)
		}
		return
	}

	// 消息体验证 并解析数据
	encrypt_type := values.Get("encrypt_type")
	msg_signature := values.Get("msg_signature")
	if encrypt_type == "aes" && msg_signature == "" {
		err = errors.New("not found signature query parameter")
		return
	}
	if encrypt_type == "aes" && ctx.encodingAESKey == "" {
		err = errors.New("AESKey was not set, see NewWechat function or SetEncodingAESKey method")
		return
	}

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	defer ctx.Request.Body.Close()

	// aes 加密 只验证消息体
	if encrypt_type == "aes" {
		// 密文信息
		ctx.MsgCiphertext = body
		encryptMsg := &EncryptMsg{}
		err = xml.Unmarshal(body, encryptMsg)
		if err != nil {
			return
		}
		wantSignature := Sha1Sign(ctx.token, timestamp, nonce, encryptMsg.Encrypt.String())
		if msg_signature != wantSignature {
			err = fmt.Errorf("check msg_signature failed, have: %s, want: %s", msg_signature, wantSignature)
			return
		}

		// 解析加密数据，以后都走一样的逻辑
		body, err = ctx.Decrypt(encryptMsg.Encrypt.String())
		if err != nil {
			return
		}
		// 标示为加密信息
		ctx.MsgSigned = true

	} else { // 不加密 验证url参数
		wantSignature := Sha1Sign(ctx.token, timestamp, nonce)
		if signature != wantSignature {
			err = fmt.Errorf("check signature failed, have: %s, want: %s", signature, wantSignature)
			return
		}
	}
	//明文信息
	ctx.MsgPlaintext = body
	// 解析xml数据
	err = xml.Unmarshal(body, ctx.WXMsg)

	return
}

// 获取消息类型
func (ctx *Context) GetMsgType() string {
	return ctx.WXMsg.MsgType.Cdata
}

// 获取事件类型
func (ctx *Context) GetMsgEvent() string {
	return ctx.WXMsg.Event.Cdata
}
package wechat

import (
	"encoding/xml"
	"fmt"
)

type ResponseEncryptMsg struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	Encrypt      Cdata
	MsgSignature Cdata
	Nonce        Cdata
	TimeStamp    string
}

func (ctx *Context) responseWecaht(data interface{}) (err error) {
	var by []byte
	by, err = xml.MarshalIndent(data, "", "  ")
	//by, err = xml.Marshal(xmlData)
	if err != nil {
		return
	}
	// 加密返回
	if ctx.MsgSigned {
		var (
			msg_encrypt   string
			msg_signature string
		)
		msg_encrypt, err = ctx.Encrypt(by)
		if err != nil {
			return
		}
		msg_signature = Sha1Sign(ctx.token, ctx.Values.Get("timestamp"), ctx.Values.Get("nonce"), msg_encrypt)
		response := ResponseEncryptMsg{TimeStamp: ctx.Values.Get("timestamp")}
		response.Encrypt.Cdata = msg_encrypt
		response.MsgSignature.Cdata = msg_signature
		response.Nonce.Cdata = ctx.Values.Get("nonce")
		by, err = xml.Marshal(response)
		if err != nil {
			return
		}
	}
	// TODO debug
	fmt.Println(string(by))
	ctx.ResponseWriter.Header().Set("Content-Type", "text/xml")
	_, err = fmt.Fprint(ctx.ResponseWriter, string(by))
	return
}

func (ctx *Context) ResponseOK() (err error) {
	_, err = fmt.Fprint(ctx.ResponseWriter, "success")
	return
}

func (ctx *Context) getResponseMsgHeader(MsgType string) MsgHeader {
	msgHeader := MsgHeader{
		CreateTime:   ctx.WXMsg.CreateTime,
		FromUserName: ctx.WXMsg.ToUserName,
		ToUserName:   ctx.WXMsg.FromUserName,
	}
	msgHeader.MsgType.Cdata = MsgType
	return msgHeader
}

// 文本消息
type RText struct {
	MsgHeader
	Content Cdata // 回复的消息内容(换行: 在content中能够换行, 微信客户端支持换行显示)
}

// http://mp.weixin.qq.com/wiki/1/6239b44c206cab9145b1d52c67e6c551.html
func (ctx *Context) ResponseText(content string) error {
	data := RText{MsgHeader: ctx.getResponseMsgHeader(MSGTYPE_TEXT)}
	data.Content.Cdata = content
	return ctx.responseWecaht(data)
}

// 图片消息
type RImage struct {
	MsgHeader
	Image struct {
		MediaId Cdata // 通过素材管理接口上传多媒体文件得到 MediaId
	}
}

func (ctx *Context) ResponseImage(mediaId string) error {
	data := RImage{MsgHeader: ctx.getResponseMsgHeader(MSGTYPE_IMAGE)}
	data.Image.MediaId.Cdata = mediaId
	return ctx.responseWecaht(data)
}

// 语音消息
type RVoice struct {
	MsgHeader
	Voice struct {
		MediaId Cdata // 通过素材管理接口上传多媒体文件得到 MediaId
	}
}

func (ctx *Context) ResponseVoice(mediaId string) error {
	data := RVoice{MsgHeader: ctx.getResponseMsgHeader(MSGTYPE_VOICE)}
	data.Voice.MediaId.Cdata = mediaId
	return ctx.responseWecaht(data)
}

// 视频消息
type RVideo struct {
	MsgHeader
	Video struct {
		MediaId     Cdata // 通过素材管理接口上传多媒体文件得到 MediaId
		Title       Cdata // 视频消息的标题, 可以为空
		Description Cdata // 视频消息的描述, 可以为空
	}
}

func (ctx *Context) ResponseVideo(mediaId, title, description string) error {
	data := RVideo{MsgHeader: ctx.getResponseMsgHeader(MSGTYPE_VIDEO)}
	data.Video.MediaId.Cdata = mediaId
	data.Video.Title.Cdata = title
	data.Video.Description.Cdata = description
	return ctx.responseWecaht(data)
}

// 音乐消息
type RMusic struct {
	MsgHeader
	Music struct {
		Title        Cdata // 音乐标题
		Description  Cdata // 音乐描述
		MusicURL     Cdata // 音乐链接
		HQMusicURL   Cdata // 高质量音乐链接, WIFI环境优先使用该链接播放音乐
		ThumbMediaId Cdata // 缩略图的媒体id，通过素材管理接口上传多媒体文件，得到的id
	}
}

func (ctx *Context) ResponseMusic(title, description, musicURL, hqMusicURL, thumbMediaId string) error {
	data := RMusic{MsgHeader: ctx.getResponseMsgHeader(MSGTYPE_MUSIC)}
	data.Music.Title.Cdata = title
	data.Music.Description.Cdata = description
	data.Music.MusicURL.Cdata = musicURL
	data.Music.HQMusicURL.Cdata = hqMusicURL
	data.Music.ThumbMediaId.Cdata = thumbMediaId
	return ctx.responseWecaht(data)
}

// 图文消息里的 Article
type RArticle struct {
	Title       Cdata // 图文消息标题
	Description Cdata // 图文消息描述
	PicURL      Cdata // 图片链接, 支持JPG, PNG格式, 较好的效果为大图360*200, 小图200*200
	URL         Cdata // 点击图文消息跳转链接
}

func SetArticle(title, description, picurl, url string) RArticle {
	a := RArticle{}
	a.Title.Cdata = title
	a.Description.Cdata = description
	a.PicURL.Cdata = picurl
	a.URL.Cdata = url
	return a
}

// 图文消息
type News struct {
	MsgHeader
	ArticleCount int        // 图文消息个数, 限制为10条以内
	Articles     []RArticle `xml:"Articles>item,omitempty"` // 多条图文消息信息, 默认第一个item为大图, 注意, 如果图文数超过10, 则将会无响应
}

func (ctx *Context) ResponseNews(articles []RArticle) error {
	data := News{MsgHeader: ctx.getResponseMsgHeader(MSGTYPE_NEWS)}
	data.ArticleCount = len(articles)
	data.Articles = articles
	return ctx.responseWecaht(data)
}

// 将消息转发到多客服, 参见多客服模块
type RTransferCustomerService struct {
	MsgHeader
	TransInfo struct {
		KfAccount Cdata
	}
}

// http://mp.weixin.qq.com/wiki/11/f0e34a15cec66fefb28cf1c0388f68ab.html
// 如果不指定客服则 kfAccount 留空.
func (ctx *Context) ResponseTransferCustomerService(kfAccount string) error {
	if kfAccount == "" {
		data := ctx.getResponseMsgHeader(MSGTYPE_DUOKEFU)
		return ctx.responseWecaht(data)
	} else {
		data := RTransferCustomerService{MsgHeader: ctx.getResponseMsgHeader(MSGTYPE_DUOKEFU)}
		data.TransInfo.KfAccount.Cdata = kfAccount
		return ctx.responseWecaht(data)
	}
}

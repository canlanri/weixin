package wechat

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// 客服接口
// http://mp.weixin.qq.com/wiki/11/c88c270ae8935291626538f9c64bd123.html

// Add 添加客服账号.
//  account:         完整客服账号，格式为：账号前缀@公众号微信号，账号前缀最多10个字符，必须是英文或者数字字符。
//  nickname:        客服昵称，最长6个汉字或12个英文字符
//  password:        客服账号登录密码
func (wx *Wechat) AddKF(account, nickname, password string) (err error) {

	urlstr := fmt.Sprintf("%s/customservice/kfaccount/add?access_token=%s", wx.apiUrl, wx.accessToken)

	md5Sum := md5.Sum([]byte(password))
	password = hex.EncodeToString(md5Sum[:])

	request := struct {
		Account  string `json:"kf_account"`
		Nickname string `json:"nickname"`
		Password string `json:"password"`
	}{
		Account:  account,
		Nickname: nickname,
		Password: password,
	}

	var result WechatErr
	err = PostJSON(urlstr, &request, &result)

	return
}

// Update 设置客服信息(增量更新, 不更新的可以留空).
//  account:         完整客服账号，格式为：账号前缀@公众号微信号
//  nickname:        客服昵称，最长6个汉字或12个英文字符
//  password:        客服账号登录密码
func (wx *Wechat) UpdateKF(account, nickname, password string) (err error) {
	urlstr := fmt.Sprintf("%s/customservice/kfaccount/update?access_token=%s", wx.apiUrl, wx.accessToken)

	md5Sum := md5.Sum([]byte(password))
	password = hex.EncodeToString(md5Sum[:])

	request := struct {
		Account  string `json:"kf_account"`
		Nickname string `json:"nickname,omitempty"`
		Password string `json:"password,omitempty"`
	}{
		Account:  account,
		Nickname: nickname,
		Password: password,
	}

	var result WechatErr
	err = PostJSON(urlstr, request, &result)

	return
}

// Delete 删除客服账号
func (wx *Wechat) DeleteKF(kfAccount string) (err error) {
	// TODO 官方文档描述不明确
	urlstr := fmt.Sprintf("%s/customservice/kfaccount/del?kf_account=%s&access_token=%s", wx.apiUrl, kfAccount, wx.accessToken)

	var result WechatErr
	err = GetJSON(urlstr, &result)

	return
}

/*
TODO 微信部分接口，返回参数比较坑爹
errcode 一般为整型，这里却返回了字符串
在解析json时,就会报错
{
"errcode": invalid kf_account,
"errmsg": "invalid kf_account"
}
*/

// UploadHeadImage 上传客服头像
// 头像图片文件必须是jpg格式，推荐使用640*640大小的图片以达到最佳效果
func (wx *Wechat) UploadKFHeadImage(kfAccount, filePath string) (err error) {

	urlstr := fmt.Sprintf("%s/customservice/kfaccount/uploadheadimg?kf_account=%s&access_token=%s", wx.apiUrl, kfAccount, wx.accessToken)

	var result WechatErr

	multipartFormField := []MultipartFormField{
		MultipartFormField{IsFile: true, Fieldname: "media", Value: filePath},
	}

	err = PostMultipartForm(urlstr, multipartFormField, &result)

	return
}

// 客服基本信息
type KFInfo struct {
	Id           int    `json:"kf_id"`         // 客服工号
	Account      string `json:"kf_account"`    // 完整客服账号，格式为：账号前缀@公众号微信号
	Nickname     string `json:"kf_nick"`       // 客服昵称
	HeadImageURL string `json:"kf_headimgurl"` // 客服头像
}

// KFList 获取客服基本信息.
func (wx *Wechat) GetKFList() (list []KFInfo, err error) {

	urlstr := fmt.Sprintf("%s/customservice/getkflist?access_token=%s", wx.apiUrl, wx.accessToken)

	var result struct {
		KFList []KFInfo `json:"kf_list"`
	}
	err = GetJSON(urlstr, &result)
	if err != nil {
		return
	}
	list = result.KFList

	return
}

// 客服接口-发消息
// 如果需要以某个客服帐号来发消息（在微信6.0.2及以上版本中显示自定义头像）则需在JSON数据包的后半部分加入kf_account
type commonKFHeader struct {
	Touser        string `json:"touser"`
	Msgtype       string `json:"msgtype"`
	CustomService struct {
		KF_account string `json:"kf_account,omitempty"`
	} `json:"customservice,omitempty"`
}

// 发送文本消息
func (wx *Wechat) SendKFText(openid, text, kf_account string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		commonKFHeader
		Text struct {
			Content string `json:"content"`
		} `json:"text"`
	}{}

	request.Touser = openid
	request.Msgtype = MSGTYPE_TEXT
	request.Text.Content = text

	if kf_account != "" {
		request.CustomService.KF_account = kf_account
	}

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送图片消息
func (wx *Wechat) SendKFImage(openid, media_id, kf_account string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		commonKFHeader
		Image struct {
			MediaId string `json:"media_id"`
		} `json:"image"`
	}{}

	request.Touser = openid
	request.Msgtype = MSGTYPE_IMAGE
	request.Image.MediaId = media_id

	if kf_account != "" {
		request.CustomService.KF_account = kf_account
	}

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送语音消息
func (wx *Wechat) SendKFVoice(openid, media_id, kf_account string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		commonKFHeader
		Voice struct {
			MediaId string `json:"media_id"`
		} `json:"voice"`
	}{}

	request.Touser = openid
	request.Msgtype = MSGTYPE_VOICE
	request.Voice.MediaId = media_id

	if kf_account != "" {
		request.CustomService.KF_account = kf_account
	}

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送视频消息
func (wx *Wechat) SendKFVideo(openid, media_id, thumb_media_id, title, description, kf_account string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		commonKFHeader
		Video struct {
			MediaId        string `json:"media_id"`
			Thumb_media_id string `json:"thumb_media_id"`
			Title          string `json:"title"`
			Description    string `json:"description"`
		} `json:"video"`
	}{}

	request.Touser = openid
	request.Msgtype = MSGTYPE_VIDEO

	request.Video.MediaId = media_id
	request.Video.Thumb_media_id = thumb_media_id
	request.Video.Title = title
	request.Video.Description = description

	if kf_account != "" {
		request.CustomService.KF_account = kf_account
	}

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送音乐消息
func (wx *Wechat) SendKFMusic(openid, title, description, musicurl, hqmusicurl, thumb_media_id, kf_account string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		commonKFHeader
		Music struct {
			Title          string `json:"title"`
			Description    string `json:"description"`
			Musicurl       string `json:"musicurl"`
			Hqmusicurl     string `json:"hqmusicurl"`
			Thumb_media_id string `json:"thumb_media_id"`
		} `json:"music"`
	}{}

	request.Touser = openid
	request.Msgtype = MSGTYPE_MUSIC

	request.Music.Title = title
	request.Music.Description = description
	request.Music.Musicurl = musicurl
	request.Music.Hqmusicurl = hqmusicurl
	request.Music.Thumb_media_id = thumb_media_id

	if kf_account != "" {
		request.CustomService.KF_account = kf_account
	}

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 图文消息里的 Article
type KFArticle struct {
	Title       string `json:"title"`       // 图文消息标题
	Description string `json:"description"` // 图文消息描述
	PicURL      string `json:"picurl"`      // 图片链接, 支持JPG, PNG格式, 较好的效果为大图360*200, 小图200*200
	URL         string `json:"url"`         // 点击图文消息跳转链接
}

// 发送图文消息（点击跳转到外链）
// 图文消息条数限制在8条以内，注意，如果图文数超过8，则将会无响应。
func (wx *Wechat) SendKFNews(openid string, KFArticles []KFArticle, kf_account string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		commonKFHeader
		News struct {
			KFArticles []KFArticle `json:"articles"`
		} `json:"news"`
	}{}

	request.Touser = openid
	request.Msgtype = MSGTYPE_NEWS

	request.News.KFArticles = KFArticles

	if kf_account != "" {
		request.CustomService.KF_account = kf_account
	}

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送图文消息（点击跳转到图文消息页面）
//  图文消息条数限制在8条以内，注意，如果图文数超过8，则将会无响应。
func (wx *Wechat) SendKFMpnews(openid, media_id, kf_account string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		commonKFHeader
		Mpnews struct {
			MediaId string `json:"media_id"`
		} `json:"mpnews"`
	}{}

	request.Touser = openid
	request.Msgtype = MSGTYPE_MPNEWS

	request.Mpnews.MediaId = media_id

	if kf_account != "" {
		request.CustomService.KF_account = kf_account
	}

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送卡券
// 特别注意客服消息接口投放卡券仅支持非自定义Code码的卡券
func (wx *Wechat) SendKFWxcard(openid, card_id, card_ext, kf_account string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		commonKFHeader
		Wxcard struct {
			Card_id  string `json:"card_id"`
			Card_ext string `json:"card_ext"`
		} `json:"wxcard"`
	}{}

	request.Touser = openid
	request.Msgtype = MSGTYPE_WXCARD

	request.Wxcard.Card_id = card_id
	request.Wxcard.Card_ext = card_ext

	if kf_account != "" {
		request.CustomService.KF_account = kf_account
	}

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 高级群发接口
// http://mp.weixin.qq.com/wiki/15/40b6865b893947b764e2de8e4a1fb55f.html

//上传图文消息内的图片获取URL【订阅号与服务号认证后均可用】
//请注意，本接口所上传的图片不占用公众号的素材库中图片数量的5000个的限制。图片仅支持jpg/png格式，大小必须在1MB以下。
func (wx *Wechat) UploadMassImage(filePath string) (url string, err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/media/uploadimg?access_token=%s", wx.apiUrl, wx.accessToken)

	var result struct {
		Url string `json:"url"`
	}

	multipartFormField := []MultipartFormField{
		MultipartFormField{IsFile: true, Fieldname: "media", Value: filePath},
	}

	err = PostMultipartForm(urlstr, multipartFormField, &result)
	if err != nil {
		return
	}

	url = result.Url
	return
}

//上传图文消息素材【订阅号与服务号认证后均可用】
type MassNews struct { // 图文消息，一个图文消息支持1到8条图文
	Thumb_media_id     string `json:"thumb_media_id"`     //图文消息缩略图的media_id，可以在基础支持-上传多媒体文件接口中获得
	Author             string `json:"author"`             //图文消息的作者
	Title              string `json:"title"`              //图文消息的标题
	Content_source_url string `json:"content_source_url"` //在图文消息页面点击“阅读原文”后的页面
	Content            string `json:"content"`            //文消息页面的内容，支持HTML标签。具备微信支付权限的公众号，可以使用a标签，其他公众号不能使用
	Digest             string `json:"digest"`             //图文消息的描述
	Show_cover_pic     string `json:"show_cover_pic"`     //是否显示封面，1为显示，0为不显示
}

type MassResponse struct {
	Type       string `json:"type"`       //媒体文件类型，分别有图片（image）、语音（voice）、视频（video）和缩略图（thumb），次数为news，即图文消息
	Media_id   string `json:"media_id"`   //媒体文件/图文消息上传后获取的唯一标识
	Created_at string `json:"created_at"` //媒体文件上传时间
}

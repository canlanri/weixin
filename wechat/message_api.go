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
func (wx *Wechat) UploadHeadImage(kfAccount, filePath string) (err error) {

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

// 发送文本消息
func (wx *Wechat) SendKFText(openid, text ,kf_account string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}{
		Touser:  openid,
		Msgtype: MSGTYPE_TEXT,
	}

	request.Text.Content = openid

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送图片消息
func (wx *Wechat) SendKFImage(openid, media_id string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		Image   struct {
			MediaId string `json:"media_id"`
		} `json:"image"`
	}{
		Touser:  openid,
		Msgtype: MSGTYPE_TEXT,
	}

	request.Image.MediaId = media_id

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送语音消息
func (wx *Wechat) SendKFVoice(openid, media_id string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		Voice   struct {
			MediaId string `json:"media_id"`
		} `json:"voice"`
	}{
		Touser:  openid,
		Msgtype: MSGTYPE_VOICE,
	}

	request.Voice.MediaId = media_id

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送视频消息
func (wx *Wechat) SendKFVideo(openid, media_id, thumb_media_id, title, description string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		Video   struct {
			MediaId        string `json:"media_id"`
			Thumb_media_id string `json:"thumb_media_id"`
			Title          string `json:"title"`
			Description    string `json:"description"`
		} `json:"video"`
	}{
		Touser:  openid,
		Msgtype: MSGTYPE_VIDEO,
	}

	request.Video.MediaId = media_id
	request.Video.Thumb_media_id = thumb_media_id
	request.Video.Title = title
	request.Video.Description = description

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送音乐消息
func (wx *Wechat) SendKFMusic(openid, title, description, musicurl, hqmusicurl, thumb_media_id string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		Music   struct {
			Title          string `json:"title"`
			Description    string `json:"description"`
			Musicurl       string `json:"musicurl"`
			Hqmusicurl     string `json:"hqmusicurl"`
			Thumb_media_id string `json:"thumb_media_id"`
		} `json:"music"`
	}{
		Touser:  openid,
		Msgtype: MSGTYPE_MUSIC,
	}

	request.Music.Title = title
	request.Music.Description = description
	request.Music.Musicurl = musicurl
	request.Music.Hqmusicurl = hqmusicurl
	request.Music.Thumb_media_id = thumb_media_id

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
func (wx *Wechat) SendKFNews(openid string, KFArticles []KFArticle) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		News    struct {
			KFArticles []KFArticle `json:"articles"`
		} `json:"news"`
	}{
		Touser:  openid,
		Msgtype: MSGTYPE_NEWS,
	}

	request.News.KFArticles = KFArticles

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送图文消息（点击跳转到图文消息页面）
//  图文消息条数限制在8条以内，注意，如果图文数超过8，则将会无响应。
func (wx *Wechat) SendKFMpnews(openid, media_id string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		Mpnews  struct {
			MediaId string `json:"media_id"`
		} `json:"mpnews"`
	}{
		Touser:  openid,
		Msgtype: MSGTYPE_MPNEWS,
	}

	request.Mpnews.MediaId = media_id

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

// 发送卡券
// 特别注意客服消息接口投放卡券仅支持非自定义Code码的卡券
func (wx *Wechat) SendKFWxcard(openid, card_id, card_ext string) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/message/custom/send?access_token=%s", wx.apiUrl, wx.accessToken)

	request := struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		Wxcard  struct {
			Card_id  string `json:"card_id"`
			Card_ext string `json:"card_ext"`
		} `json:"wxcard"`
	}{
		Touser:  openid,
		Msgtype: MSGTYPE_WXCARD,
	}

	request.Wxcard.Card_id = card_id
	request.Wxcard.Card_ext = card_ext

	var result WechatErr
	err = PostJSON(urlstr, request, &result)
	return
}

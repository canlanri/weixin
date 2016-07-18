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
	err = PostJSON(urlstr, &request, &result)

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





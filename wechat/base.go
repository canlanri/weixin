package wechat

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"
)

const API_URL = "https://api.weixin.qq.com"
const API_URL_SH = "https://api.weixin.qq.com"
const API_URL_SZ = "https://api.weixin.qq.com"
const API_URL_HK = "https://api.weixin.qq.com"

type Wechat struct {
	lock           sync.RWMutex
	apiUrl         string
	token          string
	encodingAESKey string
	preAESKey      []byte
	nowAESKey      []byte
	appid          string
	appsecret      string
	accessToken    string
	expiresIn      time.Time
	refreshToken   int32
}

func NewWechat(token, appid, appsecret string) *Wechat {
	wx := &Wechat{
		apiUrl:    API_URL,
		token:     token,
		appid:     appid,
		appsecret: appsecret,
	}
	wx.CheckToken()
	return wx
}

// SetToken 设置签名token.
func (wx *Wechat) SetToken(token string) (err error) {
	if token == "" {
		return errors.New("empty token")
	}
	wx.lock.Lock()
	wx.token = token
	wx.lock.Unlock()
	return
}

func (wx *Wechat) SetEncodingAESKey(encodingAESKey string) (err error) {
	if encodingAESKey == "" {
		return errors.New("empty encodingAESKey")
	}
	aesKey, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil {
		return err
	}
	if len(aesKey) != 32 {
		return errors.New("encodingAESKey error")
	}
	wx.lock.Lock()
	wx.preAESKey = wx.nowAESKey
	wx.nowAESKey = aesKey
	wx.lock.Unlock()
	return
}

// 获取access_token
// http://mp.weixin.qq.com/wiki/14/9f9c82c1af308e3b14ba9b973f99a8ba.html
func (wx *Wechat) getAccessToken(code string) (err error) {
	if code == "" {
		code = "client_credential"
	}

	urlstr := fmt.Sprintf("%s/cgi-bin/token?grant_type=%s&appid=%s&secret=%s", wx.apiUrl, code, wx.appid, wx.appsecret)

	var result struct {
		AccessToken string        `json:"access_token"`
		ExpiresIn   time.Duration `json:"expires_in"`
	}

	err = GetJSON(urlstr, &result)
	if err != nil {
		return
	}

	wx.accessToken = result.AccessToken
	wx.expiresIn = time.Now().Add((result.ExpiresIn - 30) * time.Second)

	return
}

// 检测access_token是否过期
// 如果过期就自动更新
func (wx *Wechat) CheckToken() (err error) {
	if wx.expiresIn.After(time.Now()) {
		return nil
	}
	wx.lock.Lock()
	defer wx.lock.Unlock()
	// 同时更新时进行二次判断
	if wx.expiresIn.After(time.Now()) {
		return nil
	}
	err = wx.getAccessToken("")
	return
}

// 获取微信服务器IP地址
// http://mp.weixin.qq.com/wiki/4/41ef0843d6e108cf6b5649480207561c.html
func (wx *Wechat) GetCallbackIP() (ip []string, err error) {
	urlstr := fmt.Sprintf("%s/cgi-bin/getcallbackip?access_token=%s", wx.apiUrl, wx.accessToken)

	var result struct {
		List []string `json:"ip_list"`
	}

	err = GetJSON(urlstr, &result)
	if err != nil {
		return
	}

	ip = result.List
	return
}

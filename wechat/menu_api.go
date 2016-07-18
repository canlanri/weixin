package wechat

import (
	"fmt"
)

type Menu struct {
	Buttons   []Button   `json:"button,omitempty"`
	MatchRule *MatchRule `json:"matchrule,omitempty"`
	MenuId    int64      `json:"menuid,omitempty"` // 有个性化菜单时查询接口返回值包含这个字段
}

type MatchRule struct {
	GroupId            *int64 `json:"group_id,omitempty"`
	Sex                *int   `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType *int   `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
}

type Button struct {
	Type       string   `json:"type,omitempty"`       // 非必须; 菜单的响应动作类型
	Name       string   `json:"name,omitempty"`       // 必须;  菜单标题
	Key        string   `json:"key,omitempty"`        // 非必须; 菜单KEY值, 用于消息接口推送
	URL        string   `json:"url,omitempty"`        // 非必须; 网页链接, 用户点击菜单可打开链接
	MediaId    string   `json:"media_id,omitempty"`   // 非必须; 调用新增永久素材接口返回的合法media_id
	SubButtons []Button `json:"sub_button,omitempty"` // 非必须; 二级菜单数组
}

// 设置 btn 指向的 Button 为 子菜单 类型按钮.
func (btn *Button) SetAsSubMenuButton(name string, subButtons []Button) {
	btn.Name = name
	btn.SubButtons = subButtons

	btn.Type = ""
	btn.Key = ""
	btn.URL = ""
	btn.MediaId = ""
}

// 设置 btn 指向的 Button 为 click 类型按钮.
func (btn *Button) SetAsClickButton(name, key string) {
	btn.Type = "click" // 点击推事件
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 view 类型按钮.
func (btn *Button) SetAsViewButton(name, url string) {
	btn.Type = "view" // 跳转URL
	btn.Name = name
	btn.URL = url

	btn.Key = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 下面的按钮类型仅支持微信 iPhone5.4.1 以上版本, 和 Android5.4 以上版本的微信用户,
// 旧版本微信用户点击后将没有回应, 开发者也不能正常接收到事件推送.
// 设置 btn 指向的 Button 为 扫码推事件 类型按钮.
func (btn *Button) SetAsScanCodePushButton(name, key string) {
	btn.Type = "scancode_push" // 扫码推事件
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 扫码推事件且弹出"消息接收中"提示框 类型按钮.
func (btn *Button) SetAsScanCodeWaitMsgButton(name, key string) {
	btn.Type = "scancode_waitmsg" // 扫码带提示
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 弹出系统拍照发图 类型按钮.
func (btn *Button) SetAsPicSysPhotoButton(name, key string) {
	btn.Type = "pic_sysphoto" // 系统拍照发图
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 弹出拍照或者相册发图 类型按钮.
func (btn *Button) SetAsPicPhotoOrAlbumButton(name, key string) {
	btn.Type = "pic_photo_or_album" // 拍照或者相册发图
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 弹出微信相册发图器 类型按钮.
func (btn *Button) SetAsPicWeixinButton(name, key string) {
	btn.Type = "pic_weixin" // 微信相册发图
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 弹出地理位置选择器 类型按钮.
func (btn *Button) SetAsLocationSelectButton(name, key string) {
	btn.Type = "location_select" // 发送位置
	btn.Name = name
	btn.Key = key

	btn.URL = ""
	btn.MediaId = ""
	btn.SubButtons = nil
}

// 下面的按钮类型专门给第三方平台旗下未微信认证(具体而言, 是资质认证未通过)的订阅号准备的事件类型,
// 它们是没有事件推送的, 能力相对受限, 其他类型的公众号不必使用.
// 设置 btn 指向的 Button 为 下发消息(除文本消息) 类型按钮.
func (btn *Button) SetAsMediaIdButton(name, mediaId string) {
	btn.Type = "media_id" // 下发消息
	btn.Name = name
	btn.MediaId = mediaId

	btn.Key = ""
	btn.URL = ""
	btn.SubButtons = nil
}

// 设置 btn 指向的 Button 为 跳转图文消息URL 类型按钮.
func (btn *Button) SetAsViewLimitedButton(name, mediaId string) {
	btn.Type = "view_limited" // 跳转图文消息URL
	btn.Name = name
	btn.MediaId = mediaId

	btn.Key = ""
	btn.URL = ""
	btn.SubButtons = nil
}

// 自定义菜单创建接口
// http://mp.weixin.qq.com/wiki/10/0234e39a2025342c17a7d23595c6b40a.html
func (wx *Wechat) CreateMenu(menu Menu) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/menu/create?access_token=%s", wx.apiUrl, wx.accessToken)

	var result WechatErr
	err = PostJSON(urlstr, menu, &result)
	if err != nil {
		return
	}
	return
}

// 自定义菜单查询接口
// http://mp.weixin.qq.com/wiki/5/f287d1a5b78a35a8884326312ac3e4ed.html
func (wx *Wechat) GetMenu() (menu Menu, conditionalMenus []Menu, err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/menu/get?access_token=%s", wx.apiUrl, wx.accessToken)

	var result struct {
		Menu             Menu   `json:"menu"`
		ConditionalMenus []Menu `json:"conditionalmenu"`
	}
	if err = GetJSON(urlstr, &result); err != nil {
		return
	}

	menu = result.Menu
	conditionalMenus = result.ConditionalMenus
	return
}

// 自定义菜单删除接口
// http://mp.weixin.qq.com/wiki/3/de21624f2d0d3dafde085dafaa226743.html
func (wx *Wechat) DeleteMenu() (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/menu/delete?access_token=%s", wx.apiUrl, wx.accessToken)

	var result WechatErr
	if err = GetJSON(urlstr, &result); err != nil {
		return
	}
	return
}

// 创建个性化菜单
// http://mp.weixin.qq.com/wiki/0/c48ccd12b69ae023159b4bfaa7c39c20.html#.E5.88.9B.E5.BB.BA.E4.B8.AA.E6.80.A7.E5.8C.96.E8.8F.9C.E5.8D.95
func (wx *Wechat) AddConditionalMenu(menu Menu) (menuId int64, err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/menu/addconditional?access_token=%s", wx.apiUrl, wx.accessToken)

	var result struct {
		MenuId int64 `json:"menuId"`
	}
	if err = PostJSON(urlstr, menu, &result); err != nil {
		return
	}

	menuId = result.MenuId
	return
}

// 删除个性化菜单
// http://mp.weixin.qq.com/wiki/0/c48ccd12b69ae023159b4bfaa7c39c20.html#.E5.88.A0.E9.99.A4.E4.B8.AA.E6.80.A7.E5.8C.96.E8.8F.9C.E5.8D.95
func (wx *Wechat) DeleteConditionalMenu(menuId int64) (err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/menu/delconditional?access_token=%s", wx.apiUrl, wx.accessToken)

	var request = fmt.Sprintf(`{"menuid":"%d"}`, menuId)

	var result WechatErr
	if err = PostJSON(urlstr, request, &result); err != nil {
		return
	}
	return
}

// 测试个性化菜单匹配结果
// userId 可以是粉丝的 OpenID, 也可以是粉丝的微信号
// http://mp.weixin.qq.com/wiki/0/c48ccd12b69ae023159b4bfaa7c39c20.html#.E6.B5.8B.E8.AF.95.E4.B8.AA.E6.80.A7.E5.8C.96.E8.8F.9C.E5.8D.95.E5.8C.B9.E9.85.8D.E7.BB.93.E6.9E.9C
func (wx *Wechat) TryMatchConditionalMenu(userId string) (menu Menu, err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/menu/trymatch?access_token=%s", wx.apiUrl, wx.accessToken)

	var request = fmt.Sprintf(`{"user_id":"%d"}`, userId)

	var result struct {
		Menu `json:"menu"`
	}
	if err = PostJSON(urlstr, request, &result); err != nil {
		return
	}

	menu = result.Menu
	return
}

type MenuInfo struct {
	Buttons []ButtonEx `json:"button,omitempty"`
}

type ButtonEx struct {
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Key     string `json:"key,omitempty"`
	URL     string `json:"url,omitempty"`
	MediaId string `json:"media_id,omitempty"`

	Value    string `json:"value,omitempty"`
	NewsInfo struct {
		Articles []Article `json:"list,omitempty"`
	} `json:"news_info"`

	SubButton struct {
		Buttons []ButtonEx `json:"list,omitempty"`
	} `json:"sub_button"`
}

type Article struct {
	Title      string `json:"title,omitempty"`       // 图文消息的标题
	Author     string `json:"author,omitempty"`      // 作者
	Digest     string `json:"digest,omitempty"`      // 摘要
	ShowCover  int    `json:"show_cover"`            // 是否显示封面, 0为不显示, 1为显示
	CoverURL   string `json:"cover_url,omitempty"`   // 封面图片的URL
	ContentURL string `json:"content_url,omitempty"` // 正文的URL
	SourceURL  string `json:"source_url,omitempty"`  // 原文的URL, 若置空则无查看原文入口
}

// 获取自定义菜单配置接口
// http://mp.weixin.qq.com/wiki/14/293d0cb8de95e916d1216a33fcb81fd6.html
func (wx *Wechat) GetMenuInfo() (info MenuInfo, isMenuOpen bool, err error) {

	urlstr := fmt.Sprintf("%s/cgi-bin/get_current_selfmenu_info?access_token=%s", wx.apiUrl, wx.accessToken)

	var result struct {
		IsMenuOpen bool     `json:"is_menu_open"`
		MenuInfo   MenuInfo `json:"selfmenu_info"`
	}
	if err = GetJSON(urlstr, &result); err != nil {
		return
	}

	info = result.MenuInfo
	isMenuOpen = result.IsMenuOpen
	return
}




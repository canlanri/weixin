package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func GetJSON(urlstr string, v interface{}) error {
	return RequestJSON("GET", urlstr, nil, v)
}

func PostJSON(urlstr string, body, v interface{}) error {
	return RequestJSON("POST", urlstr, body, v)
}

func RequestJSON(method, urlstr string, bodyinterface, v interface{}) error {
	// TODO
	//fmt.Println(method, urlstr)

	var bodys io.Reader
	if bodyinterface != nil {
		switch bodyinterface.(type) {
		case string:
			bodys = bytes.NewBufferString(bodyinterface.(string))
		default:
			by, _ := json.Marshal(bodyinterface)
			bodys = bytes.NewBuffer(by)
		}
	}
	// TODO
	//fmt.Println(bodys)

	req, err := http.NewRequest(method, urlstr, bodys)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //一定要关闭resp.Body

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http.Status: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// 首先判断错误码
	r := WechatResponse{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	if r.ErrCode != ErrCodeOK {
		return r
	}
	if _, ok := v.(WechatResponse); ok {
		return nil
	}

	// 然后判断,其他数据格式
	err = json.Unmarshal(b, v)
	//fmt.Println(v)
	return err
}

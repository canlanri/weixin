package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func GetJSON(urlstr string, v interface{}) error {
	return RequestJSON("GET", urlstr, "application/json; charset=utf-8", nil, v)
}

func PostJSON(urlstr string, body, v interface{}) error {
	return RequestJSON("POST", urlstr, "application/json; charset=utf-8", body, v)
}

func RequestJSON(method, urlstr, bodyType string, bodyinterface, v interface{}) error {

	var bodys io.Reader
	if bodyinterface != nil {
		switch bodyinterface.(type) {
		case string:
			bodys = bytes.NewBufferString(bodyinterface.(string))
		case io.Reader:
			bodys = bodyinterface.(io.Reader)
		default:
			by, _ := json.Marshal(bodyinterface)
			bodys = bytes.NewBuffer(by)
		}
	}

	req, err := http.NewRequest(method, urlstr, bodys)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", bodyType)
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

	//fmt.Println(string(b))

	// 首先判断错误码
	r := &WechatErr{}
	err = json.Unmarshal(b, r)
	if err != nil {
		return fmt.Errorf("err:%s, json:%s",err,string(b))
	}

	if r.ErrCode != ErrCodeOK {
		return r
	}

	if _, ok := v.(WechatErr); ok {
		return nil
	}

	// 然后判断,其他数据格式
	err = json.Unmarshal(b, v)
	//fmt.Println(v)
	return err
}


type MultipartFormField struct {
	IsFile    bool
	Fieldname string
	Value     string
}

func PostMultipartForm(urlstr string, multipartFormField []MultipartFormField, v interface{}) error {

	body := bytes.NewBuffer(nil)
	mu := multipart.NewWriter(body)

	if len(multipartFormField) < 1 {
		return fmt.Errorf("multipartFormField is error")
	}

	for _, field := range multipartFormField {
		// 上传文件
		if field.IsFile {
			form, err := mu.CreateFormFile(field.Fieldname, field.Value)
			if err != nil {
				return err
			}

			file, err := os.Open(field.Value)
			if err != nil {
				return err
			}
			_, err = io.Copy(form, file)
			if err != nil {
				return err
			}
		} else { // 其他参数
			err := mu.WriteField(field.Fieldname, field.Value)
			if err != nil {
				return err
			}
		}
	}

	mu.Close()

	return RequestJSON("POST", urlstr, mu.FormDataContentType(), body, v)
}

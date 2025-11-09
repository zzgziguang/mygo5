package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	//md5
	MD5 := md5.New()
	_, _ = io.WriteString(MD5, "abc")
	md5Str := hex.EncodeToString(MD5.Sum(nil))
	fmt.Println(md5Str)

	//base64编码
	str := "abc1 2%=3你好"
	res := base64.StdEncoding.EncodeToString([]byte(str))
	fmt.Println(res)

	//base64解码
	s, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		return
	}
	fmt.Println(string(s))

	//base64编码
	res = base64.URLEncoding.EncodeToString([]byte(str))
	fmt.Println(res)

	//base64解码
	s, err = base64.URLEncoding.DecodeString(res)
	if err != nil {
		return
	}
	fmt.Println(string(s))

	//hexs := "3156EF"

	//解码
	res = hex.EncodeToString([]byte(str))
	fmt.Println(res)

	//编码
	s, err = hex.DecodeString(res)
	if err != nil {
		return
	}
	fmt.Println(string(s))

	return

	//go发送get请求
	resp, err := http.Get("https://www.baidu.com/")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应片段: %.300s...\n", string(body))

	//post urlencode请求
	data := url.Values{}
	data.Set("name", "张 三")
	data.Add("email", "zhang@123.com")
	strEncode := data.Encode()
	fmt.Printf("bianma%s", strEncode)
	resp, err = http.Post("http://httpbin.org/post", "application/x-www-form-urlencoded", strings.NewReader(strEncode))
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return
	}
	//req, err := http.NewRequest("POST", "http://httpbin.org/post", strings.NewReader(strEncode))
	// if err != nil {
	// 	panic(err)
	// }
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// client := &http.Client{}
	// resp, err = client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)

	//post json请求
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	resp, err = http.Post("http://httpbin.org/post", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
}

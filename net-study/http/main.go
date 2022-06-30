package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	//getExample()
	//postExample()
	clientExample()

}

func getExample() {
	urlString := "https://www.hengtiansoft.com/"
	data := url.Values{}
	data.Set("name", "Young")
	data.Set("age", "10")
	requestURI, err := url.ParseRequestURI(urlString)
	if err != nil {
		fmt.Println(err.Error())
	}
	requestURI.RawPath = data.Encode()
	fmt.Println(requestURI.String())
	resp, err := http.Get(requestURI.String())
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
}

func postExample() {
	urlString := "https://www.hengtiansoft.com/admin/doLogin"
	data := url.Values{}
	data.Set("user", "Young")
	data.Set("pwd", "10")
	data.Set("captcha", "TART")

	response, err := http.PostForm(urlString, data)
	printResInfo(response, err)
	defer response.Body.Close()
}

func clientExample() {
	//getUrlString := "https://www.hengtiansoft.com/"
	postUrlString := "https://account.adream.org/cas/loginApi"
	client := &http.Client{
		Transport:     http.DefaultTransport,
		Jar:           nil,
		CheckRedirect: nil,
		Timeout:       10 * time.Second,
	}
	//fmt.Println("client.Get")
	//getRes, err := client.Get(getUrlString)
	//printResInfo(getRes, err)
	//defer getRes.Body.Close()
	var data = url.Values{}
	data.Add("username", "18668042850")
	data.Add("password", "mayang1996")
	data.Add("checkcode", "cuuj")
	data.Add("service", "https://www.adream.org/wp-admin/")

	//fmt.Println("client.Post")
	//postRes, err := client.PostForm(postUrlString, data)
	//printResInfo(postRes, err)
	//defer postRes.Body.Close()
	//readAll, _ := ioutil.ReadAll(postRes.Body)
	//fmt.Println(string(readAll))

	fmt.Println("client.Do")
	req, err := http.NewRequest("POST", postUrlString, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	doRes, err := client.Do(req)
	printResInfo(doRes, err)
	defer doRes.Body.Close()
	readAll, _ := ioutil.ReadAll(doRes.Body)
	fmt.Println(string(readAll))

	//fmt.Println("client.Head")
	//headRes, err := client.Head(getUrlString)
	//printResInfo(headRes, err)

}

func printResInfo(response *http.Response, err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s\t%d\n", response.Status, response.StatusCode)
	fmt.Println("============Response Header============")
	for key, value := range response.Header {
		fmt.Printf("%s:%s\n", key, value)
	}
	fmt.Println("=======================================")
}

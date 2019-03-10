package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type BingImg struct {
	Images []struct {
		Startdate     string        `json:"startdate"`
		Fullstartdate string        `json:"fullstartdate"`
		Enddate       string        `json:"enddate"`
		URL           string        `json:"url"`
		Urlbase       string        `json:"urlbase"`
		Copyright     string        `json:"copyright"`
		Copyrightlink string        `json:"copyrightlink"`
		Title         string        `json:"title"`
		Quiz          string        `json:"quiz"`
		Wp            bool          `json:"wp"`
		Hsh           string        `json:"hsh"`
		Drk           int           `json:"drk"`
		Top           int           `json:"top"`
		Bot           int           `json:"bot"`
		Hs            []interface{} `json:"hs"`
	} `json:"images"`
	Tooltips struct {
		Loading  string `json:"loading"`
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Walle    string `json:"walle"`
		Walls    string `json:"walls"`
	} `json:"tooltips"`
}

func main() {
	var (
		bingJson BingImg
		imgUrl   string
	)

	bingUrl := "https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=zh-CN"
	resp, err := http.Get(bingUrl)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	//decode json
	err = json.NewDecoder(resp.Body).Decode(&bingJson)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("%+v\n", bingJson)
	// parse image url
	parse, err := url.Parse(bingUrl)
	if err != nil {
		fmt.Println(err)
	}

	imgUrl = parse.Scheme + "://" + parse.Host + bingJson.Images[0].URL

	// if contains https
	if strings.Contains(bingJson.Images[0].URL, `https://`) {
		imgUrl = bingJson.Images[0].URL
	}

	fmt.Println("图片地址:", imgUrl)

	// download image content
	// fix x509: certificate signed by unknown
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{
		Transport: tr,
	}
	req, err := http.NewRequest("GET", imgUrl, nil)

	imgResp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer imgResp.Body.Close()

	// save image into wallpaper.jpg
	imgFile, err := os.Create("./download/wallpaper.jpg")
	if err != nil {
		fmt.Println(err)
	}

	written, err := io.Copy(imgFile, imgResp.Body)
	if err != nil {
		fmt.Println(err)
	}
	if written > 10000 {
		fmt.Println("Copyright", bingJson.Images[0].Copyright)
		fmt.Println("当前时间", time.Now().Format("2006-01-02"))
		fmt.Println("图片下载完成，保存为 wallpaper.jpg ")
	}
}

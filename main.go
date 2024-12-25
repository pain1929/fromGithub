package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

import "C"

// Asset 定义 assets 的结构
type Asset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
}

// Release 定义整体结构
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

//export FromGithub
func FromGithub(uname string, rep string, currentVer string, saveTo string) bool {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", uname, rep)
	res, err := http.Get(url)
	if err != nil {
		return false
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return false
	}

	// 读取响应体内容
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return false
	}

	// 解析 JSON 数据
	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return false
	}

	// 如果版本一致则不更新
	if release.TagName == currentVer {
		return false
	}

	resp, err := http.Get(release.Assets[0].BrowserDownloadURL)
	if err != nil {
		return false
	}

	create, err := os.Create(saveTo)
	if err != nil {
		return false
	}
	defer create.Close()

	_, err = io.Copy(create, resp.Body)

	if err != nil {
		return false
	}

	return true

}

func main() {
	//FromGithub("pain1929", "deepRockHack1929", "1.0.1", "data.zip")
}

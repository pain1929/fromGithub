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
		println(err.Error())
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			println(err.Error())
		}
	}(res.Body)

	if res.StatusCode != 200 {
		println(res.Status)
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
		println(err.Error())
		return false
	}

	create, err := os.Create(saveTo)
	if err != nil {
		println(err.Error())
		return false
	}
	defer func(create *os.File) {
		err := create.Close()
		if err != nil {
			println(err.Error())
		}
	}(create)

	_, err = io.Copy(create, resp.Body)

	if err != nil {
		println(err.Error())
		return false
	}

	return true

}

// 定义结构体，用于解析 GitHub API 返回的 JSON 数据
type GitHubFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadURL string `json:"download_url"`
}

//export GetTextFromGithub
func GetTextFromGithub(uname string, rep string, path string) *C.char {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", uname, rep, path)

	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		println(err.Error())
		return nil
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return nil
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		println("statusCode error")
		return nil
	}

	// 解析 JSON 响应
	var file GitHubFile
	err = json.Unmarshal(body, &file)
	if err != nil {
		println(err.Error())
		return nil
	}

	// 获取 README.md 文件内容
	if file.DownloadURL != "" {
		// 发送请求以获取文件内容
		resp, err := http.Get(file.DownloadURL)
		if err != nil {
			println(err.Error())
			return nil
		}
		defer resp.Body.Close()

		// 读取文件内容
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			println(err.Error())
			return nil
		}
		// 将 Go 字符串转换为 C 字符串
		return C.CString(string(content))
	}

	return nil
}

func main() {
	//FromGithub("pain1929", "deepRockHack1929", "1.0.1", "data.zip")
	github := GetTextFromGithub("pain1929", "deepRockHack1929", "README.md")

	println(github)
}

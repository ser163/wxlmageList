package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"time"
)

// Config 保存从 config.yaml 读取的配置
type Config struct {
	AppID  string `yaml:"appid"`
	Secret string `yaml:"secret"`
}

// AccessToken 表示本地保存的访问令牌
type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresAt int64  `json:"expires_at"`
}

// TokenResponse 表示获取新令牌时的 API 返回结果
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

// MaterialRequest 是获取素材的请求体
type MaterialRequest struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Count  int    `json:"count"`
}

// MaterialItem 表示素材列表中的单个素材
type MaterialItem struct {
	MediaID string `json:"media_id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
}

// MaterialResponse 表示素材 API 的返回结果
type MaterialResponse struct {
	Item    []MaterialItem `json:"item"`
	ErrCode int            `json:"errcode"`
	ErrMsg  string         `json:"errmsg"`
}

const tokenFile = "access_token.json"

func main() {
	// 加载配置文件
	log.Println("加载配置文件...")
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 获取访问令牌
	log.Println("获取访问令牌...")
	accessToken, err := getAccessToken(config)
	if err != nil {
		log.Fatalf("获取访问令牌失败: %v", err)
	}

	// 获取素材
	log.Println("获取媒体素材...")
	items, err := fetchMediaItems(accessToken)
	if err != nil {
		log.Fatalf("获取媒体素材失败: %v", err)
	}

	// 打印素材列表
	log.Println("打印素材列表:")
	for _, item := range items {
		fmt.Printf("媒体ID: %s, 名称: %s, 链接: %s\n", item.MediaID, item.Name, item.URL)
	}
}

// loadConfig 从 config.yaml 读取配置
func loadConfig(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("读取配置文件失败: %v", err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("解析配置文件失败: %v", err)
	}
	if config.AppID == "" || config.Secret == "" {
		return Config{}, fmt.Errorf("配置文件中缺少 appid 或 secret")
	}
	return config, nil
}

// getAccessToken 获取有效的访问令牌
func getAccessToken(config Config) (string, error) {
	token, err := readAccessToken()
	if err == nil && token.ExpiresAt > time.Now().Unix() {
		log.Println("使用本地缓存的访问令牌")
		return token.Token, nil
	}
	log.Println("本地令牌过期或不存在，重新获取...")
	newToken, err := fetchNewAccessToken(config)
	if err != nil {
		return "", err
	}
	err = saveAccessToken(newToken)
	if err != nil {
		return "", fmt.Errorf("保存访问令牌失败: %v", err)
	}
	return newToken.Token, nil
}

// readAccessToken 从本地文件读取访问令牌
func readAccessToken() (AccessToken, error) {
	data, err := os.ReadFile(tokenFile)
	if err != nil {
		return AccessToken{}, err // 文件不存在或无法读取
	}
	var token AccessToken
	err = json.Unmarshal(data, &token)
	if err != nil {
		return AccessToken{}, fmt.Errorf("解析访问令牌文件失败: %v", err)
	}
	return token, nil
}

// saveAccessToken 将访问令牌保存到本地文件
func saveAccessToken(token AccessToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("序列化访问令牌失败: %v", err)
	}
	err = os.WriteFile(tokenFile, data, 0644)
	if err != nil {
		return fmt.Errorf("写入访问令牌文件失败: %v", err)
	}
	return nil
}

// fetchNewAccessToken 从微信 API 获取新的访问令牌
func fetchNewAccessToken(config Config) (AccessToken, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", config.AppID, config.Secret)
	resp, err := http.Get(url)
	if err != nil {
		return AccessToken{}, fmt.Errorf("网络请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return AccessToken{}, fmt.Errorf("获取令牌失败，状态码: %s", resp.Status)
	}
	var tokenResp TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return AccessToken{}, fmt.Errorf("解析令牌响应失败: %v", err)
	}
	if tokenResp.AccessToken != "" && tokenResp.ExpiresIn != 0 {
		expiresAt := time.Now().Unix() + int64(tokenResp.ExpiresIn-5)
		return AccessToken{Token: tokenResp.AccessToken, ExpiresAt: expiresAt}, nil
	} else if tokenResp.ErrCode != 0 {
		return AccessToken{}, fmt.Errorf("API错误: %d - %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}
	return AccessToken{}, fmt.Errorf("令牌响应无效")
}

// fetchMediaItems 从微信 API 获取图片素材
func fetchMediaItems(accessToken string) ([]MaterialItem, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=%s", accessToken)
	reqBody, err := json.Marshal(MaterialRequest{Type: "image", Offset: 0, Count: 20})
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("网络请求失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取素材失败，状态码: %s", resp.Status)
	}
	var materialResp MaterialResponse
	err = json.NewDecoder(resp.Body).Decode(&materialResp)
	if err != nil {
		return nil, fmt.Errorf("解析素材响应失败: %v", err)
	}
	if materialResp.Item != nil {
		return materialResp.Item, nil
	} else if materialResp.ErrCode != 0 {
		return nil, fmt.Errorf("API错误: %d - %s", materialResp.ErrCode, materialResp.ErrMsg)
	}
	return nil, fmt.Errorf("素材响应无效")
}

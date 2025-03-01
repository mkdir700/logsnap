package uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"logsnap/remote"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// CloudreveUploader 实现Cloudreve存储上传
type CloudreveUploader struct {
	config  remote.UploadConfigProvider
	session *http.Client
}

func NewCloudreveUploader(config remote.UploadConfigProvider) *CloudreveUploader {
	return &CloudreveUploader{config: config}
}

func (c *CloudreveUploader) login() error {
	// 创建cookie jar用于存储会话cookie
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("创建cookie jar失败: %w", err)
	}

	// 准备登录数据
	loginData := map[string]string{
		"userName": c.config.Username,
		"Password": c.config.Password,
	}
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("序列化登录数据失败: %w", err)
	}

	// 创建HTTP请求
	loginURL := fmt.Sprintf("%s/api/v3/user/session", c.config.Endpoint)
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建登录请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 创建临时客户端发送登录请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("登录失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 获取并保存cookies
	cookies := resp.Cookies()
	hasSessionCookie := false
	for _, cookie := range cookies {
		if strings.Contains(cookie.Name, "cloudreve-session") {
			hasSessionCookie = true
			break
		}
	}

	if !hasSessionCookie {
		return fmt.Errorf("登录成功但未获取到cloudreve-session cookie")
	}

	// 设置cookies到jar
	cookieJar.SetCookies(req.URL, cookies)

	// 创建带有cookie jar的会话客户端
	c.session = &http.Client{
		Jar:     cookieJar,
		Timeout: 30 * time.Second,
	}

	logrus.Infof("Cloudreve登录成功，已保存会话cookie")
	return nil
}

// 获取webdav账户信息
// 返回账户ID和账户密码
//
//	{
//	    "code": 0,
//	    "data": {
//	        "accounts": [
//	            {
//	                "ID": 3,
//	                "CreatedAt": "2022-07-13T13:50:16.733477315+08:00",
//	                "UpdatedAt": "2022-07-13T13:50:16.733477315+08:00",
//	                "DeletedAt": null,
//	                "Name": "HFR-Cloud挂载",
//	                "Password": "xxxxx",
//	                "UserID": 1,
//	                "Root": "/HFR-Cloud",
//	                "Readonly": false,
//	                "UseProxy": false
//	            }
//	        ],
//	        "folders": [
//	            {
//	                "id": "abcd",
//	                "name": "/",
//	                "policy_name": "HFR-Cloud"
//	            }
//	        ]
//	    },
//	    "msg": ""
//	}
func (c *CloudreveUploader) getWebdavCredentials() (remote.UploadConfigProvider, error) {
	// 创建HTTP请求
	loginURL := fmt.Sprintf("%s/api/v3/webdav/accounts", c.config.Endpoint)
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		return remote.UploadConfigProvider{}, fmt.Errorf("创建登录请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 发送请求，使用带有cookie的会话
	resp, err := c.session.Do(req)
	if err != nil {
		return remote.UploadConfigProvider{}, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return remote.UploadConfigProvider{}, fmt.Errorf("获取webdav账户信息失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response struct {
		Code int `json:"code"`
		Data struct {
			Accounts []struct {
				Password string `json:"Password"`
				Root     string `json:"Root"`
			} `json:"accounts"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return remote.UploadConfigProvider{}, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(response.Data.Accounts) == 0 {
		return remote.UploadConfigProvider{}, fmt.Errorf("未找到webdav账户")
	}

	return remote.UploadConfigProvider{
		Endpoint:   c.config.Endpoint + "/dav",
		Username:   c.config.Username,
		Password:   response.Data.Accounts[0].Password,
		FolderPath: response.Data.Accounts[0].Root,
	}, nil
}

// /api/v3/file/search/keywords%2F11baf85e5b680d85ea338599e309ab64_logs.zip?path=%2Fsnapshots%2F2025%2F03%2F02
func (c *CloudreveUploader) getFileID(fileURL string) (string, error) {
	// /dav/snapshots/2025/03/02/de322746df94435e5b403a6e629dfc34_logs.zip
	// 需要获取de322746df94435e5b403a6e629dfc34_logs.zip的fileID
	logrus.Infof("开始获取文件ID，文件URL: %s", fileURL)

	// 从 WebDAV URL 中提取路径
	objectKey := strings.TrimPrefix(fileURL, c.config.Endpoint+"/dav/")
	logrus.Infof("提取的对象键: %s", objectKey)

	splitItems := strings.Split(objectKey, "/")
	if len(splitItems) == 0 {
		return "", fmt.Errorf("无效的文件URL: %s", fileURL)
	}

	// de322746df94435e5b403a6e629dfc34_logs.zip
	fileName := splitItems[len(splitItems)-1]

	logrus.Infof("解析结果 - 文件名: %s", fileName)

	// URL 编码文件名和路径
	encodedFileName := url.QueryEscape(fileName)

	// 创建HTTP请求
	searchURL := fmt.Sprintf("%s/api/v3/file/search/keywords%%2F%s",
		c.config.Endpoint, encodedFileName)

	logrus.Infof("搜索URL: %s", searchURL)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建搜索请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 发送请求，使用带有cookie的会话
	resp, err := c.session.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %w", err)
	}

	logrus.Infof("搜索响应: %s", string(respBody))

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("搜索文件失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var response struct {
		Code int `json:"code"`
		Data struct {
			Objects []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Path string `json:"path"`
			} `json:"objects"`
		} `json:"data"`
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(response.Data.Objects) == 0 {
		return "", fmt.Errorf("未找到文件")
	}

	fileID := response.Data.Objects[0].ID
	logrus.Infof("找到文件 ID: %s, 名称: %s, 路径: %s",
		fileID, response.Data.Objects[0].Name, response.Data.Objects[0].Path)

	return fileID, nil
}

func (c *CloudreveUploader) createShareURL(fileID string) (string, error) {
	// 创建HTTP请求
	loginURL := fmt.Sprintf("%s/api/v3/share", c.config.Endpoint)
	shareData := map[string]interface{}{
		"id":        fileID,
		"is_dir":    false,
		"password":  "",
		"downloads": 10,
		"expire":    3600, // 1小时
		"preview":   false,
	}
	jsonData, err := json.Marshal(shareData)
	if err != nil {
		return "", fmt.Errorf("序列化分享数据失败: %w", err)
	}
	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建登录请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 发送请求，使用带有cookie的会话
	resp, err := c.session.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var response struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	return response.Data, nil
}

func (c *CloudreveUploader) Upload(localPath, objectKey string) (string, error) {
	// 确保已登录
	if c.session == nil {
		if err := c.login(); err != nil {
			return "", fmt.Errorf("登录失败: %w", err)
		}
	}

	webdavConfig, err := c.getWebdavCredentials()
	if err != nil {
		return "", fmt.Errorf("获取webdav账户信息失败: %w", err)
	}

	logrus.Infof("开始通过 WebDAV 上传文件: %s", localPath)
	webdavUploader := NewWebdavUploader(webdavConfig)
	webdavURL, err := webdavUploader.Upload(localPath, objectKey)
	if err != nil {
		return "", fmt.Errorf("上传失败: %w", err)
	}
	logrus.Infof("WebDAV 上传成功，URL: %s", webdavURL)

	// 等待文件索引完成
	time.Sleep(1 * time.Second)

	logrus.Infof("开始获取文件 ID...")
	fileID, err := c.getFileID(webdavURL)
	if err != nil {
		logrus.Warnf("获取文件ID失败: %v，将使用 WebDAV URL", err)
		return webdavURL, nil
	}

	logrus.Infof("获取到文件ID: %s", fileID)

	// 构建文件分享链接
	shareURL, err := c.createShareURL(fileID)
	if err != nil {
		return "", fmt.Errorf("创建分享链接失败: %w", err)
	}

	logrus.Infof("创建分享链接成功: %s", shareURL)

	return shareURL, nil
}

package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"pkg/config"
	"pkg/utils"
	"pkg/utils/logger"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileUploadController struct{}
type FilePaths struct {
	FilePathMap map[string]string `json:"filePaths"`
}

func (FileUploadController) UploadFile(c *gin.Context) {
	_, _ = strconv.Atoi(c.PostForm("index"))
	chunk, _ := c.FormFile("chunk")
	hash := c.PostForm("hash")
	dst := filepath.Join(config.ChunkDir, hash)
	if err := c.SaveUploadedFile(chunk, dst); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "文件保存失败")
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", "分块上传成功!")
}
func (FileUploadController) DeleteFile(c *gin.Context) {
	userID, _ := ParserToken(c)
	filename := c.Query("filename")
	fileDir := filepath.Join(config.UploadDir, userID)
	if filename == "" {
		utils.ReturnError(c, http.StatusBadRequest, "文件名不能为空")
		return
	}
	targetPath := filepath.Join(fileDir, filename)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		utils.ReturnError(c, http.StatusNotFound, "文件不存在")
		return
	}
	if err := os.RemoveAll(targetPath); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "文件删除失败")
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", "删除成功")
}
func (FileUploadController) MergeFileChunk(c *gin.Context) {
	filename := c.PostForm("filename")
	userID, _ := ParserToken(c)
	hashList := strings.Split(c.PostForm("hashList"), ",")
	fileDir := filepath.Join(config.UploadDir, userID)
	finalFilePath := filepath.Join(config.UploadDir, userID, filename)
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			utils.ReturnError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	os.Remove(finalFilePath)
	file, err := os.Create(finalFilePath)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "文件上传成功", nil)
	defer file.Close()
	for _, hash := range hashList {
		chunkPath := filepath.Join(config.ChunkDir, hash)
		matches, _ := filepath.Glob(chunkPath)
		chunkFile, err := os.Open(matches[0])
		if err != nil {
			utils.ReturnError(c, http.StatusInternalServerError, "打开文件块失败"+err.Error())
			return
		}
		defer chunkFile.Close()
		_, err = io.Copy(file, chunkFile)
		if err != nil {
			utils.ReturnError(c, http.StatusInternalServerError, "合并文件块失败"+err.Error())
			return
		}
		os.Remove(matches[0])
	}

}

func (FileUploadController) UploadFolder(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		utils.ReturnError(c, 1, "获取表单数据失败")
		return
	}
	files := form.File["files"]
	filePathsJSON := c.PostForm("filePaths")
	var filePaths FilePaths
	if filePathsJSON != "" {
		if err := json.Unmarshal([]byte(filePathsJSON), &filePaths.FilePathMap); err != nil {
			utils.ReturnError(c, 1, "解析文件路径映射失败")
			return
		}
	} else {
		utils.ReturnError(c, 1, "缺少文件路径映射")
		return
	}
	os.MkdirAll(config.UploadDir, os.ModePerm)
	logger.Debug("开始上传文件夹，文件路径映射：%v", filePaths.FilePathMap)
	for _, file := range files {
		relativePath, exists := filePaths.FilePathMap[file.Filename]
		if !exists {
			utils.ReturnError(c, 1, fmt.Sprintf("找不到文件 '%s' 的相对路径", file.Filename))
			return
		}
		saveDir := filepath.Join(config.UploadDir, filepath.Dir(relativePath))
		os.MkdirAll(saveDir, os.ModePerm)
		dst := filepath.Join(config.UploadDir, relativePath)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			utils.ReturnError(c, 1, fmt.Sprintf("保存文件 '%s' 失败: %v", file.Filename, err))
			return
		}
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", "文件夹上传成功")
}

type FileInfo struct {
	Path    string
	Name    string
	IsDir   bool
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time
}

func (FileUploadController) GetDirStruct(c *gin.Context) {
	var DirInfo []FileInfo
	userID, _ := ParserToken(c)
	if userID == "" {
		utils.ReturnError(c, http.StatusUnauthorized, "状态出错")
		return
	}
	personDir := filepath.Join(config.UploadDir, userID)
	err := filepath.Walk(personDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		DirInfo = append(DirInfo, FileInfo{
			Path:    path,
			Name:    info.Name(),
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			Mode:    info.Mode(),
			ModTime: info.ModTime(),
		})
		return nil
	})
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, nil)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", DirInfo)
}

func (FileUploadController) MarkdownImage(c *gin.Context) {
	DirPath := config.MarkdownImagesDir
	file, err := c.FormFile("file")
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "无法解析数据")
	}
	u := strings.ReplaceAll(uuid.New().String(), "-", "")
	extension := strings.Split(file.Filename, ".")
	var dst string
	if len(extension) > 1 {
		dst = filepath.Join(DirPath, fmt.Sprintf("/%s.%s", u, extension[len(extension)-1]))
	} else {
		dst = filepath.Join(DirPath, fmt.Sprintf("/%s.png", u))
	}
	if err := c.SaveUploadedFile(file, dst); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "文件保存失败")
	}
	updatedURL := strings.Replace(dst, "/var/data/oj/media", fmt.Sprintf("%s:%s", config.Address, config.Port), 1)
	utils.ReturnSuccess(c, http.StatusOK, "文件上传成功", updatedURL)
}

type FileShare struct {
	Pwd       string `json:"pwd"`
	Filename  string `json:"filename"`
	OwnerID   string `json:"owner_id"`
	OwnerName string `json:"owner_name"`
}

func (FileUploadController) CreateShareFile(c *gin.Context) { // 创建文件分享
	userID, _ := ParserToken(c)
	pwd := c.PostForm("pwd")
	filename := c.PostForm("filename") // 	拼出这个文件的地址
	if pwd == "" {
		utils.ReturnError(c, http.StatusBadRequest, "密钥不能为空")
		return
	}
	fileDir := filepath.Join(config.UploadDir, userID)
	targetPath := filepath.Join(fileDir, filename)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		utils.ReturnError(c, http.StatusNotFound, "文件不存在")
		return
	}
	// fileShare := FileShare{
	// 	Pwd:       pwd,
	// 	Filename:  filename,
	// 	OwnerID:   x.UserID,
	// 	OwnerName: "匿名用户",
	// }
}

func (FileUploadController) GetShareFile(c *gin.Context) {
	userID, _ := ParserToken(c)
	// key := c.Query("key")
	pwd := c.Query("pwd")
	// cache.Set(key, 24*time.Hour)
	filename := c.Query("filename")
	if pwd == "" {
		utils.ReturnError(c, http.StatusBadRequest, "密钥不能为空")
		return
	}
	fileDir := filepath.Join(config.UploadDir, userID)
	targetPath := filepath.Join(fileDir, filename)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		utils.ReturnError(c, http.StatusNotFound, "文件不存在")
		return
	}
	fileURL := fmt.Sprintf("%s:%s/%s/%s", config.Address, config.Port, userID, filename)
	utils.ReturnSuccess(c, http.StatusOK, "success", fileURL)
}

package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pkg/config"
	"pkg/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileDownloadController struct{}

func (FileDownloadController) GetTest(c *gin.Context) {
	filePath := "./Makefile"
	file, err := os.Open(filePath)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "无法打开文件")
		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileInfo.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	c.File(filePath)
}

func (FileDownloadController) GetChunk(c *gin.Context) {
	x, _ := ParserToken(c.Request.Header.Get("Authorization"))
	filepath := filepath.Join(config.UploadDir, x.UserID, c.Query("filename"))
	Range := strings.Split(c.Request.Header.Get("Range"), "-")
	start, err := strconv.ParseInt(Range[0], 10, 64)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "范围出错")
		return
	}
	end, err := strconv.ParseInt(Range[1], 10, 64)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "范围出错")
		return
	}
	if start > end {
		utils.ReturnError(c, http.StatusInternalServerError, "范围出错")
		return
	}
	fileBytes, err := readFileByRange(filepath, int64(start), int64(end))
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "找不到文件,请确认文件名是否正确")
		return
	}
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, int64(start)+int64(len(fileBytes))-1))
	c.Header("Content-Length", strconv.FormatInt(int64(len(fileBytes)), 10))
	c.Header("Content-Type", "application/octet-stream")
	c.Status(http.StatusPartialContent)
	c.Writer.Write(fileBytes)
	// utils.ReturnSuccess(c, http.StatusOK, "success", fileBytes, 1)
}

func readFileByRange(path string, start, end int64) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fs, _ := file.Stat()
	defer file.Close()
	_, Err := file.Seek(start, io.SeekStart)
	if Err != nil {
		return nil, Err
	}
	if end == 0 || end > fs.Size() {
		end = fs.Size() - 1
	}
	Length := end - start + 1
	buf := make([]byte, Length)
	bytesRead, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf[:bytesRead], nil
}

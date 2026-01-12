/*
 * @Author: x_wq3337 854541540@qq.com
 * @Date: 2025-07-10 15:34:29
 * @LastEditors: x_wq3337 854541540@qq.com
 * @LastEditTime: 2026-01-12 17:48:55
 * @FilePath: /backend/controllers/record.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package controllers

import (
	"net/http"
	"pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecodeController struct{}

func (RecodeController) GetRecodeInfo(c *gin.Context) {
	id := c.Param("id")

	record, err := models.QueryRecordById(id)
	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", record)
		return
	}
	ReturnError(c, http.StatusNotFound, err)
}

func (RecodeController) GetList(c *gin.Context) {
	// 获取分页参数，默认第1页，每页10条
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	// 转换分页参数
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 || pageSizeInt > 100 {
		pageSizeInt = 10
	}

	// 获取查询参数
	search := c.Query("search")
	verdict := c.Query("verdict")
	language := c.Query("language")

	// 调用模型层获取数据
	results, total, err := models.GetAllRecord(pageInt, pageSizeInt, search, verdict, language)

	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", gin.H{
			"data":      results,
			"total":     total,
			"page":      pageInt,
			"page_size": pageSizeInt,
		})
		return
	}
	ReturnError(c, http.StatusInternalServerError, err)
}
func (RecodeController) CreateRecord(c *gin.Context) {
	data := &models.Record{}
	_ = c.BindJSON(&data)
}

func (RecodeController) GetRecodeListByUser(c *gin.Context) {
	results, err := models.QueryRecordByUserId(c.Param("id"))
	if err != nil {
		ReturnError(c, http.StatusNotFound, err)
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", results)

}

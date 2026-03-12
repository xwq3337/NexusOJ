package microapps

import "github.com/gin-gonic/gin"

func MicroRouter(r *gin.Engine) {
	r.GET("/formatter-languages", FormatterListHandler)
	r.POST("/formatter", FormatterHandler)

}

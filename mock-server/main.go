package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/api/admin/login", func(ctx *gin.Context) {
		ctx.SetCookie("dumb-dumb", "you", 3600, "", "localhost:8080", false, true)
		ctx.JSON(http.StatusOK, gin.H{
			"dumb": "dumb",
		})
	})

	r.POST("/api/Home/GetAction", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "<p>Missing Attendance record 欠缺出入勤紀錄:<br /> 2023-12-21<br />2023-12-28<br />2024-01-08<br />2024-01-11<br />2024-01-12<br />2024-01-15</p><p>Early leave:<br /> 2023-12-18<br />2024-01-17</p><p>Lateness 遲到:<br /> 2023-12-18<br />2024-01-04</p>")
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

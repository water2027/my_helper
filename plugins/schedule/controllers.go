package schedule

import (
	"github.com/gin-gonic/gin"
)

type Controller struct{}

type ScheduleController interface {
	AddOnce(c *gin.Context)
	AddLong(c *gin.Context)

	AddPage(c *gin.Context)
	BrowsePage(c *gin.Context)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authCode := c.GetHeader("Authorization")
		if authCode != "1234" {
			c.JSON(400, gin.H{})
			return
		}
		c.Next()
	}
}

func (sp *SchedulePlugin) RegisterRoutes(r *gin.Engine) {
	c := &Controller{}
	r.Use(AuthMiddleware())
	r.POST("/schedule/add_once", c.AddOnce)
	r.POST("/schedule/add_long", c.AddLong)
	r.GET("/schedule/add_page", c.AddPage)
	r.GET("/schedule/browse_page", c.BrowsePage)
}

func (controller *Controller) AddOnce(c *gin.Context) {
	// AddOnce
}

func (controller *Controller) AddLong(c *gin.Context) {

}

func (controller *Controller) AddPage(c *gin.Context) {

}

func (controller *Controller) BrowsePage(c *gin.Context) {

}


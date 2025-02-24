package schedule

import (
	"github.com/gin-gonic/gin"
)

type Controller struct{
	Service ScheduleService
}

type ScheduleController interface {
	AddOnce(c *gin.Context)
	AddLong(c *gin.Context)
	DeleteTask(c *gin.Context)

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
	c := &Controller{
		Service: *NewScheduleService(),
	}
	group := r.Group("/schedule")
	group.Use(AuthMiddleware())
	group.GET("/", c.AddPage)
	group.GET("/browse_page", c.BrowsePage)

	group.POST("/add_once", c.AddOnce)
	group.POST("/add_long", c.AddLong)
	
	group.DELETE("/", c.DeleteTask)
}

func (controller *Controller) AddOnce(c *gin.Context) {
	// AddOnce
}

func (controller *Controller) AddLong(c *gin.Context) {

}

func (controller *Controller) DeleteTask(c *gin.Context) {

}

func (controller *Controller) AddPage(c *gin.Context) {

}

func (controller *Controller) BrowsePage(c *gin.Context) {

}

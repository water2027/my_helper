package schedule

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Controller struct{
	Service ScheduleService
}

type ScheduleController interface {
	AddOnce(c *gin.Context)
	AddLong(c *gin.Context)
	DeleteTask(c *gin.Context)

	Browse(c *gin.Context)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authCode := c.GetHeader("Authorization")
		if authCode == "1234" {
			c.Next()
			return
		}
		c.HTML(400, "schedule.html", nil)
		c.Abort()
	}
}

func (sp *SchedulePlugin) RegisterRoutes(r *gin.Engine) {
	c := &Controller{
		Service: *NewScheduleService(),
	}
	group := r.Group("/schedule")
	group.Use(AuthMiddleware())
	group.GET("/", func(c *gin.Context) {
		c.HTML(200, "schedule.html", nil)
	})
	group.GET("/browse", c.Browse)
	group.GET("/auth", func(c *gin.Context) {
		c.JSON(200, nil)
	})

	group.POST("/add_once", c.AddOnce)
	group.POST("/add_long", c.AddLong)
	
	group.DELETE("/", c.DeleteTask)
}

func (controller *Controller) AddOnce(c *gin.Context) {
	// AddOnce
	var date Date
	c.BindJSON(&date)
	err := controller.Service.AddOnce(date.Year, date.Month, date.Day, date.Hour, date.Minute, date.Content)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{})
		return
	}
	c.JSON(200, gin.H{})
}

func (controller *Controller) AddLong(c *gin.Context) {
	// AddLong
	var date Date
	c.BindJSON(&date)
	err := controller.Service.AddLong(date.Hour, date.Minute, date.Weekday, date.Content)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{})
		return
	}
	c.JSON(200, gin.H{})
}

func (controller *Controller) DeleteTask(c *gin.Context) {
	// DeleteTask
	var date Date
	c.BindJSON(&date)
	err := controller.Service.DeleteTask(date.Id)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{})
		return
	}
	c.JSON(200, gin.H{})
}

func (controller *Controller) Browse(c *gin.Context) {
	var date Date
	c.BindQuery(&date)
    // 从服务层获取所有任务
    tasks, err := controller.Service.GetAllTasks(date.Year, date.Month, date.Day, date.Weekday)
    if err != nil {
        log.Println("获取任务失败:", err)
        c.JSON(500, gin.H{"error": "无法获取任务列表"})
        return
    }
    // 渲染浏览页面，并传递任务数据，假设模板位于 templates/schedule/browse.html
    c.HTML(200, "schedule/browse.html", gin.H{
        "Tasks": tasks,
    })
}

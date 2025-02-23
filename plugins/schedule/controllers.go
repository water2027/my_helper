package schedule

import (
	"github.com/gin-gonic/gin"
)

type Controller struct {}

type ScheduleController interface {
	Login(c *gin.Context)
	AddOnce(c *gin.Context)
	AddLong(c *gin.Context)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		authCode := c.GetHeader("Authorization")
		if authCode != "1234" {
			c.JSON(400, gin.H{
				
			})
		}
	}	
}

func (sc *Controller) Login(r *gin.Engine) {

}


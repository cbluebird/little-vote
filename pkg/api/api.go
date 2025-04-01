package api

import (
	"github.com/gin-gonic/gin"
	"little-vote/pkg/dao"
	"little-vote/pkg/kafka"
	"little-vote/pkg/ticket"
	"log"
)

func Cas(c *gin.Context) {
	ticket.ServerTicket.RLock()
	defer ticket.ServerTicket.RUnlock()
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"ticket": ticket.ServerTicket.TicketId,
		},
	})
}

func Query(c *gin.Context) {
	name := c.Param("name")
	count, err := dao.GetUserInCache(name)
	if err != nil {
		user, err := dao.GetUserInfo(name)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "failed to get user info in database",
			})
			return
		}
		err = dao.SetUserInCache(user.Name)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "fail to set cache",
			})
			return
		}
		count = user.Count
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"count": count,
		},
	})
}

type Req struct {
	Name   []string `json:"name" binding:"required"`
	Ticket string   `json:"ticket" binding:"required"`
}

func Vote(c *gin.Context) {
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code": 1,
			"msg":  "fail to bind",
		})
		return
	}

	ticket.ServerTicket.RLock()
	if ticket.ServerTicket.TicketId != req.Ticket {
		c.JSON(400, gin.H{
			"code": 1,
			"msg":  "fail to check ticket",
		})
		ticket.ServerTicket.RUnlock()
		return
	}
	if ticket.ServerTicket.Count >= ticket.MAX {
		c.JSON(401, gin.H{
			"code": 2,
			"msg":  "over ticket limit",
		})
		ticket.ServerTicket.RUnlock()
		return
	}
	ticket.ServerTicket.Count++
	ticket.ServerTicket.RUnlock()

	for _, name := range req.Name {
		_, err := dao.GetUserInCache(name)
		if err != nil {
			err = dao.SetUserInCache(name)
			if err != nil {
				c.JSON(500, gin.H{
					"code": 1,
					"msg":  "fail to set cache",
				})
				return
			}
		}
		err = dao.IncrUserInCache(name)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 1,
				"msg":  "fail to incr",
			})
			return
		}
		err = kafka.Send(name)
		if err != nil {
			log.Println("fail to send message in kafka", err)
		}
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

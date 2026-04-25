package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/saintbyte/home-ctrl/internal/scheduler"
)

type TaskHandler struct {
	config *config.Config
	sched  *scheduler.Scheduler
}

func NewTaskHandler(cfg *config.Config, sched *scheduler.Scheduler) *TaskHandler {
	return &TaskHandler{
		config: cfg,
		sched:  sched,
	}
}

func (h *TaskHandler) SetupRoutes(router *gin.RouterGroup) {
	taskGroup := router.Group("/tasks")
	{
		taskGroup.GET("", h.listTasks)
		taskGroup.GET("/:name", h.getTask)
		taskGroup.POST("/:name/run", h.runTask)
		taskGroup.PATCH("/:name/enable", h.enableTask)
		taskGroup.PATCH("/:name/disable", h.disableTask)
	}
}

func (h *TaskHandler) listTasks(c *gin.Context) {
	tasks := make([]gin.H, 0, len(h.config.Tasks))
	for _, task := range h.config.Tasks {
		tasks = append(tasks, gin.H{
			"name":     task.Name,
			"schedule": task.Schedule,
			"enabled":  task.Enabled,
			"command":  task.Command,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"total": len(tasks),
	})
}

func (h *TaskHandler) getTask(c *gin.Context) {
	name := c.Param("name")

	for _, task := range h.config.Tasks {
		if task.Name == name {
			c.JSON(http.StatusOK, gin.H{
				"name":     task.Name,
				"schedule": task.Schedule,
				"enabled":  task.Enabled,
				"command":  task.Command,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error":   "Not Found",
		"message": "Task not found",
	})
}

func (h *TaskHandler) runTask(c *gin.Context) {
	name := c.Param("name")

	if h.sched != nil {
		err := h.sched.RunTask(name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Task triggered",
			"name":    name,
			"status":  "running",
		})
		return
	}

	for _, task := range h.config.Tasks {
		if task.Name == name {
			c.JSON(http.StatusOK, gin.H{
				"message": "Task triggered",
				"name":    name,
				"status":  "pending",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error":   "Not Found",
		"message": "Task not found",
	})
}

func (h *TaskHandler) enableTask(c *gin.Context) {
	name := c.Param("name")

	if h.sched != nil {
		err := h.sched.EnableTask(name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"name":    name,
			"enabled": true,
		})
		return
	}

	for i, task := range h.config.Tasks {
		if task.Name == name {
			h.config.Tasks[i].Enabled = true
			c.JSON(http.StatusOK, gin.H{
				"name":    name,
				"enabled": true,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error":   "Not Found",
		"message": "Task not found",
	})
}

func (h *TaskHandler) disableTask(c *gin.Context) {
	name := c.Param("name")

	if h.sched != nil {
		err := h.sched.DisableTask(name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Not Found",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"name":    name,
			"enabled": false,
		})
		return
	}

	for i, task := range h.config.Tasks {
		if task.Name == name {
			h.config.Tasks[i].Enabled = false
			c.JSON(http.StatusOK, gin.H{
				"name":    name,
				"enabled": false,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error":   "Not Found",
		"message": "Task not found",
	})
}

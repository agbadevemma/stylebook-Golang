package users

import "github.com/gin-gonic/gin"

type Module struct {
	handler *Handler
}

func NewModule(handler *Handler) *Module {
	return &Module{handler: handler}
}

func (m *Module) ModuleRoutes(api *gin.RouterGroup) {
	user := api.Group("/auth/user")
	{
		user.GET("/signup", m.handler.Signup)
		user.POST("/login", m.handler.Login)
	}

}

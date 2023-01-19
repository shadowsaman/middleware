package api

import (
	"errors"
	"net/http"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"

	_ "app/api/docs"
	"app/api/handler"
	"app/storage"
)

func NewApi(r *gin.Engine, storage storage.StorageI) {

	handlerV1 := handler.NewHandler(storage)

	r.Use(customCORSMiddleware())

	r.POST("/product", handlerV1.CreateProduct)
	r.GET("/product/:id", handlerV1.GetByIDProduct)
	r.GET("/product", handlerV1.GetListProduct)
	r.PUT("/product/:id", handlerV1.UpdateProduct)
	r.DELETE("/product/:id", handlerV1.DeleteProduct)

	r.POST("/category", handlerV1.CreateCategory)
	r.GET("/category/:id", handlerV1.GetByIdCategory)
	r.GET("/category", handlerV1.GetListCategory)
	r.DELETE("/category/:id", handlerV1.DeleteCategory)
	r.PUT("/category/:id", handlerV1.UpdateCategory)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}

func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if len(c.Request.Header["Auth"]) > 0 {
			key := c.Request.Header["Auth"][0]
			if key != "12345" {
				c.JSON(http.StatusUnauthorized, struct {
					Code int
					Err  string
				}{
					Code: http.StatusUnauthorized,
					Err:  errors.New("error access denied").Error(),
				})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, struct {
				Code int
				Err  string
			}{
				Code: http.StatusUnauthorized,
				Err:  errors.New("error access denied").Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Allow-Headers", "Platform-Id, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Max-Age", "3600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

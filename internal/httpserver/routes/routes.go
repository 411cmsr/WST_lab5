package routes

import (
	"WST_lab5_server/internal/database/postgres"
	"WST_lab5_server/internal/handlers"
	
	"github.com/gin-gonic/gin"
)

func Init(httpserver *gin.Engine) {

	//Восстановление после паники
	httpserver.Use(gin.Recovery())
	//Логгирование 
	httpserver.Use(gin.Logger())
	//Подключение к БД
	db := postgres.Init()
	storage := &postgres.Storage{DB: db}
	route := &handlers.StorageHandler{Storage: storage}
	//routes по запросам
	apiv1 := httpserver.Group("/api/v1")
	apiv1.GET("/persons", route.SearchPersonHandler)
	apiv1.POST("/persons", route.AddPersonHandler)
	apiv1.GET("/persons/list", route.GetAllPersonsHandler)
	apiv1.GET("/person/:id", route.GetPersonHandler)
	apiv1.PUT("/person/:id", route.UpdatePersonHandler)
	apiv1.DELETE("/person/:id", route.DeletePersonHandler)

}

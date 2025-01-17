package handlers

import (
	"WST_lab5_server/internal/database"
	"WST_lab5_server/internal/database/postgres"
	"WST_lab5_server/internal/models"
	"errors"
	"fmt"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	
)

/*
Структура обработчика для разделения логики обработки запросов от доступа к данным
*/
type StorageHandler struct {
	Storage *postgres.Storage
}

/*
//
Метод обработки запроса на поиск по всем полям данных
//
*/
func (sh *StorageHandler) SearchPersonHandler(context *gin.Context) {
	//Получаем строку поиска из запроса
	searchString := context.Query("query")
	//Проверяем что строка не пустая
	if searchString == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Search query cannot be empty."})
		return
	}
	//Ищем в базе данных по этой строке и возвращаем результат
	persons, err := sh.Storage.SearchPerson(searchString)
	if err != nil {
		//При ошибке возвращаем статус InternalServerError (500) и сообщение об ошибке
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve persons."})
		return
	}
	if len(persons) == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"code":    "not_found",
			"message": fmt.Sprintf("Person for '%s' request was not found.", searchString),
		})
		return
	}
	//Возвращаем статус ОК (200) и результат поиска (массив)
	context.JSON(http.StatusOK, persons)
}

/*
Метод обработки запроса на получение всех данных
*/
func (sh *StorageHandler) GetAllPersonsHandler(context *gin.Context) {
	//Получаем все данные из базы данных
	persons, err := sh.Storage.GetAllPersons()
	if err != nil {
		//При ошибке возвращаем статус InternalServerError (500) и сообщение об ошибке
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch persons. Try again later."})
		return
	}
	//Возвращаем статус ОК (200) и результат поиска (массив)
	context.JSON(http.StatusOK, persons)
}

/*
Метод обработки запроса на получение всех данных
*/
func (sh *StorageHandler) GetPersonHandler(context *gin.Context) {
	//Получаем id из запроса
	personIdStr := context.Param("id")
	//Преобразуем строку id в uint
	personId, err := strconv.ParseUint(personIdStr, 10, 64)
	if err != nil {
		//При ошибке возвращаем статус Bad Request (400) и сообщение об ошибке
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse person id."})
		return
	}
	//Получаем данные из базы данных
	person, err := sh.Storage.GetPerson(uint(personId))
	if err != nil {
		//При ошибке, проверяем тип ошибки
		if errors.Is(err, database.ErrPersonNotFound) {
			//Если запись не найдена, возвращаем статус  Not Found (404) и сообщение об ошибке
			context.JSON(http.StatusNotFound, gin.H{"message": "Person not found."})
			return
		}
		fmt.Println(err)
		//Если другая ошибка, возвращаем статус Internal Server Error (500) и сообщение об ошибке
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch person."})
		return
	}

	//Возвращаем статус ОК (200) и результат
	context.JSON(http.StatusOK, person)

}

/*
Метод обработки запроса на добавление новых данных
*/
func (sh *StorageHandler) AddPersonHandler(context *gin.Context) {
	var newPerson models.Person
	//Привязываем к структуре
	err := context.ShouldBindJSON(&newPerson)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}
	//Проверяем полученные данные
	//Возраст > 0
	if newPerson.Age < 0 {
		context.Error(gin.Error{
			Type: gin.ErrorTypePublic,
			Err:  fmt.Errorf("invalid age value"),
		})
		return
	}
	
	//Добавляем в БД
	id, err := sh.Storage.AddPerson(&newPerson)
	if err != nil {
		//Проверяем тип ошибки
		if errors.Is(err, database.ErrEmailExists) {
			//Если запись с таким email уже существует, возвращаем статус Conflict (409) и сообщение об ошибке
			context.JSON(http.StatusConflict, gin.H{"message": "Email already in use."})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create person."})
		return
	}
	//Возвращаем статус Created (201) и id новой записи
	context.JSON(http.StatusCreated, gin.H{"id": id})
}

/*
Метод обработки запроса на изменение данных
*/
func (sh *StorageHandler) UpdatePersonHandler(context *gin.Context) {
	//Получаем id из запроса
	personIdStr := context.Param("id")
	//Преобразуем строку id в uint
	personId, err := strconv.ParseUint(personIdStr, 10, 64)
	if err != nil {
		//При ошибке возвращаем статус Bad Request (400) и сообщение об ошибке
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse person id."})
		return

	}
	//Создаем структуру обновленных данных
	var updatedPerson models.Person
	//Привязываем данные из запроса к структуре
	err = context.ShouldBindJSON(&updatedPerson)
	if err != nil {
		//При ошибке возвращаем статус Bad Request (400) и сообщение об ошибке
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data."})
		return
	}
	//Присваиваем id обновляемой записи в структуре
	updatedPerson.ID = uint(personId)


	//Обновляем данные в базе данных
	err = sh.Storage.UpdatePerson(&updatedPerson)
	if err != nil {
		//Проверяем на отсутстве в базе данных
		if errors.Is(err, database.ErrPersonNotFound) {
			//Если запись не найдена, возвращаем статус Not Found (404) и сообщение об ошибке
			context.JSON(http.StatusNotFound, gin.H{"message": "Person not found."})
			return
		}
		//Если другая ошибка, возвращаем статус Internal Server Error (500) и сообщение об ошибке
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update person."})
		return
	}
	//Возвращаем статус OK (200) и сообщение об успешном обновлении данных
	context.JSON(http.StatusOK, gin.H{"message": "Person updated successfully!"})
}

/*
Метод обработки запроса на удаление данных
*/
func (sh *StorageHandler) DeletePersonHandler(context *gin.Context) {
	//Получаем id из запроса
	personIdStr := context.Param("id")
	//Преобразуем строку id в uint
	personId, err := strconv.ParseUint(personIdStr, 10, 64)
	if err != nil {
		//При ошибке возвращаем статус Bad Request (400) и сообщение об ошибке
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse person id."})
		return

	}
		
	//Удаляем запись из базы данных
	err = sh.Storage.DeletePerson(&models.Person{ID: uint(personId)})
	if err != nil {
		//При ошибке возвращаем статус Internal Server Error (500) и сообщение об ошибке
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could Not Delete Person"})
		return
	}
	//Возвращаем статус OK (200) и сообщение об успешном удалении данных
	context.JSON(http.StatusOK, gin.H{"message": "Deleted Successfully"})

}
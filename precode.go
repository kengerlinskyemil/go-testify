package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4 // максимальное количество кафе, которое сервис может вернуть

	// Создаём запрос с параметрами count=10 и city=moscow
	req, err := http.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
	require.NoError(t, err, "Ошибка при создании запроса") // Используем require, чтобы сразу остановить тест в случае ошибки

	// Записываем ответ
	responseRecorder := httptest.NewRecorder()

	// Создаём обработчик
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка, что статус-код ответа - 200
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Неверный статус-код ответа")

	// Проверка, что в ответе не пустое тело
	assert.NotEmpty(t, responseRecorder.Body.String(), "Тело ответа пустое")

	// Проверка, что количество кафе не больше доступного
	cafesReturned := strings.Split(responseRecorder.Body.String(), ",")
	assert.Len(t, cafesReturned, totalCount, "Количество кафе в ответе неверное")
}

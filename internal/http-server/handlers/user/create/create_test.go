package create_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"tstUser/internal/http-server/handlers/user/create"
	"tstUser/internal/http-server/handlers/user/create/mocks"
	"tstUser/internal/http-server/transport/userDTO"
	"tstUser/internal/lib/logger/handlers/slogdiscard"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		nameTest  string // Имя теста
		usr       userDTO.UDTO
		respError string // Какую ошибку мы должны получить?
		mockError error  // Ошибку, которую вернёт мок
	}{
		{
			nameTest: "Success",
			usr: userDTO.UDTO{
				Name:    "Dmitrii",
				Surname: "Shilenko",
				Mail:    "Prettypepe@mail.ru",
				Cash:    100000000,
				Date:    "20011007",
			},
			// Тут поля respError и mockError оставляем пустыми,
			// т.к. это успешный запрос
		},
		// Другие кейсы ...
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.nameTest, func(t *testing.T) {
			t.Parallel()
			// Создаем объект мока стораджа
			urlSaverMock := mocks.NewUserCreator(t)

			// Если ожидается успешный ответ, значит к моку точно будет вызов
			// Либо даже если в ответе ожидаем ошибку,
			// но мок должен ответить с ошибкой, к нему тоже будет запрос:
			if tc.respError == "" || tc.mockError != nil {
				// Сообщаем моку, какой к нему будет запрос, и что надо вернуть
				urlSaverMock.On("CreateUser", tc.usr.Name, tc.usr.Surname, tc.usr.Mail, tc.usr.Date, tc.usr.Cash).
					Return(int64(1), tc.mockError).
					Once() // Запрос будет ровно один
			}

			// Создаем наш хэндлер
			handler := create.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

			// Формируем тело запроса
			input := fmt.Sprintf(`{"name": "%s", "surname": "%s", "mail": "%s", "cash": %v, "date": "%s"}`,
				tc.usr.Name, tc.usr.Surname, tc.usr.Mail, tc.usr.Cash, tc.usr.Date)

			// Создаем объект запроса
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			// Создаем ResponseRecorder для записи ответа хэндлера
			rr := httptest.NewRecorder()
			// Обрабатываем запрос, записывая ответ в рекордер
			handler.ServeHTTP(rr, req)

			// Проверяем, что статус ответа корректный
			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp create.Response

			// Анмаршаллим тело, и проверяем что при этом не возникло ошибок
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			// Проверяем наличие требуемой ошибки в ответе
			require.Equal(t, tc.respError, resp.Error)

			// Другие проверки
		})
	}
}

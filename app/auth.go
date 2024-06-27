package app

import (
	"context"
	//"fmt"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var JWTAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonAuth := []string{"/api/user/new", "/api/user/login"} //Список эндпоинтов, для которых не требуется авторизация
		requestPath := r.URL.Path                               //текущий путь запроса

		//проверяем, не требует ли запрос аутентификации, обслуживаем запрос, если он не нужен
		for _, value := range nonAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" { //токен пустой. Возвращаем (403)("UnAuthorized")
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)                //доступ запрещён
			w.Header().Add("Content-Type", "application/json") //устанавливает заголовок ответа "Content-Type" со значением "application/json".
			//четко сообщает клиенту, что ответ содержит данные в формате JSON.
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		//fmt.Print(splitted)
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Неправильный токен, как правило, возвращает 403 http-код
			response = u.Message(false, "Token is not valid")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid { //токен недействителен, возможно, не подписан на этом сервере
			response = u.Message(false, "Token is not valid")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		//Всё прошло хорошо, продолжаем выполнение запроса
		//fmt.Sprintf("User %/n", tk.UserId) //Полезно для мониторинга
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //передать управление следующему обработчику!

	})
}

/*
JWTAuthentication является middleware-функцией,
которая обрабатывает каждый входящий HTTP-запрос перед передачей его следующему обработчику в цепочке.
*/

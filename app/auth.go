package app

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/shoksin/go-contacts-REST-API-/models"
	u "github.com/shoksin/go-contacts-REST-API-/utils"

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

		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" { //токен пустой. Возвращаем (403)("UnAuthorized")
			w.WriteHeader(http.StatusForbidden)                //доступ запрещён
			w.Header().Add("Content-Type", "application/json") //устанавливает заголовок ответа "Content-Type" со значением "application/json".
			//сообщает клиенту, что ответ содержит данные в формате JSON.
			u.Respond(w, u.Message(false, "Missing auth token"))
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, u.Message(false, "Invalid auth token"))
			return
		}

		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Неправильный токен, как правило, возвращает 403 http-код
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, u.Message(false, "Token is not valid"))
			return
		}

		if !token.Valid { //токен недействителен, возможно, не подписан на этом сервере
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, u.Message(false, "Token is not valid"))
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

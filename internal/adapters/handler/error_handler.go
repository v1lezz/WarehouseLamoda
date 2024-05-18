package handler

import (
	"encoding/json"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	ans := map[string]interface{}{
		"data":  nil,
		"error": err.Error(),
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ans)
	//не обрабатываю ошибку тут, т.к. может возникнуть случай,
	//в котором у нас какая-то часть уже записалась в w (т.к. запись потоковая)
	//и чуть позже вознилка ошибка
	//когда в http.ResponseWriter хоть что-то записывается, автоматически
	//подставляется статус 200 - OK, что не позволит нам отследить ошибку на
	//стороне клиента и он получит битые данные
	//в данном случае есть 2 варианта:
	//1) использовать json.Marshal, если нам действительно важно отслеживать ошибку
	//2) не ловить ошибку
	//я выбрал второй вариант, поскольку вероятность прихода енкодеру битых данных очень мала
	//поскольку используется map[string]interface{}
}

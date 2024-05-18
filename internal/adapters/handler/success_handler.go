package handler

import (
	"encoding/json"
	"net/http"
)

func SuccessHandler(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	ans := map[string]interface{}{
		"data":  data,
		"error": nil,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ans)
	//объяснение того, почему не обрабатываю ошибку тут такое же как в error_handler.go в этом же пакете
}

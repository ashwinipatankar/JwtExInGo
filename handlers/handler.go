package handler

import "net/http"

var GetLoginPageHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/login.html")
})

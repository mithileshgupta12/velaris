package helper

import "net/http"

const (
	path       = "/"
	isHttpOnly = true
	sameSite   = http.SameSiteLaxMode
)

func SetCookie(w http.ResponseWriter, name, value string, maxAge int, isSecure bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		MaxAge:   maxAge,
		HttpOnly: isHttpOnly,
		Secure:   isSecure,
		SameSite: sameSite,
	}

	http.SetCookie(w, cookie)
}

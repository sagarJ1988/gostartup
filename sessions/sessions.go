package sessions

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

type SessionData struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
}

var CookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func Setsession(sess interface{}, w http.ResponseWriter) {
	s, _ := sess.(*SessionData)
	value := map[string]string{
		"FirstName": s.FirstName,
		"LastName":  s.LastName,
		"Email":     s.Email,
	}
	if encoded, err := CookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}

	//fmt.Println("Email \t" + s.Email + " \tfirst name \t " + s.FirstName + "\tlast name\t " + s.LastName)
}

func ClearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func Get(r *http.Request) {

}

func GetAll(r *http.Request) (string, string, string) {
	var FirstName, LastName, Email string
	if cookie, err := r.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			FirstName = cookieValue["FirstName"]
			LastName = cookieValue["LastName"]
			Email = cookieValue["Email"]
		}
	}
	return FirstName, LastName, Email
}

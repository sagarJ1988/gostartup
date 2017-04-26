package main

import (
	"crypto/md5"
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"

	"fmt"

	"gostartup/sessions"
	"gostartup/validation"

	_ "github.com/go-sql-driver/mysql"
)

/*
	Must() : It is intended for use in variable initializations such as tmpl
	ParseGlob() : ParseGlob creates a new Template and parses the template definitions
	from the files identified by the pattern,
	which must match at least one file
	ParseFiles() : ParseFiles creates a new Template and parses the template definitions from
	the named files.The returned template's name will have the (base) name and (parsed)
	contents of the first file
*/

var tmpl = template.Must(template.ParseGlob("templates/*"))
var db *sql.DB
var err error

type UserData struct {
	Fname string
	Lname string
	Email string
}

func main() {
	log.Println("server started on: http://localhost:9000")

	db, err = sql.Open("mysql", "root:password@/gostartup")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	err = db.Ping()

	if err != nil {
		panic(err.Error())
	}

	//http.Handle("/resources/", http.FileServer(http.Dir("resources")))
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	//url management
	http.HandleFunc("/", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/dashboard", Dashboard)

	http.HandleFunc("/user", User)
	//http.HandleFunc("/form-validation", formValidation)
	http.HandleFunc("/listing", Listing)
	http.HandleFunc("/profile", Profile)
	//start server
	http.ListenAndServe(":9000", nil)
	// :9000 is address  nil  is handler
	/* it used for https connection
	cert.pem is certificate file
	key.pem is matching key for the server must be provided
	http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	*/
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessions.ClearSession(w)
	http.Redirect(w, r, "/", 301)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// Join each row on struct inside the Slice

	Fname, Lname, Email := sessions.GetAll(r)
	log.Println("session data ", Fname+" "+Lname+" "+Email)

	if r.Method != "POST" {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	msg := &validation.Message{
		Email:    r.FormValue("userName"),
		Password: r.FormValue("password"),
	}

	if msg.ValidateLogin() == false {
		tmpl.ExecuteTemplate(w, "login.html", msg)
		return
	}

	username := r.FormValue("userName")
	password := r.FormValue("password")

	h := md5.New()
	io.WriteString(h, password)

	encPassword := fmt.Sprintf("%x", h.Sum(nil))

	var databaseUsername string
	var databasePassword string
	var dataID int
	var FirstName string
	var LastName string

	err := db.QueryRow("select id, firstName,lastname, userName,password from user where userName=?", username).Scan(&dataID, &FirstName, &LastName, &databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(w, r, "/", 301)
		return
	}

	//log.Println("dbpass : " + databasePassword + " \tencpassword : " + encPassword)

	if databasePassword != encPassword {
		log.Println("Password not matched")
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	sess := &sessions.SessionData{dataID, FirstName, LastName, databaseUsername}

	sessions.Setsession(sess, w)

	http.Redirect(w, r, "/dashboard", 301)

	//tmpl.ExecuteTemplate(w, "login.html", nil)
}

func Signup(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
		return
	}

	msg := &validation.SignupForm{
		FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
		UserName:  r.FormValue("userName"),
		Password:  r.FormValue("password"),
	}

	if msg.ValidateSignup() == false {
		//log.Println("errors ", msg)
		tmpl.ExecuteTemplate(w, "login.html", msg)
		return
	}

	firstName := r.FormValue("fname")
	lastName := r.FormValue("lname")
	userName := r.FormValue("userName")
	password := r.FormValue("password")

	h := md5.New()
	io.WriteString(h, password)

	//newPass :=fmt.Printf("%x", md5.Sum(data))
	newPassword := fmt.Sprintf("%x", h.Sum(nil))
	//newPassword := h.Sum(nil)

	log.Println("new password : ", newPassword)

	// Prepare sql insert
	insForm, err := db.Prepare("INSERT INTO user(firstName, lastname,userName,password) VALUES(?,?,?,?)")

	if err != nil {
		panic(err.Error())
	}

	insForm.Exec(firstName, lastName, userName, newPassword)

	log.Println("Inserted Successfully")

	http.Redirect(w, r, "/", 301)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	FirstName, LastName, Email := sessions.GetAll(r)
	log.Println("session data ", FirstName+" "+LastName+" "+Email)
	u := &UserData{FirstName, LastName, Email}
	tmpl.ExecuteTemplate(w, "index.html", u)
}

func User(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "user.html", nil)
}

// func formValidation(w http.ResponseWriter, r *http.Request) {
// 	tmpl.ExecuteTemplate(w, "form_validation.html", nil)
// }

func Listing(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "listing.html", nil)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "profile.html", nil)
}

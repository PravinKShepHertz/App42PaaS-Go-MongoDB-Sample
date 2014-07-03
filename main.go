package main

import (
	"fmt"
	"net/http"
	"log"
	"html/template"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	db *mgo.Session
)

type User struct {
	Name string
	Email string
	Description string
}

func setupDB() *mgo.Session{
	db, err := mgo.Dial("ax1j4z8drbya3zg2:agcne68j8xrzot71oyq37nji8dbv3973@192.168.3.241:4841/demo_db")
	PanicIf(err)
	return db
}

func PanicIf(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request ) {
	fmt.Println("In indexHandler")
	coll := db.DB("demo_db").C("user")
	fmt.Println("Collection: ", coll)

	users := []User{}
	err := coll.Find(bson.M{}).All(&users)
	PanicIf(err)
	fmt.Println(users)

	t := template.New("index.html")
	t.ParseFiles("templates/index.html")
	t.Execute(w, users)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In newHanlder")
	t := template.New("new.html")
	t.ParseFiles("templates/new.html")
	t.Execute(w, t)
}

func saveHandler (w http.ResponseWriter, r *http.Request) {
	fmt.Println("In saveHandler")
	name := r.FormValue("name")
	email := r.FormValue("email")
	description := r.FormValue("description")

	coll := db.DB("demo_db").C("user")
	fmt.Println("Collection: ", coll)

	err := coll.Insert(&User{name, email, description})
	PanicIf(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	db = setupDB()
	fmt.Println("Database connected")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/save/", saveHandler)

	fmt.Println("Listening on port 3000......")
	http.Handle("/public/css/", http.StripPrefix("/public/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/public/images/", http.StripPrefix("/public/images/", http.FileServer(http.Dir("public/images"))))
	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Error to listen 3000: ", err)
	}

}

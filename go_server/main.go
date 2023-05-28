package main

import (
	"fmt"
	"log"
	"net/http"
)

func formHandler(w http.ResponseWriter, r *http.Request){
	if err := r.ParseForm(); err!=nil {
		fmt.Fprintf(w, "ParseForm() err: %v",err)
		return
	}
	fmt.Fprintf(w, "POST request successful")
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s\n", address)
}
func helloHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/hello" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != "GET"{
		http.Error(w,"method not supported", http.StatusNotFound)
		return 
	}
	fmt.Fprintf(w,"hello!")
}
func main(){
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/",fileServer)
	http.HandleFunc("/form",formHandler)
	http.HandleFunc("/hello",helloHandler)
	fmt.Println("Starting server")
	if err := http.ListenAndServe(":8084",nil); err!=nil {
		log.Fatal(err)
	}
}



// package main

// import (
//     "fmt"
//     "html/template"
//     // "log"
//     "net/http"
//     "strings"
// )

// func sayhelloName(w http.ResponseWriter, r *http.Request) {
//     r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
//     // attention: If you do not call ParseForm method, the following data can not be obtained form
//     fmt.Println(r.Form) // print information on server side.
//     fmt.Println("path", r.URL.Path)
//     fmt.Println("scheme", r.URL.Scheme)
//     fmt.Println(r.Form["url_long"])
//     for k, v := range r.Form {
//         fmt.Println("key:", k)
//         fmt.Println("val:", strings.Join(v, ""))
//     }
//     fmt.Fprintf(w, "Hello astaxie!") // write data to response
// }

// func login(w http.ResponseWriter, r *http.Request) {
//     fmt.Println("method:", r.Method) //get request method
//     if r.Method == "GET" {
//         t, err := template.ParseFiles("login.html")
//         t.Execute(w, nil)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
//     } //else {
//     //     r.ParseForm()
//     //     // logic part of log in
//     //     fmt.Println("username:", r.Form["username"])
//     //     fmt.Println("password:", r.Form["password"])
//     // }
// }

// // func main() {
// //     http.HandleFunc("/", sayhelloName) // setting router rule
// //     http.HandleFunc("/login", login)
// //     err := http.ListenAndServe(":9090", nil) // setting listening port
// //     if err != nil {
// //         log.Fatal("ListenAndServe: ", err)
// //     }
// // }


// type ContactDetails struct {
// 	Email string
// 	Subject string
// 	Message string
// }
// func main() {
// 	tmpl := template.Must(template.ParseFiles("forms.html"))
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			tmpl.Execute(w, nil)
// 			return
// 		}
// 		details := ContactDetails{
// 			Email: r.FormValue("email"),
// 			Subject: r.FormValue("subject"),
// 			Message: r.FormValue("message"),
// 		}
// 		_ = details
// 		tmpl.Execute(w, struct{Success bool}{true})

// 	})
// 	http.ListenAndServe(":8080", nil)
// }
package main

import (
	"errors"
	"log"
	"net/http"
)

func clearCookieHandler(w http.ResponseWriter, r *http.Request) {

	// MaxAge=0 means no 'Max-Age' attribute specified.

	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'

	// MaxAge>0 means Max-Age attribute present and given in seconds

	c := &http.Cookie{
		Name:     "exampleCookie",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, c)

	w.Write([]byte("Cookies Cleard Now!"))
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize a new cookie containing the string "Hello Shubham Mishra!" and some
	// non-default attributes.
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "Hello Shubham Mishra!",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	// Use the http.SetCookie() function to send the cookie to the client.
	// Behind the scenes this adds a `Set-Cookie` header to the response
	// containing the necessary cookie data.

	log.Println(cookie.String())
	log.Println(cookie.Valid())
	log.Println(cookie.Value)

	http.SetCookie(w, &cookie)

	// Write a HTTP response as normal.
	w.Write([]byte("cookie set!"))
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the cookie from the request using its name (which in our case is
	// "exampleCookie"). If no matching cookie is found, this will return a
	// http.ErrNoCookie error. We check for this, and return a 400 Bad Request
	// response to the client.
	cookie, err := r.Cookie("exampleCookie")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	// Echo out the cookie value in the response body.
	w.Write([]byte(cookie.Value))
}

func main() {
	// Start a web server with the two endpoints.
	mux := http.NewServeMux()

	mux.HandleFunc("/set", setCookieHandler)

	mux.HandleFunc("/get", getCookieHandler)

	mux.HandleFunc("/clear", clearCookieHandler)

	log.Print("Listening...")

	err := http.ListenAndServe(":3000", mux)

	if err != nil {
		log.Fatal(err)
	}
}

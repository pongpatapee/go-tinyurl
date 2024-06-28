package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const HOST string = "http://localhost:8080"

type URLShortener struct {
	urls map[string]string
}

func (us *URLShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "Missing url", http.StatusBadRequest)
		return
	}

	// Generate unique key for this URL
	shortKey := generateShortKey()
	us.urls[shortKey] = url

	shortUrl := fmt.Sprintf("%v/short/%v", HOST, shortKey)

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
        <h2>URL Shortener</h2>
        <p>Original URL: %s</p>
        <p>Shortened URL: <a href="%s">%s</a></p>
        <form method="post" action="/shorten">
            <input type="text" name="url" placeholder="Enter a URL">
            <input type="submit" value="Shorten">
        </form>
        `, url, shortUrl, shortUrl)

	fmt.Fprint(w, responseHTML)
}

func (us *URLShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]

	if shortKey == "" {
		http.Error(w, "Short key missing", http.StatusBadRequest)
		return
	}

	originalUrl, found := us.urls[shortKey]

	if !found {
		http.Error(w, "Shortened URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusMovedPermanently)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	shortKeyArr := make([]byte, keyLength)

	for i := range keyLength {
		shortKeyArr[i] = charset[r.Intn(len(charset))]
	}

	return string(shortKeyArr)
}

func HandleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}

	// Serve the HTML form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<form method="post" action="/shorten">
				<input type="url" name="url" placeholder="Enter a URL" required>
				<input type="submit" value="Shorten">
			</form>
		</body>
		</html>
	`)
}

func main() {
	shortener := &URLShortener{
		urls: make(map[string]string),
	}

	http.HandleFunc("/", HandleForm)
	http.HandleFunc("/shorten", shortener.HandleShorten)
	http.HandleFunc("/short/", shortener.HandleRedirect)

	fmt.Println("Listening at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

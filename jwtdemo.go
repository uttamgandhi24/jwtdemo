package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func userpassword() (userPasswordMap map[string]string) {
	userPasswordMap = make(map[string]string)
	userPasswordMap["Synerzip1"] = "Password1"
	userPasswordMap["Synerzip2"] = "Password2"
	userPasswordMap["Synerzip3"] = "Password3"
	return
}

func getPage(pageNum int) interface{} {

	pages := []string{"This is page 1",
		"This is page 2",
		"This is page 3"}

	if pageNum > len(pages) {
		log.Fatal("Invalid PageNum")
	}
	return pages[pageNum]
}

func getPageHandler(w http.ResponseWriter, r *http.Request) {
	if !authenticateJWT(r.Header) {
		fmt.Println("StatusUnauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	pageNum, _ := strconv.Atoi(vars["pageNum"])
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println(getPage(pageNum))
	if err := json.NewEncoder(w).Encode(getPage(pageNum)); err != nil {
		panic(err)
	}
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	authHeader := r.Header.Get("Authentication")
	if len(authHeader) == 0 || !authenticateRequest(authHeader) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := generateJWT()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(token); err != nil {
		panic(err)
	}
}

func generateJWT() (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 1).Unix()})

	// Sign and get the complete encoded token as a string
	tokenString, err = token.SignedString([]byte("secret"))
	return tokenString, err
}

func authenticateRequest(authHeader string) bool {

	data, err := base64.StdEncoding.DecodeString(authHeader)
	if err != nil {
		fmt.Println("error:", err)
		return false
	}
	fmt.Printf("%q\n", data)
	userpwd := strings.Split(string(data), ":")

	userpwdmap := userpassword()
	if userpwdmap[userpwd[0]] == userpwd[1] {
		return true
	}
	return false
}

func authenticateJWT(header http.Header) bool {
	jwtString := header.Get("Authentication")
	if len(jwtString) == 0 {
		return false
	}
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err == nil && token.Valid {
		return true
	} else {
		return false
	}
	return false
}

func main() {
	h := mux.NewRouter()
	h.HandleFunc("/book/page/{pageNum:[0-9][0-9]*}", getPageHandler).Methods("GET")
	h.HandleFunc("/authenticate", authenticateHandler).Methods("POST")

	fmt.Println("Listening on 3080....")
	http.ListenAndServe(":3080", h)
}

package main

import (
  "fmt"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "log"
  "os"
  "time"
  "encoding/json"
  "encoding/base64"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/dgrijalva/jwt-go"
  "strings"
  "strconv"
)

func userpassword() (userPasswordMap map[string]string) {
  userPasswordMap = make(map[string]string)
  userPasswordMap["Synerzip1"]="Password1"
  userPasswordMap["Synerzip2"]="Password2"
  userPasswordMap["Synerzip3"]="Password3"
  return
}

func connect() (session *mgo.Session) {
  connectURL := "localhost"
  session, err := mgo.Dial(connectURL)
  if err != nil {
    fmt.Printf("Can't connect to mongo, go error %v\n", err)
    os.Exit(1)
  }
  session.SetSafe(&mgo.Safe{})
  return session
}

func getPage(pageNum int) (result interface{}) {
  session := connect()
  defer session.Close()

  collection := session.DB("books").C("book")

  err := collection.Find(bson.M{"PageNum": pageNum}).One(&result)
  if err != nil {
    fmt.Printf("Error in find")
    log.Fatal(err)
  }
  return result
}

func getPageHandler(w http.ResponseWriter, r *http.Request) {
  if !authenticateJWT(r.Header) {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }
  vars := mux.Vars(r)
  pageNum, _ := strconv.Atoi(vars["pageNum"])
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if err := json.NewEncoder(w).Encode(getPage(pageNum)); err != nil {
    panic(err)
  }
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if !authenticateRequest(r.Header) {
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

func generateJWT()(tokenString string, err error) {
  token := jwt.New(jwt.SigningMethodHS256)
    // Set some claims
    token.Claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
    // Sign and get the complete encoded token as a string
    tokenString, err = token.SignedString([]byte("secret"))
    return tokenString, err
}

func authenticateRequest(header map[string][]string) bool {
  if header["Authentication"] == nil {
    return false
  }
  encodedValue := header["Authentication"][0]

  data, err := base64.StdEncoding.DecodeString(encodedValue)
  if err != nil {
    fmt.Println("error:", err)
    return false
  }
  //fmt.Printf("%q\n", data)
  userpwd := strings.Split(string(data),":")

   userpwdmap := userpassword()
   if userpwdmap[userpwd[0]] == userpwd[1] {
    return true
  }
  return false
}

func authenticateJWT(header map[string][]string) bool {
  if header["Authentication"] == nil {
    return false
  }

  jwtString := header["Authentication"][0]
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
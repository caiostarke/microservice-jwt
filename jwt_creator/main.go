package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var MySigninKey = []byte(os.Getenv("SECRET_KEY"))

func GetJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Kaer"
	claims["aud"] = "billing,wjs.dsa"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(MySigninKey)

	if err != nil {
		_ = fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil

}

func Index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJWT()
	if err != nil {
		fmt.Println("Failed to generate the token")
	}
	fmt.Println(validToken)
	fmt.Fprintf(w, string(validToken))
}

func handleRequests() {
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}

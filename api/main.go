package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

var MySigninKey = []byte(os.Getenv("SECRET_KEY"))

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Super Secret information")
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf(("invalid signing method"))
				}

				aud := "billing,wjs.dsa"
				iss := "jwtgo.io"

				checkAudience := t.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAudience {
					return nil, fmt.Errorf(("invalid aud"))
				}

				checkIssuer := t.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIssuer {
					return nil, fmt.Errorf(("invalid issuer"))
				}

				return MySigninKey, nil
			})

			if err != nil {
				fmt.Fprintln(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "No authorization token provided =( ")
		}
	})
}

func handleRequests() {
	http.Handle("/", isAuthorized(homePage))
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	handleRequests()
}

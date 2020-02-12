package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT() (string, error) {
	str := "keeek"
	mySigningKey := []byte(str)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["_id"] = "3"
	claims["name"] = "Evgsol"
	claims["score"] = 20

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func main() {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println("Failed to generate token")
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://ec2-52-91-188-222.compute-1.amazonaws.com:8000/add", nil)
	req.Header.Set("Token", validToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf(string(body))
}

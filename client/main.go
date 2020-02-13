package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(id, name string, score int) (string, error) {
	str := "keeek"
	mySigningKey := []byte(str)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["_id"] = id
	claims["name"] = name
	claims["score"] = score

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func call(id, name string, score int) error {
	validToken, err := GenerateJWT(id, name, score)
	if err != nil {
		return fmt.Errorf("Failed to generate token: %s", err)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://ec2-52-91-188-222.compute-1.amazonaws.com:8000/add", nil)
	req.Header.Set("Token", validToken)
	_, err = client.Do(req)
	return err
}


var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    return string(b)
}

func randStr() string {
	return RandStringRunes(10)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	type score struct {
		name  string
		score int
	}
	for i := 0; i < 10; i++ {
		err := call(randStr(), fmt.Sprintf("Docker%d", i), i * 10 + 30)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
	}
}

package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/atakanozceviz/cpypst-secure/model"
	"github.com/dgrijalva/jwt-go"
)

var SecretKey []byte

func Sign(data model.Data) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Action":  data.Action,
		"Content": data.Content,
		"From":    data.From,
		"Time":    data.Time,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(SecretKey)

	return tokenString, err
}

func Parse(ts string) (jwt.MapClaims, error) {
	// sample token string taken from the New example
	tokenString := ts

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
	return nil, nil
}

func EncSend(action, from, content, ip string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Duration(time.Second * 5),
	}

	enc, err := Sign(model.Data{Action: action, From: from, Content: content, Time: time.Now().Format(time.UnixDate)})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "http://"+ip+":8080/"+action, bytes.NewBuffer([]byte(enc)))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

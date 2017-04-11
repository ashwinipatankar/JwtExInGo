package authentication

import (
	"crypto/rsa"
	"io/ioutil"
	"log"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

//RSA Keys and Initialisation

const (
	privateKeyPath = "key/app.rsa"
	publicKeyPath  = "key/app.rsa.pub"
)

var (
	VerifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
)

func GetSignKey() *rsa.PrivateKey {
	return SignKey
}
func InitKeys() {
	signKeyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Error reading private key file")
		return
	}

	SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signKeyBytes)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	verifyKeyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal("error reading public key file")
	}

	VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyKeyBytes)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}
}

//MiddleWare
var ValidateToken = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return VerifyKey, nil
	},
	SigningMethod: jwt.SigningMethodRS256,
})

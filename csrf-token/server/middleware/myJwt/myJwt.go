package myjwt

import (
	"io/ioutil"
	"time"
	"errors"
	"log"
	"crypto/rsa"
	"github.com/samridhi/golang-csrf/db/models"
	"github.com/samridhi/golang-csrf/db"
	"github.com/dgrijalva/jwt-go"
	
)
const(
	privKeyPath = "keys/app.rsa"
	pubKeyPath = "keys/app.rsa.pub"
)
func InitJWT() error{
  signBytes, err := ioutil.ReadFile(privKeyPath)
  if err != nil{
	return err
  }
  signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
  if err != nil{
	return err
  }
  verifyBytes, err := ioutil.ReadFile(pubKeyPath)
  if err != nil{
	return err
  }
  verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
  if err != nil{
	return err
  }
  return nil
}

func CreateNewToken(uuid string, role string)(authTokenString, refreshTokenString, csrfSecret string, err error){
  //generating csrf secret
  csrfSecret, err  = models.GenerateCSRFSecret()
  if err != nil{
	return 
  }
  //generating refresh token
  refreshTokenString, err = createRefreshTokenString(uuid, role, csrfSecret)
  authTokenString, err = createAuthTokenString(uuid, role, csrfSecret)
  if err != nil{
      return
  }
  return
  //generating auth token
}

func CheckAndRefreshTokens(oldAuthTokenString string, oldRefreshTokenString string, oldCSRFToken string)(newAuthTokenString, newRefreshTokenString, newCSRFSecret, err error){
   if oldCSRFToken == ""{
	log.Println("No CSRF Token")
	err = errors.New("Unauthorised")
	return
   }
   jwt.ParseWithClaims(oldAuthTokenString, &models.TokenClaims{}, func( oken *jwt.Token)()  {
	
   }
}

func createAuthTokenString(uuid string, role string, csrfSecret string)(authTokenString, err error){
    authTokenExp := time.Now().Add(models.AuthToken).Unix()
	authClaims := models.TokenClaims{
		jwt.StandardClaims{
			Subject: uuid,
			ExpiresAt: authTokenExp,
		},
		role,
		csrfSecret,
	}
    authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"),authClaims)
	authTokenString, err = authJwt.SignedString(signKey)
	return

}

func createRefreshTokenString(uuid string, role string, csrfString string)(refreshTokenString string, err error){
  refreshTokenExp := time.Now().Add(models.RefreshTokenValidTime).Unix()
  refreshJti, er := db.StoreRefreshToken()
  if err != nil{
	return
  }
  refreshClaims := models.TokenClaims{
	jwt.StandardClaims{
		Id: refreshJti,
		Subject: uuid,
		ExpiresAt: refreshTokenExp,
	},
    role,
	csrfString,
  }
  refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)
  refreshTokenString := refreshJwt.SignedString(signKey)
  return
}

func updateRefreshTokenExp()(){

}

func updateAuthTokenString()(){

}

func RevokeRefreshToken() error{

}

func updateRefreshTokenCsrf()(){

}

func GrabUUid()(){

}
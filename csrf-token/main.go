package main

import "log"
var host = "localhost"
var port = "9010"
func main(){
  db.InitDB()
  jwtErr := myJwt.InitJWT()
  if jwtErr != nil{
	log.Println("Error initialising the JWT!")
	log.Fatal(jwtErr)
  }
  serverErr := server.StartServer(host,port)
  if serverErr != nil{
	log.Println("Error Starting server")
	log.Fatal(serverErr)
  }
}
package db
import(
   "github.com/"
)

var users = map[string]models.User{}

func InitDB(){

}

func StoreUser(username string, password string, role string)(uuid string, err error){

}

func FetchUserByUsername(username string)(models.User, string, error){
   for k,v := range users{
	if v.username == username{
		return v,k,nil
	}
   }
   return models.User{}, "", errors.New("User not foumd that matches given username")
}

func deleteuser(){

}

func FetchUserById()(){

}

func StoreRefreshToken()(){

}

func DeleteRefreshToken(){

}

func CheckRefreshToken() bool{

}

func LogUserIn()(){

}

func generateBcryptHash()(){

}

func checkPasswordAgainstHash() error{

}
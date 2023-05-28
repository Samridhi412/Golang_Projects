package middleware

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/justinas/alice"
)

func NewHandler() http.Handler{
 return alice.New(recoverHandler, authHandler).ThenFunction(logicHandler)
}

//prevent from crash
func recoverHandler(next http.Handler) http.Handler{
  fn := func (w http.ResponseWriter, r *http.Request)  {
	defer func(){
      if err := recover(); err != nil{
		log.Panic("Recovered ! Panic: %+v", err)
		http.Error(w, http.StatusText(500), 500)
	  }
	}()
	next.ServeHTTP(w,r)
  }
  return http.HandlerFunc(fn)
}

// case 1 check for jwt, whem logout cookies are deleted
func authHandler(next http.Handler) http.Handler{
   fn := func (w http.ResponseWriter, r *http.Request)  {
	switch r.URL.Path{
	case "/restricted", "/logout", "/deleteuser":
		log.Println("In auth restricted section")
		AuthCookie, authErr := r.Cookie("AuthToken")
		if authErr == http.ErrNoCookie{
			log.Println("Unauthorised  attempt! no auth cookie")
			nullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(401), 401)
			return
		}else if authErr != nil{
			log.Panic("panic: %+v",authErr)
			nullifyTokenCookies(&w, r)
			http.Error(w,http.StatusText(500), 500)
			return
		}
		RefreshCookie, refreshErr := r.Cookie("RefreshToken")
		if refreshErr == http.ErrNoCookie{
			log.Println("Unuthorised attempt! no refresh cookie found")
			nullifyTokenCookies(&w, r)
			http.Redirect(w, r, "/login", 302)
			return
		}else if refreshErr != nil{
			log.Panic("panic: %+v", refreshErr)
			nullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		requestCSRFToken := grabCSRFFromReq(r)
		log.Println(requestCSRFToken)
		authTokenString, refreshTokenString, csrfSecret, err := myJwt.CheckAndRefreshTokens(AuthCookie.Value, RefreshCookie.Value, requestCSRFToken)
		if err != nil{
		  if err.Error() == "Unauthorised"{
			log.Println("Unauthorised attempt ! JWTs not valid")
			http.Error(w, http.StatusText(401), 401)
			return
		  }else{
			log.Panic("error not nil")
			log.Panic("panic: %+v", err)
			http.Error(w, http.StatusText(500), 500)
			return
		  }
		}
		log.Println("Successfully recreated jwts")
		w.Header().Set("Access-Control-Allow-Origin","")
		setAuthAndRefreshCookies(&w, authTokenString, refreshTokenString)
		w.Header().Set("X-CSRF-Token", csrfSecret)
	default:
	}
	next.ServeHTTP(w,r)
   }
   return http.HandlerFunc(fn)
}

func logicHandler(w http.ResponseWriter, r *http.Request){
	switch r.URL.Path{
	case "/restricted":
		csrfSecret := grabCSRFFromReq(r)
		templates.RenderTemplate(w,"restricted",&templates.RestrictedPage{csrfSecret,"Hello Samridhi"})
	case "/login":
		switch r.Method{
		case "GET":
		case "POST":
		default:
		}
	case "/register":
		switch r.Method{
		case "GET":
			templates.RenderTemplate(w,"register",&templates.RegisterPage{false,""})
		case "POST":
			r.ParseForm()
			log.Println(r.Form)
			_, uuid, err := db.FetchUserByUsername(strings.Join(r.Form["username"],""))
			if err == nil{
				w.WriteHeader(http.StatusUnauthorized)
			}else{
				role := "user"
				uuid, err := db.StoreUser(strings.Join(r.Form["username"],""),strings.Join(r.Form["password"],""),role)
				if err != nil{
					http.Error(w,http.StatusText(500),500)
				}
				log.Println("uuid: " + uuid)
				//register user get cookies
				authTokenString, refreshTokenString, csrfSecret, err := myJwt.CreateNewToken(uuid, role)
				if err != nil{
					http.Error(w, http.StatusText(500), 500)
				}
				setAuthAndRefreshCookies(&w, authTokenString, refreshTokenString)
				w.Header().Set("X-CSRF-Token", csrfSecret)
				w.WriteHeader(http.StatusOK)
			}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)

		}	
	case "/logout":
		nullifyTokenCookies(&w, r)
		http.Redirect(w, r, "/login", 302)
	case "/deleteuser":
    default:
	}
}

func nullifyTokenCookies(w *http.ResponseWriter, r *http.Request){
    authCookie := http.Cookie{
		Name:"AuthToken",
		Value:"",
		Expires:time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(*w, &authCookie)
	refreshCookie := http.Cookie{
		Name:"RefreshToken",
		Value: "",
		Expires: time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(*w, &refreshCookie)
	RefreshCookie, refreshErr := r.Cookie("RefreshToken")
	if refreshErr == http.ErrNoCookie{
        //do nothing
		return
	}else if refreshErr != nil{
		log.Panic("panic: %+v",refreshErr)
		http.Error(*w, http.StatusText(500), 500)
	}
	myJwt.RevokeRefreshToken(RefreshCookie.Value)
}

func setAuthAndRefreshCookies(w *http.ResponseWriter, authTokenString string, refreshTokenString string){
  authCookie := http.Cookie{
	Name:"AuthToken",
	Value:authTokenString,
	HttpOnly: true,
  }
  http.SetCookie(*w, &authCookie)
  refreshCookie := http.Cookie{
	Name:"RefreshToken",
	Value: refreshTokenString,
	HttpOnly: true,
  }
  http.SetCookie(*w, &refreshCookie)
}

func grabCSRFFromReq( r *http.Request) string{
    csrfFrom := r.FormValue("X-CSRF-Token")
	if csrfFrom != ""{
		return csrfFrom
	}else{
		return r.Header.Get("X-CSRF-Token")
	}
}
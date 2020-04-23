package handler

import (
	"encoding/base64"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/rithikjain/SocialRecipe/api/middleware"
	"github.com/rithikjain/SocialRecipe/api/view"
	"github.com/rithikjain/SocialRecipe/pkg/user"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func register(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		_ = r.ParseMultipartForm(10 << 20)
		_ = r.ParseForm()

		var user user.User

		file, handler, err := r.FormFile("image")
		if err != nil {
			user.Name = r.FormValue("name")
			user.PhoneNumber = r.FormValue("phone_number")
			user.Email = r.FormValue("email")
			user.Password = r.FormValue("password")
			user.ProfileImgUrl = "https://res.cloudinary.com/dvn1hxflu/image/upload/v1587624171/blank-profile-picture-973460_640_bgnkjn.png"
			user.ProfileImgPublicID = ""
			user.Bio = r.FormValue("bio")
		} else {
			defer file.Close()

			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				view.Wrap(view.ErrFile, w)
			}
			imgBase64 := base64.StdEncoding.EncodeToString(fileBytes)

			imgUrl := format(imgBase64, handler.Header.Get("Content-Type"))

			// Uploading the image on cloudinary
			form := url.Values{}
			form.Add("file", imgUrl)
			form.Add("upload_preset", os.Getenv("uploadPreset"))

			response, err := http.PostForm(os.Getenv("cloudinaryUrl"), form)
			if err != nil {
				view.Wrap(view.ErrFile, w)
				return
			}
			defer response.Body.Close()

			var resJson map[string]interface{}
			err = json.NewDecoder(response.Body).Decode(&resJson)
			if err != nil {
				view.Wrap(view.ErrUpload, w)
				return
			}

			if response.StatusCode != http.StatusOK {
				view.Wrap(view.ErrUpload, w)
				return
			}

			user.Name = r.FormValue("name")
			user.PhoneNumber = r.FormValue("phone_number")
			user.Email = r.FormValue("email")
			user.Password = r.FormValue("password")
			user.ProfileImgUrl = resJson["secure_url"].(string)
			user.ProfileImgPublicID = resJson["public_id"].(string)
			user.Bio = r.FormValue("bio")
		}

		u, err := svc.Register(&user)
		if err != nil {
			view.Wrap(err, w)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   u.ID,
			"role": "user",
		})
		tokenString, err := token.SignedString([]byte(os.Getenv("jwt_secret")))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		u.Password = ""
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Account Created",
			"token":   tokenString,
			"user":    u,
		})
	})
}

func login(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}
		var user user.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			view.Wrap(err, w)
			return
		}

		u, err := svc.Login(user.Email, user.Password)
		if err != nil {
			view.Wrap(err, w)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   u.ID,
			"role": "user",
		})
		tokenString, err := token.SignedString([]byte(os.Getenv("jwt_secret")))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		u.Password = ""
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Login Successful",
			"token":   tokenString,
			"user":    u,
		})
	})
}

// Protected Request
func userDetails(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		u, err := svc.GetUserByID(uint(claims["id"].(float64)))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User Found",
			"user":    u,
		})
	})
}

// Handlers
func MakeUserHandler(r *http.ServeMux, svc user.Service) {
	r.Handle("/api/v1/user/register", register(svc))
	r.Handle("/api/v1/user/login", login(svc))
	r.Handle("/api/v1/user/details", middleware.Validate(userDetails(svc)))
}

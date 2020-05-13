package handler

import (
	"encoding/base64"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/rithikjain/SocialRecipe/api/middleware"
	"github.com/rithikjain/SocialRecipe/api/view"
	"github.com/rithikjain/SocialRecipe/pkg"
	"github.com/rithikjain/SocialRecipe/pkg/entities"
	"github.com/rithikjain/SocialRecipe/pkg/user"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func register(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		_ = r.ParseMultipartForm(10 << 20)
		_ = r.ParseForm()

		var user entities.User

		email := r.FormValue("email")
		username := r.FormValue("username")

		exist, err := svc.DoesEmailExist(email)
		if err != nil {
			view.Wrap(err, w)
			return
		}
		if exist {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Email exists",
			})
			return
		}

		exist, err = svc.DoesUsernameExist(username)
		if err != nil {
			view.Wrap(err, w)
			return
		}
		if exist {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Username exists",
			})
			return
		}

		file, handler, err := r.FormFile("image")
		if err != nil {
			user.Name = r.FormValue("name")
			user.PhoneNumber = r.FormValue("phone_number")
			user.Email = r.FormValue("email")
			user.Username = r.FormValue("username")
			user.Password = r.FormValue("password")
			user.ProfileImgUrl = "https://res.cloudinary.com/dvn1hxflu/image/upload/v1587624171/blank-profile-picture-973460_640_bgnkjn.png"
			user.ProfileImgPublicID = ""
			user.Verified = false
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
			user.Username = r.FormValue("username")
			user.PhoneNumber = r.FormValue("phone_number")
			user.Email = r.FormValue("email")
			user.Password = r.FormValue("password")
			user.ProfileImgUrl = resJson["secure_url"].(string)
			user.ProfileImgPublicID = resJson["public_id"].(string)
			user.Verified = false
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
		var user entities.User
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
func updateProfile(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		_ = r.ParseMultipartForm(10 << 20)
		_ = r.ParseForm()

		username := r.FormValue("username")
		exist, err := svc.DoesUsernameExist(username)
		if err != nil {
			view.Wrap(err, w)
			return
		}
		if exist {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Username exists",
			})
			return
		}

		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		userID := uint(claims["id"].(float64))

		u, err := svc.GetUserByID(userID)

		if err != nil {
			view.Wrap(err, w)
			return
		}

		var name, un, bio string
		if r.FormValue("name") == "" {
			name = u.Name
		} else {
			name = r.FormValue("name")
		}
		if r.FormValue("username") == "" {
			un = u.Username
		} else {
			un = r.FormValue("username")
		}
		if r.FormValue("bio") == "" {
			bio = u.Bio
		} else {
			bio = r.FormValue("bio")
		}

		file, handler, err := r.FormFile("image")
		if err != nil {
			u.Name = name
			u.Username = un
			u.Bio = bio
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

			u.Name = name
			u.Username = un
			u.ProfileImgUrl = resJson["secure_url"].(string)
			u.ProfileImgPublicID = resJson["public_id"].(string)
			u.Bio = bio
		}

		us, err := svc.UpdateUser(u)
		if err != nil {
			view.Wrap(err, w)
			return
		}
		us.Password = ""
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User updated",
			"user":    us,
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

// Protected Request
func addRecipeToFav(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}

		recipeIDStr := r.URL.Query().Get("recipe_id")
		if recipeIDStr == "" {
			view.Wrap(pkg.ErrNoContent, w)
			return
		}
		recipeID, _ := strconv.Atoi(recipeIDStr)
		err = svc.AddRecipeToFav(uint(claims["id"].(float64)), uint(recipeID))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Recipe added to favorites",
		})
	})
}

// Protected Request
func removeRecipeFromFav(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}

		recipeIDStr := r.URL.Query().Get("recipe_id")
		if recipeIDStr == "" {
			view.Wrap(pkg.ErrNoContent, w)
			return
		}
		recipeID, _ := strconv.Atoi(recipeIDStr)
		err = svc.RemoveRecipeFromFav(uint(claims["id"].(float64)), uint(recipeID))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Recipe removed from favorites",
		})
	})
}

// Protected Request
func followUser(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}

		otherUserIDStr := r.URL.Query().Get("user_id")
		if otherUserIDStr == "" {
			view.Wrap(pkg.ErrNoContent, w)
			return
		}
		otherUserID, _ := strconv.Atoi(otherUserIDStr)

		err = svc.FollowUser(uint(claims["id"].(float64)), uint(otherUserID))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User followed",
		})
	})
}

// Protected Request
func unFollowUser(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}

		otherUserIDStr := r.URL.Query().Get("user_id")
		if otherUserIDStr == "" {
			view.Wrap(pkg.ErrNoContent, w)
			return
		}
		otherUserID, _ := strconv.Atoi(otherUserIDStr)

		err = svc.UnFollowUser(uint(claims["id"].(float64)), uint(otherUserID))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User unfollowed",
		})
	})
}

// Protected Request
func viewFollowers(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr == "" {
			view.Wrap(pkg.ErrNoContent, w)
			return
		}
		userID, _ := strconv.Atoi(userIDStr)

		var pageNo = 1
		pageNoStr := r.URL.Query().Get("page")
		if pageNoStr != "" {
			pageNo, _ = strconv.Atoi(pageNoStr)
		}

		page, err := svc.ViewFollowers(uint(userID), pageNo)
		if err != nil {
			view.Wrap(err, w)
			return
		}

		hasNextPage := true
		if page.Page >= page.TotalPage {
			hasNextPage = false
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Users fetched",
			"users":         page.Records,
			"page":          page.Page,
			"has_next_page": hasNextPage,
			"total_pages":   page.TotalPage,
		})
	})
}

// Protected Request
func viewFollowing(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		userIDStr := r.URL.Query().Get("user_id")
		if userIDStr == "" {
			view.Wrap(pkg.ErrNoContent, w)
			return
		}
		userID, _ := strconv.Atoi(userIDStr)

		var pageNo = 1
		pageNoStr := r.URL.Query().Get("page")
		if pageNoStr != "" {
			pageNo, _ = strconv.Atoi(pageNoStr)
		}

		page, err := svc.ViewFollowing(uint(userID), pageNo)
		if err != nil {
			view.Wrap(err, w)
			return
		}

		hasNextPage := true
		if page.Page >= page.TotalPage {
			hasNextPage = false
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Users fetched",
			"users":         page.Records,
			"page":          page.Page,
			"has_next_page": hasNextPage,
			"total_pages":   page.TotalPage,
		})
	})
}

// Protected Request
func searchUsers(svc user.Service) http.Handler {
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
		userID := uint(claims["id"].(float64))
		query := r.URL.Query().Get("query")
		if query == "" {
			view.Wrap(pkg.ErrNoContent, w)
			return
		}
		query = strings.ToLower(query)
		var pageNo = 1
		pageNoStr := r.URL.Query().Get("page")
		if pageNoStr != "" {
			pageNo, _ = strconv.Atoi(pageNoStr)
		}

		page, err := svc.SearchUsers(userID, query, pageNo)
		if err != nil {
			view.Wrap(err, w)
			return
		}

		hasNextPage := true
		if page.Page >= page.TotalPage {
			hasNextPage = false
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Users fetched",
			"users":         page.Records,
			"page":          page.Page,
			"has_next_page": hasNextPage,
			"total_pages":   page.TotalPage,
		})
	})
}

// Protected Request
func updateBio(svc user.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}
		type Bio struct {
			Bio string `json:"bio"`
		}
		var bioResp Bio
		_ = json.NewDecoder(r.Body).Decode(&bioResp)

		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		userID := uint(claims["id"].(float64))
		err = svc.UpdateUserBio(userID, bioResp.Bio)
		if err != nil {
			view.Wrap(err, w)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Bio updated",
		})
	})
}

// Handlers
func MakeUserHandler(r *http.ServeMux, svc user.Service) {
	r.Handle("/api/v1/user/register", register(svc))
	r.Handle("/api/v1/user/login", login(svc))
	r.Handle("/api/v1/user/updateprofile", middleware.Validate(updateProfile(svc)))
	r.Handle("/api/v1/user/details", middleware.Validate(userDetails(svc)))
	r.Handle("/api/v1/user/addrecipetofav", middleware.Validate(addRecipeToFav(svc)))
	r.Handle("/api/v1/user/removerecipefromfav", middleware.Validate(removeRecipeFromFav(svc)))
	r.Handle("/api/v1/user/follow", middleware.Validate(followUser(svc)))
	r.Handle("/api/v1/user/unfollow", middleware.Validate(unFollowUser(svc)))
	r.Handle("/api/v1/user/viewfollowers", middleware.Validate(viewFollowers(svc)))
	r.Handle("/api/v1/user/viewfollowing", middleware.Validate(viewFollowing(svc)))
	r.Handle("/api/v1/user/search", middleware.Validate(searchUsers(svc)))
	r.Handle("/api/v1/user/updatebio", middleware.Validate(updateBio(svc)))
}

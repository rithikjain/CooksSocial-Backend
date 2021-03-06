package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/rithikjain/SocialRecipe/api/middleware"
	"github.com/rithikjain/SocialRecipe/api/view"
	"github.com/rithikjain/SocialRecipe/pkg"
	"github.com/rithikjain/SocialRecipe/pkg/entities"
	"github.com/rithikjain/SocialRecipe/pkg/recipe"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Protected Request
func createRecipe(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		// Get user id from claims
		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		userID := uint(claims["id"].(float64))

		_ = r.ParseMultipartForm(10 << 20)
		_ = r.ParseForm()

		file, handler, err := r.FormFile("image")
		if err != nil {
			view.Wrap(view.ErrFile, w)
			return
		}
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

		us, err := svc.FindUserByID(userID)
		if err != nil {
			view.Wrap(err, w)
			return
		}

		difficulty, _ := strconv.Atoi(r.FormValue("difficulty"))
		recipe := &entities.Recipe{
			UserID:      userID,
			RecipeName:  r.FormValue("recipe_name"),
			Description: r.FormValue("description"),
			Ingredients: r.FormValue("ingredients"),
			Difficulty:  difficulty,
			Procedure:   r.FormValue("procedure"),
			ImgUrl:      resJson["secure_url"].(string),
			ImgPublicId: resJson["public_id"].(string),
			Name:        us.Name,
			Username:    us.Username,
			UserImg:     us.ProfileImgUrl,
		}
		rec, err := svc.CreateRecipe(recipe)
		if err != nil {
			view.Wrap(err, w)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Recipe created",
			"recipe":  rec,
		})
	})
}

// Protected Request
func updateRecipe(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		// Get user id from claims
		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		userID := uint(claims["id"].(float64))

		_ = r.ParseMultipartForm(10 << 20)
		_ = r.ParseForm()

		file, handler, err := r.FormFile("image")
		if err != nil {
			view.Wrap(view.ErrFile, w)
			return
		}
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

		difficulty, _ := strconv.Atoi(r.FormValue("difficulty"))
		id, _ := strconv.ParseUint(r.FormValue("recipe_id"), 10, 32)
		rec, err := svc.FindRecipeByID(uint(id))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		rec.ID = uint(id)
		rec.UserID = userID
		rec.RecipeName = r.FormValue("recipe_name")
		rec.Description = r.FormValue("description")
		rec.Ingredients = r.FormValue("ingredients")
		rec.Difficulty = difficulty
		rec.Procedure = r.FormValue("procedure")
		rec.ImgUrl = resJson["secure_url"].(string)
		rec.ImgPublicId = resJson["public_id"].(string)

		re, err := svc.UpdateRecipe(rec)
		if err != nil {
			view.Wrap(err, w)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Recipe updated",
			"recipe":  re,
		})
	})
}

// Protected Request
func deleteRecipe(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		recipeIDStr := r.URL.Query().Get("recipe_id")
		if recipeIDStr == "" {
			view.Wrap(pkg.ErrNoContent, w)
			return
		}
		recipeID, _ := strconv.Atoi(recipeIDStr)
		rec, err := svc.FindRecipeByID(uint(recipeID))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		// Get user id from claims
		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		userID := uint(claims["id"].(float64))

		if rec.UserID != userID {
			view.Wrap(pkg.ErrUnauthorized, w)
			return
		}

		req, err := http.NewRequest("DELETE", os.Getenv("cloudinaryDeleteUrl"), nil)
		if err != nil {
			view.Wrap(err, w)
			return
		}
		q := req.URL.Query()
		q.Add("public_ids", rec.ImgPublicId)
		req.URL.RawQuery = q.Encode()

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			view.Wrap(err, w)
			return
		}
		if res.StatusCode != http.StatusOK {
			view.Wrap(view.ErrFile, w)
			return
		}
		err = svc.DeleteRecipe(uint(recipeID))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Recipe deleted",
		})
	})
}

// Protected Request
func showAllRecipesOfUser(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		userIDStr := r.URL.Query().Get("user_id")
		userID, _ := strconv.Atoi(userIDStr)

		var pageNo = 1
		pageNoStr := r.URL.Query().Get("page")
		if pageNoStr != "" {
			pageNo, _ = strconv.Atoi(pageNoStr)
		}

		page, err := svc.GetAllRecipesOfUser(uint(userID), pageNo)
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
			"message":       "Recipes fetched",
			"recipes":       page.Records,
			"page":          page.Page,
			"has_next_page": hasNextPage,
			"total_pages":   page.TotalPage,
		})
	})
}

// Protected Request
func showMyRecipes(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		// Get user id from claims
		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		userID := uint(claims["id"].(float64))

		var pageNo = 1
		pageNoStr := r.URL.Query().Get("page")
		if pageNoStr != "" {
			pageNo, _ = strconv.Atoi(pageNoStr)
		}

		page, err := svc.GetAllRecipesOfUser(userID, pageNo)
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
			"message":       "Recipes fetched",
			"recipes":       page.Records,
			"page":          page.Page,
			"has_next_page": hasNextPage,
			"total_pages":   page.TotalPage,
		})
	})
}

func showMyFavRecipes(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		// Get user id from claims
		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		userID := uint(claims["id"].(float64))

		var pageNo = 1
		pageNoStr := r.URL.Query().Get("page")
		if pageNoStr != "" {
			pageNo, _ = strconv.Atoi(pageNoStr)
		}

		page, err := svc.ShowUsersFavRecipes(userID, pageNo)
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
			"message":       "Favorite recipes fetched",
			"recipes":       page.Records,
			"page":          page.Page,
			"has_next_page": hasNextPage,
			"total_pages":   page.TotalPage,
		})
	})
}

// Protected Request
func likeRecipe(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		// Get user id from claims
		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		userID := uint(claims["id"].(float64))

		recipeIDStr := r.URL.Query().Get("recipe_id")
		recipeID, _ := strconv.Atoi(recipeIDStr)

		hasLiked, err := svc.HasUserLiked(userID, uint(recipeID))
		if err != nil {
			view.Wrap(err, w)
			return
		}
		if hasLiked {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Recipe was already liked",
			})
		} else {
			err = svc.LikeRecipe(userID, uint(recipeID))
			if err != nil {
				view.Wrap(err, w)
				return
			}

			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Recipe liked",
			})
		}
	})
}

// Protected Request
func unlikeRecipe(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		// Get user id from claims
		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		userID := uint(claims["id"].(float64))

		recipeIDStr := r.URL.Query().Get("recipe_id")
		recipeID, _ := strconv.Atoi(recipeIDStr)

		err = svc.UnlikeRecipe(userID, uint(recipeID))
		if err != nil {
			view.Wrap(err, w)
			return
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Recipe unliked",
		})
	})
}

// Protected Request
func showUsersWhoLiked(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		recipeIDStr := r.URL.Query().Get("recipe_id")
		recipeID, _ := strconv.Atoi(recipeIDStr)

		var pageNo = 1
		pageNoStr := r.URL.Query().Get("page")
		if pageNoStr != "" {
			pageNo, _ = strconv.Atoi(pageNoStr)
		}

		page, err := svc.ShowUsersWhoLiked(uint(recipeID), pageNo)
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
func showUserFeed(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		// Get user id from claims
		claims, err := middleware.ValidateAndGetClaims(r.Context(), "user")
		if err != nil {
			view.Wrap(err, w)
			return
		}
		userID := uint(claims["id"].(float64))

		var pageNo = 1
		pageNoStr := r.URL.Query().Get("page")
		if pageNoStr != "" {
			pageNo, _ = strconv.Atoi(pageNoStr)
		}

		page, err := svc.ShowUserFeed(userID, pageNo)
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
			"message":       "Recipes fetched",
			"recipes":       page.Records,
			"page":          page.Page,
			"has_next_page": hasNextPage,
			"total_pages":   page.TotalPage,
		})
	})
}

// Protected Request
func showAllLatestRecipes(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		var pageNo = 1
		pageNoStr := r.URL.Query().Get("page")
		if pageNoStr != "" {
			pageNo, _ = strconv.Atoi(pageNoStr)
		}

		page, err := svc.ShowAllLatestRecipes(pageNo)
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
			"message":       "Recipes fetched",
			"recipes":       page.Records,
			"page":          page.Page,
			"has_next_page": hasNextPage,
			"total_pages":   page.TotalPage,
		})
	})
}

// Protected Request
func searchRecipes(svc recipe.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			view.Wrap(view.ErrMethodNotAllowed, w)
			return
		}

		var pageNo = 1
		pageNoStr := r.URL.Query().Get("page")
		if pageNoStr != "" {
			pageNo, _ = strconv.Atoi(pageNoStr)
		}

		query := r.URL.Query().Get("query")
		if query == "" {
			view.Wrap(pkg.ErrNoContent, w)
			return
		}
		query = strings.ToLower(query)

		page, err := svc.SearchRecipes(query, pageNo)
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
			"message":       "Recipes fetched",
			"recipes":       page.Records,
			"page":          page.Page,
			"has_next_page": hasNextPage,
			"total_pages":   page.TotalPage,
		})
	})
}

func format(encStr string, mime string) string {
	switch mime {
	case "image/gif", "image/jpeg", "image/pjpeg", "image/png", "image/tiff":
		return fmt.Sprintf("data:%s;base64,%s", mime, encStr)
	default:
	}

	return fmt.Sprintf("data:image/png;base64,%s", encStr)
}

func MakeRecipeHandler(r *http.ServeMux, svc recipe.Service) {
	r.Handle("/api/v1/recipe/create", middleware.Validate(createRecipe(svc)))
	r.Handle("/api/v1/recipe/update", middleware.Validate(updateRecipe(svc)))
	r.Handle("/api/v1/recipe/delete", middleware.Validate(deleteRecipe(svc)))
	r.Handle("/api/v1/recipe/viewofuser", middleware.Validate(showAllRecipesOfUser(svc)))
	r.Handle("/api/v1/recipe/viewmine", middleware.Validate(showMyRecipes(svc)))
	r.Handle("/api/v1/recipe/viewmyfeed", middleware.Validate(showUserFeed(svc)))
	r.Handle("/api/v1/recipe/viewmyfav", middleware.Validate(showMyFavRecipes(svc)))
	r.Handle("/api/v1/recipe/explore", middleware.Validate(showAllLatestRecipes(svc)))
	r.Handle("/api/v1/recipe/like", middleware.Validate(likeRecipe(svc)))
	r.Handle("/api/v1/recipe/unlike", middleware.Validate(unlikeRecipe(svc)))
	r.Handle("/api/v1/recipe/viewuserlikes", middleware.Validate(showUsersWhoLiked(svc)))
	r.Handle("/api/v1/recipe/search", middleware.Validate(searchRecipes(svc)))
}

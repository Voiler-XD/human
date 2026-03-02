package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/barun-bash/human/human-studio/server/auth"
	"github.com/barun-bash/human/human-studio/server/middleware"
	"github.com/barun-bash/human/human-studio/server/models"
)

func Signup(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			jsonError(w, "email and password are required", http.StatusBadRequest)
			return
		}

		if len(req.Password) < 8 {
			jsonError(w, "password must be at least 8 characters", http.StatusBadRequest)
			return
		}

		resp, err := svc.Signup(req)
		if err != nil {
			if errors.Is(err, auth.ErrEmailTaken) {
				jsonError(w, "email already registered", http.StatusConflict)
				return
			}
			jsonError(w, "signup failed", http.StatusInternalServerError)
			return
		}

		jsonResponse(w, resp, http.StatusCreated)
	}
}

func Login(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := svc.Login(req)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				jsonError(w, "invalid email or password", http.StatusUnauthorized)
				return
			}
			jsonError(w, "login failed", http.StatusInternalServerError)
			return
		}

		jsonResponse(w, resp, http.StatusOK)
	}
}

func RefreshToken(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := svc.RefreshTokens(body.RefreshToken)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidToken) {
				jsonError(w, "invalid or expired refresh token", http.StatusUnauthorized)
				return
			}
			jsonError(w, "token refresh failed", http.StatusInternalServerError)
			return
		}

		jsonResponse(w, resp, http.StatusOK)
	}
}

func ResetPassword(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// TODO: Send password reset email
		// Always return success to prevent email enumeration
		jsonResponse(w, map[string]string{"message": "if an account exists, a reset email has been sent"}, http.StatusOK)
	}
}

func GetProfile(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		user, err := svc.GetUser(userID)
		if err != nil {
			jsonError(w, "user not found", http.StatusNotFound)
			return
		}
		jsonResponse(w, user, http.StatusOK)
	}
}

func UpdateProfile(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		var update models.ProfileUpdate
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, err := svc.UpdateProfile(userID, update)
		if err != nil {
			jsonError(w, "update failed", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, user, http.StatusOK)
	}
}

func ChangePassword(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		var change models.PasswordChange
		if err := json.NewDecoder(r.Body).Decode(&change); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if len(change.NewPassword) < 8 {
			jsonError(w, "new password must be at least 8 characters", http.StatusBadRequest)
			return
		}

		err := svc.ChangePassword(userID, change)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				jsonError(w, "current password is incorrect", http.StatusUnauthorized)
				return
			}
			jsonError(w, "password change failed", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"message": "password updated"}, http.StatusOK)
	}
}

func DeleteAccount(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		if err := svc.DeleteUser(userID); err != nil {
			jsonError(w, "account deletion failed", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"message": "account deleted"}, http.StatusOK)
	}
}

func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

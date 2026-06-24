package handlers

import (
	"log/slog"
	"net/http"
	"time"
	"fitmind/middleware"
	"fitmind/models"
	"fitmind/services"
)

type AuthHandler struct {
	supabase *services.SupabaseService
}

func NewAuthHandler(sb *services.SupabaseService) *AuthHandler {
	return &AuthHandler{supabase: sb}
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		SendJSON(w, http.StatusUnauthorized, false, nil, "Unauthorized")
		return
	}

	profile, err := h.supabase.GetUserProfile(userID)
	if err != nil {
		slog.Error("Failed to get user profile", "user_id", userID, "error", err)
		SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to fetch user profile")
		return
	}

	todayStr := time.Now().Format("2006-01-02")
	yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	if profile == nil {
		// Profile does not exist (trigger might have failed or not completed yet)
		// Create a new default profile
		profile = &models.UserProfile{
			ID:          userID,
			FullName:    "Health Enthusiast",
			StreakCount: 1,
			LastActive:  todayStr,
			CreatedAt:   time.Now(),
		}
		err = h.supabase.CreateUserProfile(profile)
		if err != nil {
			slog.Error("Failed to create default user profile", "user_id", userID, "error", err)
			SendJSON(w, http.StatusInternalServerError, false, nil, "Failed to create user profile")
			return
		}
	} else {
		// Profile exists, update streak count
		updates := make(map[string]interface{})
		needUpdate := false

		if profile.LastActive == yesterdayStr {
			profile.StreakCount += 1
			profile.LastActive = todayStr
			updates["streak_count"] = profile.StreakCount
			updates["last_active"] = todayStr
			needUpdate = true
		} else if profile.LastActive != todayStr {
			// Missed days or first time login since a while
			profile.StreakCount = 1
			profile.LastActive = todayStr
			updates["streak_count"] = 1
			updates["last_active"] = todayStr
			needUpdate = true
		}

		if needUpdate {
			err = h.supabase.UpdateUserProfile(userID, updates)
			if err != nil {
				slog.Error("Failed to update user profile streak", "user_id", userID, "error", err)
				// Do not block the request if streak update fails, just proceed
			}
		}
	}

	SendJSON(w, http.StatusOK, true, profile, "")
}

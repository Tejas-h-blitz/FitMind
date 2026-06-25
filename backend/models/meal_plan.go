package models

import "time"

type MealItem struct {
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
	Calories    int      `json:"calories"`
	Benefits    string   `json:"benefits"`
}

type DayMeals struct {
	Breakfast MealItem `json:"breakfast"`
	Lunch     MealItem `json:"lunch"`
	Dinner    MealItem `json:"dinner"`
	Snacks    MealItem `json:"snacks"`
}

type MealDay struct {
	Day      string   `json:"day"`
	Meals    DayMeals `json:"meals"`
	DailyTip string   `json:"daily_tip"`
}

type MealPlan struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	DocID               string    `json:"doc_id"`
	Reasoning           string    `json:"reasoning"`
	DailyCaloriesTarget int       `json:"daily_calories_target"`
	ProteinTargetG      int       `json:"protein_target_g"`
	Days                []MealDay `json:"days"`
	WeeklyNotes         string    `json:"weekly_notes"`
	DietaryPreference   string    `json:"dietary_preference"`
	CreatedAt           time.Time `json:"created_at"`
}

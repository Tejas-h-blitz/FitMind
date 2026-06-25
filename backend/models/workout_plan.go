package models

import "time"

type Exercise struct {
	Name               string   `json:"name"`
	Sets               int      `json:"sets"`
	Reps               string   `json:"reps"`
	RestSeconds        int      `json:"rest_seconds"`
	Instructions       string   `json:"instructions"`
	MuscleGroups       []string `json:"muscle_groups"`
	ModificationEasier string   `json:"modification_easier"`
	ModificationHarder string   `json:"modification_harder"`
}

type WorkoutDay struct {
	Day                      string     `json:"day"`
	Focus                    string     `json:"focus"`
	IsRestDay                bool       `json:"is_rest_day"`
	Exercises                []Exercise `json:"exercises"`
	EstimatedDurationMinutes int        `json:"estimated_duration_minutes"`
	Warmup                   string     `json:"warmup"`
	Cooldown                 string     `json:"cooldown"`
}

type WorkoutPlan struct {
	ID               string       `json:"id"`
	UserID           string       `json:"user_id"`
	ProgramName      string       `json:"program_name"`
	DurationWeeks    int          `json:"duration_weeks"`
	Reasoning        string       `json:"reasoning"`
	WeeklySchedule   []WorkoutDay `json:"weekly_schedule"`
	ProgressionNotes string       `json:"progression_notes"`
	SafetyNotes      string       `json:"safety_notes"`
	FitnessLevel     string       `json:"fitness_level"`
	Equipment        string       `json:"equipment"`
	DaysPerWeek      int          `json:"days_per_week"`
	CreatedAt        time.Time    `json:"created_at"`
}

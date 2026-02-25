package model

import (
	"time"

	"github.com/uptrace/bun"
)

// CtxKey is the type for context keys used in middleware.
type CtxKey string

const (
	CtxKeyVisitorIP CtxKey = "visitor_ip"
	CtxKeyVisitorUA CtxKey = "visitor_ua"
)

// Visitor represents a unique site visitor, tracked by IP address.
type Visitor struct {
	bun.BaseModel `bun:"table:visitors,alias:vis"`

	ID           int64     `json:"id" bun:"id,pk,autoincrement"`
	IP           string    `json:"ip" bun:"ip,type:text,unique,notnull"`
	UserAgent    string    `json:"user_agent" bun:"user_agent,type:text"`
	Page         string    `json:"page" bun:"page,type:text"`
	Referrer     string    `json:"referrer" bun:"referrer,type:text"`
	ScreenWidth  int       `json:"screen_width" bun:"screen_width,default:0"`
	ScreenHeight int       `json:"screen_height" bun:"screen_height,default:0"`
	Language     string    `json:"language" bun:"language,type:text"`
	VisitCount   int       `json:"visit_count" bun:"visit_count,notnull,default:1"`
	FirstVisitAt time.Time `json:"first_visit_at" bun:"first_visit_at,nullzero,notnull,default:current_timestamp"`
	LastVisitAt  time.Time `json:"last_visit_at" bun:"last_visit_at,nullzero,notnull,default:current_timestamp"`
}

// DailyVisitCount holds visit count for a single day.
type DailyVisitCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// VisitStats contains aggregated visit statistics.
type VisitStats struct {
	// Active visitors by period (by last_visit_at)
	TotalVisits int `json:"total_visits"`
	TodayVisits int `json:"today_visits"`
	WeekVisits  int `json:"week_visits"`
	MonthVisits int `json:"month_visits"`

	// New visitors by period (by first_visit_at)
	NewToday int `json:"new_today"`
	NewWeek  int `json:"new_week"`
	NewMonth int `json:"new_month"`

	// Engagement
	ReturningVisitors int     `json:"returning_visitors"`
	TotalPageViews    int     `json:"total_page_views"`
	AvgVisitsPerUser  float64 `json:"avg_visits_per_user"`

	// Entity counts
	ExhibitionCount int `json:"exhibition_count"`
	ExhibitCount    int `json:"exhibit_count"`
	NewsCount       int `json:"news_count"`

	// Daily breakdown (last 7 days)
	DailyVisits []DailyVisitCount `json:"daily_visits"`
}

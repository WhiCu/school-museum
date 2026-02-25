package storage

import (
	"context"
	"time"

	"github.com/WhiCu/school-museum/db/model"
	"github.com/uptrace/bun"
)

// VisitStorage handles visitor tracking and statistics.
type VisitStorage struct {
	db *bun.DB
}

func NewVisitStorage(db *bun.DB) *VisitStorage {
	return &VisitStorage{db: db}
}

// Record upserts a visitor record by IP.
// On first visit: creates a new record with visit_count=1.
// On return visit: increments visit_count, updates last_visit_at and other fields.
func (s *VisitStorage) Record(ctx context.Context, v model.Visitor) error {
	v.LastVisitAt = time.Now()
	if v.FirstVisitAt.IsZero() {
		v.FirstVisitAt = v.LastVisitAt
	}
	v.VisitCount = 1

	_, err := s.db.NewInsert().
		Model(&v).
		On("CONFLICT (ip) DO UPDATE").
		Set("visit_count = vis.visit_count + 1").
		Set("last_visit_at = EXCLUDED.last_visit_at").
		Set("user_agent = EXCLUDED.user_agent").
		Set("page = EXCLUDED.page").
		Set("referrer = CASE WHEN EXCLUDED.referrer != '' THEN EXCLUDED.referrer ELSE vis.referrer END").
		Set("screen_width = CASE WHEN EXCLUDED.screen_width > 0 THEN EXCLUDED.screen_width ELSE vis.screen_width END").
		Set("screen_height = CASE WHEN EXCLUDED.screen_height > 0 THEN EXCLUDED.screen_height ELSE vis.screen_height END").
		Set("language = CASE WHEN EXCLUDED.language != '' THEN EXCLUDED.language ELSE vis.language END").
		Exec(ctx)
	return err
}

// Stats returns aggregated visit statistics along with entity counts.
func (s *VisitStorage) Stats(ctx context.Context) (model.VisitStats, error) {
	var stats model.VisitStats

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := today.AddDate(0, 0, -7)
	monthAgo := today.AddDate(0, -1, 0)

	// ── Active visitors by last_visit_at ──

	total, err := s.db.NewSelect().Model((*model.Visitor)(nil)).Count(ctx)
	if err != nil {
		return stats, err
	}
	stats.TotalVisits = total

	todayCount, err := s.db.NewSelect().Model((*model.Visitor)(nil)).
		Where("last_visit_at >= ?", today).Count(ctx)
	if err != nil {
		return stats, err
	}
	stats.TodayVisits = todayCount

	weekCount, err := s.db.NewSelect().Model((*model.Visitor)(nil)).
		Where("last_visit_at >= ?", weekAgo).Count(ctx)
	if err != nil {
		return stats, err
	}
	stats.WeekVisits = weekCount

	monthCount, err := s.db.NewSelect().Model((*model.Visitor)(nil)).
		Where("last_visit_at >= ?", monthAgo).Count(ctx)
	if err != nil {
		return stats, err
	}
	stats.MonthVisits = monthCount

	// ── New visitors by first_visit_at ──

	newToday, err := s.db.NewSelect().Model((*model.Visitor)(nil)).
		Where("first_visit_at >= ?", today).Count(ctx)
	if err != nil {
		return stats, err
	}
	stats.NewToday = newToday

	newWeek, err := s.db.NewSelect().Model((*model.Visitor)(nil)).
		Where("first_visit_at >= ?", weekAgo).Count(ctx)
	if err != nil {
		return stats, err
	}
	stats.NewWeek = newWeek

	newMonth, err := s.db.NewSelect().Model((*model.Visitor)(nil)).
		Where("first_visit_at >= ?", monthAgo).Count(ctx)
	if err != nil {
		return stats, err
	}
	stats.NewMonth = newMonth

	// ── Engagement ──

	returningCount, err := s.db.NewSelect().Model((*model.Visitor)(nil)).
		Where("visit_count > 1").Count(ctx)
	if err != nil {
		return stats, err
	}
	stats.ReturningVisitors = returningCount

	var totalPageViews int
	err = s.db.NewSelect().Model((*model.Visitor)(nil)).
		ColumnExpr("COALESCE(SUM(visit_count), 0)").
		Scan(ctx, &totalPageViews)
	if err != nil {
		return stats, err
	}
	stats.TotalPageViews = totalPageViews

	if total > 0 {
		stats.AvgVisitsPerUser = float64(totalPageViews) / float64(total)
	}

	// ── Daily breakdown (last 7 days by last_visit_at) ──

	type dayRow struct {
		Day   time.Time `bun:"day"`
		Count int       `bun:"count"`
	}
	var rows []dayRow
	err = s.db.NewSelect().
		TableExpr("generate_series(?::date, ?::date, '1 day'::interval) AS day", weekAgo, today).
		ColumnExpr("day::date AS day").
		ColumnExpr("(SELECT COUNT(*) FROM visitors WHERE last_visit_at >= day AND last_visit_at < day + '1 day'::interval) AS count").
		Scan(ctx, &rows)
	if err != nil {
		return stats, err
	}
	stats.DailyVisits = make([]model.DailyVisitCount, 0, len(rows))
	for _, r := range rows {
		stats.DailyVisits = append(stats.DailyVisits, model.DailyVisitCount{
			Date:  r.Day.Format("2006-01-02"),
			Count: r.Count,
		})
	}

	// ── Entity counts ──

	exCount, _ := s.db.NewSelect().Model((*model.Exhibition)(nil)).
		Where("deleted_at IS NULL").Count(ctx)
	stats.ExhibitionCount = exCount

	eCount, _ := s.db.NewSelect().Model((*model.Exhibit)(nil)).
		Where("deleted_at IS NULL").Count(ctx)
	stats.ExhibitCount = eCount

	nCount, _ := s.db.NewSelect().Model((*model.News)(nil)).
		Where("deleted_at IS NULL").Count(ctx)
	stats.NewsCount = nCount

	return stats, nil
}

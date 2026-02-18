package utils

import (
	"fmt"
	"time"
)

func Now() time.Time {
	return time.Now().UTC()
}

func Today() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func Tomorrow() time.Time {
	return Today().AddDate(0, 0, 1)
}

func Yesterday() time.Time {
	return Today().AddDate(0, 0, -1)
}

func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return t.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
}

func EndOfWeek(t time.Time) time.Time {
	return StartOfWeek(t).AddDate(0, 0, 7).Add(-time.Second)
}

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Second)
}

func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	return time.Date(t.Year()+1, 1, 1, 0, 0, 0, -1, t.Location())
}

func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

func DaysBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours() / 24)
}

func HoursBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours())
}

func MinutesBetween(start, end time.Time) int {
	return int(end.Sub(start).Minutes())
}

func SecondsBetween(start, end time.Time) int {
	return int(end.Sub(start).Seconds())
}

func IsWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

func ParseDate(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func ParseDateISO(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func FormatDate(t time.Time, layout string) string {
	return t.Format(layout)
}

func FormatDateISO(t time.Time) string {
	return t.Format(time.RFC3339)
}

func FormatDateShort(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDateTimeShort(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func FormatDateTimeFull(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func FormatRelative(t time.Time) string {
	now := Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	}
	if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / (24 * 7))
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	}
	if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	years := int(diff.Hours() / (24 * 365))
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

func IsExpired(expiry time.Time) bool {
	return Now().After(expiry)
}

func IsNotExpired(expiry time.Time) bool {
	return !IsExpired(expiry)
}

package handlers

import (
	"strings"
	"testing"
	"time"

	"fmps/models"
)

// ─── ageAtDate ────────────────────────────────────────────────────────────────

func TestAgeAtDate(t *testing.T) {
	tests := []struct {
		name   string
		birth  string
		target string
		want   int
	}{
		// Normal cases
		{"exact birthday", "2000-03-15", "2020-03-15", 20},
		{"day before birthday", "2000-03-15", "2020-03-14", 19},
		{"day after birthday", "2000-03-15", "2020-03-16", 20},
		// Year boundary
		{"new year birth on new year", "2000-01-01", "2020-01-01", 20},
		{"dec31 birth on dec31", "2000-12-31", "2020-12-31", 20},
		{"dec31 birth day before", "2000-12-31", "2020-12-30", 19},
		// Feb 29 birth – birthday treated as Mar 1 in ALL years
		{"feb29 birth, target mar1 non-leap", "2000-02-29", "2019-03-01", 19},
		{"feb29 birth, target feb28 non-leap", "2000-02-29", "2019-02-28", 18},
		{"feb29 birth, target mar2 non-leap", "2000-02-29", "2019-03-02", 19},
		{"feb29 birth, target mar1 leap", "2000-02-29", "2020-03-01", 20},
		{"feb29 birth, target feb29 leap", "2000-02-29", "2020-02-29", 19},
		// Spec case from requirements: birth 2004-02-29, issue 2020-02-29 → age 15
		{"spec: feb29-2004 at feb29-2020 = 15", "2004-02-29", "2020-02-29", 15},
		// Verify age 16 boundary for feb29 birth
		{"feb29-2004 at mar1-2020 = 16", "2004-02-29", "2020-03-01", 16},
		// Age 0
		{"just born", "2020-06-15", "2020-06-15", 0},
		{"before first birthday", "2020-06-15", "2021-06-14", 0},
		{"first birthday", "2020-06-15", "2021-06-15", 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := time.Parse("2006-01-02", tc.birth)
			if err != nil {
				t.Fatalf("parse birth: %v", err)
			}
			target, err := time.Parse("2006-01-02", tc.target)
			if err != nil {
				t.Fatalf("parse target: %v", err)
			}
			got := ageAtDate(b.UTC(), target.UTC())
			if got != tc.want {
				t.Errorf("ageAtDate(%s, %s) = %d, want %d", tc.birth, tc.target, got, tc.want)
			}
		})
	}
}

// ─── addSameDayYears ──────────────────────────────────────────────────────────

func TestAddSameDayYears(t *testing.T) {
	tests := []struct {
		name  string
		issue string
		years int
		want  string
	}{
		{"normal +5", "2020-03-15", 5, "2025-03-15"},
		{"normal +10", "2015-07-01", 10, "2025-07-01"},
		{"normal +20", "2000-01-01", 20, "2020-01-01"},
		// Feb 29 → target year is leap
		{"feb29 +4 leap", "2000-02-29", 4, "2004-02-29"},
		{"feb29 +8 leap", "2000-02-29", 8, "2008-02-29"},
		// Feb 29 → target year not leap → Mar 1
		{"feb29 +1 non-leap", "2000-02-29", 1, "2001-03-01"},
		{"feb29 +5 non-leap", "2000-02-29", 5, "2005-03-01"},
		{"feb29 +10 non-leap", "2000-02-29", 10, "2010-03-01"},
		// Spec from requirements: issue 2020-02-29 +5 = 2025-03-01 (2025 not leap)
		{"spec: feb29-2020 +5 = 2025-03-01", "2020-02-29", 5, "2025-03-01"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			issue, err := time.Parse("2006-01-02", tc.issue)
			if err != nil {
				t.Fatalf("parse issue: %v", err)
			}
			got := addSameDayYears(issue.UTC(), tc.years)
			if got.Format("2006-01-02") != tc.want {
				t.Errorf("addSameDayYears(%s, %d) = %s, want %s",
					tc.issue, tc.years, got.Format("2006-01-02"), tc.want)
			}
		})
	}
}

// ─── validateDocumentDates ────────────────────────────────────────────────────

func makeMemberForDate(docType, birth, issue, expiry string) *models.Member {
	return &models.Member{
		IdDocType:    docType,
		BirthDate:    birth,
		IdIssueDate:  issue,
		IdExpiryDate: expiry,
	}
}

// futureDate returns a date string N days from today.
func futureDate(days int) string {
	return time.Now().UTC().AddDate(0, 0, days).Format("2006-01-02")
}

func TestValidateDocumentDates(t *testing.T) {
	// Spec required case: birth 2004-02-29, issue 2020-02-29 → type 01 → age=15 → 5yr
	// The expected expiry is 2025-03-01 (which is now in the past).
	// Submitting a different future expiry proves the 5yr rule was chosen (error says 2025-03-01).
	t.Run("spec: feb29-2004 birth, feb29-2020 issue, age=15, wrong future expiry confirms 5yr rule", func(t *testing.T) {
		m := makeMemberForDate("01", "2004-02-29", "2020-02-29", "2035-02-28")
		err := validateDocumentDates(m)
		if err == nil {
			t.Fatal("expected error (wrong expiry), got nil")
		}
		if !strings.Contains(err.Error(), "2025-03-01") {
			t.Errorf("error should mention expected expiry 2025-03-01, got: %v", err)
		}
	})

	// ── Type 01 tests using future-safe dates ─────────────────────────────────

	t.Run("01: age<16 at issue → 5yr", func(t *testing.T) {
		// born 2012-06-01, issued 2025-06-01 (age=13), expiry = addSameDayYears(2025-06-01,5) = 2030-06-01
		m := makeMemberForDate("01", "2012-06-01", "2025-06-01", "2030-06-01")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("01: age<16 → 5yr → wrong expiry", func(t *testing.T) {
		m := makeMemberForDate("01", "2012-06-01", "2025-06-01", "2035-06-01")
		if err := validateDocumentDates(m); err == nil {
			t.Error("expected error but got nil")
		}
	})

	t.Run("01: age 20 (16-25) → 10yr", func(t *testing.T) {
		// born 2006-01-01, issued 2026-01-01 (age=20), expiry = 2036-01-01
		m := makeMemberForDate("01", "2006-01-01", "2026-01-01", "2036-01-01")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("01: age 40 (26-45) → 20yr", func(t *testing.T) {
		// born 1986-01-01, issued 2026-01-01 (age=40), expiry = 2046-01-01
		m := makeMemberForDate("01", "1986-01-01", "2026-01-01", "2046-01-01")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("01: age>=46 → long-term 2099-12-31", func(t *testing.T) {
		// born 1970-01-01, issued 2026-01-01 (age=56), expiry = 2099-12-31
		m := makeMemberForDate("01", "1970-01-01", "2026-01-01", "2099-12-31")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("01: age>=46 wrong expiry → error mentioning long-term", func(t *testing.T) {
		m := makeMemberForDate("01", "1970-01-01", "2026-01-01", "2046-01-01")
		err := validateDocumentDates(m)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "长期") {
			t.Errorf("error should mention 长期, got: %v", err)
		}
	})

	t.Run("01: feb29 birth, feb29 issue (age<16) → 5yr non-leap expiry", func(t *testing.T) {
		// born 2008-02-29, issued 2024-02-29 (age = 15 because feb29 birthday treated as Mar 1)
		// expiry = addSameDayYears(2024-02-29, 5) = 2029-03-01 (2029 not leap)
		m := makeMemberForDate("01", "2008-02-29", "2024-02-29", "2029-03-01")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// ── Type 91 ───────────────────────────────────────────────────────────────

	t.Run("91: under 16, expiry = 16th birthday", func(t *testing.T) {
		// born 2015-05-01, 16th birthday 2031-05-01; issued 2024-05-01
		m := makeMemberForDate("91", "2015-05-01", "2024-05-01", "2031-05-01")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("91: under 16, wrong expiry by 1 day", func(t *testing.T) {
		m := makeMemberForDate("91", "2015-05-01", "2024-05-01", "2031-05-02")
		if err := validateDocumentDates(m); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("91: holder already 16+ → error", func(t *testing.T) {
		// born 2008-01-01 → age today is >=16 (2026); expiry 2031-01-01
		m := makeMemberForDate("91", "2008-01-01", "2024-01-01", "2031-01-01")
		err := validateDocumentDates(m)
		if err == nil {
			t.Fatal("expected error for age>=16 type-91")
		}
		if !strings.Contains(err.Error(), "已满16周岁") {
			t.Errorf("error should mention 已满16周岁, got: %v", err)
		}
	})

	// ── Type 11 ───────────────────────────────────────────────────────────────

	t.Run("11: +5yr correct", func(t *testing.T) {
		// issue 2024-06-01, expiry 2029-06-01
		m := makeMemberForDate("11", "1990-01-01", "2024-06-01", "2029-06-01")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("11: +5yr wrong by 1 day", func(t *testing.T) {
		m := makeMemberForDate("11", "1990-01-01", "2024-06-01", "2029-06-02")
		if err := validateDocumentDates(m); err == nil {
			t.Error("expected error, got nil")
		}
	})

	// ── Type 21 ───────────────────────────────────────────────────────────────

	t.Run("21: +5yr correct", func(t *testing.T) {
		m := makeMemberForDate("21", "1990-01-01", "2024-06-15", "2029-06-15")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// ── Type 31 ───────────────────────────────────────────────────────────────

	t.Run("31: age<18 at issue → 5yr-1day", func(t *testing.T) {
		// born 2008-01-01, issued 2024-01-01 (age=16) → addYearsMinusOneDay(2024-01-01,5)=2028-12-31
		m := makeMemberForDate("31", "2008-01-01", "2024-01-01", "2028-12-31")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("31: age>=18 at issue → 10yr-1day", func(t *testing.T) {
		// born 2004-01-01, issued 2024-01-01 (age=20) → addYearsMinusOneDay(2024-01-01,10)=2033-12-31
		m := makeMemberForDate("31", "2004-01-01", "2024-01-01", "2033-12-31")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// ── Type 02 ───────────────────────────────────────────────────────────────

	t.Run("02: age<18 at issue → 5yr-1day", func(t *testing.T) {
		// born 2008-06-01, issued 2024-06-01 (age=16) → addYearsMinusOneDay(2024-06-01,5)=2029-05-31
		m := makeMemberForDate("02", "2008-06-01", "2024-06-01", "2029-05-31")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("02: age>=18 at issue → 10yr-1day", func(t *testing.T) {
		// born 2000-01-01, issued 2024-01-01 (age=24) → addYearsMinusOneDay(2024-01-01,10)=2033-12-31
		m := makeMemberForDate("02", "2000-01-01", "2024-01-01", "2033-12-31")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// ── Type 03 ───────────────────────────────────────────────────────────────

	t.Run("03: 5yr-1day correct", func(t *testing.T) {
		// issue 2024-03-15 → addYearsMinusOneDay(2024-03-15,5)=2029-03-14
		m := makeMemberForDate("03", "1990-01-01", "2024-03-15", "2029-03-14")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("03: 5yr-1day wrong by 1 day", func(t *testing.T) {
		m := makeMemberForDate("03", "1990-01-01", "2024-03-15", "2029-03-15")
		if err := validateDocumentDates(m); err == nil {
			t.Error("expected error, got nil")
		}
	})

	// ── Type 52 ───────────────────────────────────────────────────────────────

	t.Run("52: 5yr-1day correct", func(t *testing.T) {
		// issue 2024-07-01 → addYearsMinusOneDay(2024-07-01,5)=2029-06-30
		m := makeMemberForDate("52", "1990-01-01", "2024-07-01", "2029-06-30")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// ── Type 04 ───────────────────────────────────────────────────────────────

	t.Run("04: age<16 at issue → 5yr-1day", func(t *testing.T) {
		// born 2012-01-01, issued 2026-01-01 (age=14) → addYearsMinusOneDay(2026-01-01,5)=2030-12-31
		m := makeMemberForDate("04", "2012-01-01", "2026-01-01", "2030-12-31")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("04: age>=16 at issue → 10yr-1day", func(t *testing.T) {
		// born 2004-01-01, issued 2024-01-01 (age=20) → addYearsMinusOneDay(2024-01-01,10)=2033-12-31
		m := makeMemberForDate("04", "2004-01-01", "2024-01-01", "2033-12-31")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// ── Type 05 ───────────────────────────────────────────────────────────────

	t.Run("05: within 10yr ok", func(t *testing.T) {
		m := makeMemberForDate("05", "1990-01-01", "2024-06-01", "2030-06-01")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("05: exactly 10yr ok", func(t *testing.T) {
		m := makeMemberForDate("05", "1990-01-01", "2026-01-01", "2036-01-01")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("05: exceeds 10yr → error", func(t *testing.T) {
		m := makeMemberForDate("05", "1990-01-01", "2026-01-01", "2036-01-02")
		err := validateDocumentDates(m)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "10年") {
			t.Errorf("error should mention 10年, got: %v", err)
		}
	})

	// ── Basic validations ─────────────────────────────────────────────────────

	t.Run("birth date in future → error", func(t *testing.T) {
		m := makeMemberForDate("05", futureDate(365), "2024-06-01", "2030-06-01")
		err := validateDocumentDates(m)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "出生日期") {
			t.Errorf("error should mention 出生日期, got: %v", err)
		}
	})

	t.Run("issue date in future → error", func(t *testing.T) {
		m := makeMemberForDate("05", "1990-01-01", futureDate(30), "2035-01-01")
		err := validateDocumentDates(m)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "签发日期") {
			t.Errorf("error should mention 签发日期, got: %v", err)
		}
	})

	t.Run("expiry date in past → error", func(t *testing.T) {
		m := makeMemberForDate("05", "1990-01-01", "2020-01-01", "2021-01-01")
		err := validateDocumentDates(m)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "晚于今日") {
			t.Errorf("error should mention 晚于今日, got: %v", err)
		}
	})

	t.Run("feb29 leap year edge: type 01 feb29 issue +10yr → non-leap → mar1", func(t *testing.T) {
		// born 2004-02-29 (age 20 at issue 2024-02-29): ageAtDate = 19 (before Mar 1)
		// → 16 <= 19 <= 25 → 10yr → addSameDayYears(2024-02-29, 10)
		// 2034 is not a leap year → expiry = 2034-03-01
		m := makeMemberForDate("01", "2004-02-29", "2024-02-29", "2034-03-01")
		if err := validateDocumentDates(m); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

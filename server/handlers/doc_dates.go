package handlers

import (
	"fmt"
	"time"

	"fmps/models"
)

// isLeapYear returns true if year is a leap year.
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// ageAtDate returns the completed age (in years) of a person born on birthDate as of targetDate.
// For Feb 29 births, the birthday is ALWAYS treated as Mar 1 in ALL years (not just non-leap).
func ageAtDate(birthDate, targetDate time.Time) int {
	years := targetDate.Year() - birthDate.Year()
	var effectiveBirthday time.Time
	if birthDate.Month() == time.February && birthDate.Day() == 29 {
		effectiveBirthday = time.Date(targetDate.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
	} else {
		effectiveBirthday = time.Date(targetDate.Year(), birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, time.UTC)
	}
	if targetDate.Before(effectiveBirthday) {
		years--
	}
	return years
}

// addSameDayYears returns issueDate + years years with the same month/day.
// If the result would be Feb 29 but that year has no Feb 29, returns Mar 1 that year.
func addSameDayYears(issueDate time.Time, years int) time.Time {
	targetYear := issueDate.Year() + years
	if issueDate.Month() == time.February && issueDate.Day() == 29 && !isLeapYear(targetYear) {
		return time.Date(targetYear, time.March, 1, 0, 0, 0, 0, time.UTC)
	}
	return time.Date(targetYear, issueDate.Month(), issueDate.Day(), 0, 0, 0, 0, time.UTC)
}

// addYearsMinusOneDay returns issueDate + years years - 1 day.
func addYearsMinusOneDay(issueDate time.Time, years int) time.Time {
	return issueDate.AddDate(years, 0, -1)
}

// validateDocumentDates validates birth date, issue date, expiry date, and type-specific
// expiry rules for a member.
func validateDocumentDates(m *models.Member) error {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// ── Birth date ────────────────────────────────────────────────────────────
	birthDate, err := time.Parse("2006-01-02", m.BirthDate)
	if err != nil {
		return fmt.Errorf("出生日期格式无效（需为 YYYY-MM-DD）")
	}
	birthDate = birthDate.UTC()
	if !birthDate.Before(today) {
		return fmt.Errorf("出生日期必须早于今日")
	}
	if ageAtDate(birthDate, today) > 100 {
		return fmt.Errorf("出生日期不得早于100年前")
	}

	// ── Issue date ────────────────────────────────────────────────────────────
	issueDate, err := time.Parse("2006-01-02", m.IdIssueDate)
	if err != nil {
		return fmt.Errorf("签发日期格式无效（需为 YYYY-MM-DD）")
	}
	issueDate = issueDate.UTC()
	if issueDate.After(today) {
		return fmt.Errorf("签发日期不得晚于今日")
	}
	issuedYearsAgo := today.Year() - issueDate.Year()
	if today.Before(time.Date(today.Year(), issueDate.Month(), issueDate.Day(), 0, 0, 0, 0, time.UTC)) {
		issuedYearsAgo--
	}
	if issuedYearsAgo > 20 {
		return fmt.Errorf("签发日期不得早于20年前")
	}

	// ── Expiry date ───────────────────────────────────────────────────────────
	expiryDate, err := time.Parse("2006-01-02", m.IdExpiryDate)
	if err != nil {
		return fmt.Errorf("有效期格式无效（需为 YYYY-MM-DD）")
	}
	expiryDate = expiryDate.UTC()

	// Special case: type 01 with long-term (2099-12-31) is validated below
	longTermDate := time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)
	isLongTerm := expiryDate.Equal(longTermDate)

	if !isLongTerm && !expiryDate.After(today) {
		return fmt.Errorf("有效期必须晚于今日")
	}

	// ── Type-specific expiry validation ───────────────────────────────────────
	docType := m.IdDocType

	switch docType {
	case "01":
		ageAtIssue := ageAtDate(birthDate, issueDate)
		var expectedExpiry time.Time
		if ageAtIssue >= 46 {
			expectedExpiry = longTermDate
		} else if ageAtIssue >= 26 {
			expectedExpiry = addSameDayYears(issueDate, 20)
		} else if ageAtIssue >= 16 {
			expectedExpiry = addSameDayYears(issueDate, 10)
		} else {
			expectedExpiry = addSameDayYears(issueDate, 5)
		}
		if !expiryDate.Equal(expectedExpiry) {
			if ageAtIssue >= 46 {
				return fmt.Errorf("01类证件（签发时年满46周岁）有效期应为长期（2099-12-31）")
			}
			return fmt.Errorf("01类证件有效期应为 %s（按签发时年龄 %d 岁计算）", expectedExpiry.Format("2006-01-02"), ageAtIssue)
		}

	case "91":
		// Expiry = 16th birthday of the holder
		// Since birth year + 16 may or may not be a leap year, and Feb 29 births are
		// guaranteed to have a leap year at birth+16 (only if birth year is leap and
		// birth+16 is also leap – but actually not guaranteed, so we follow addSameDayYears logic).
		expiry91 := time.Date(birthDate.Year()+16, birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, time.UTC)
		if birthDate.Month() == time.February && birthDate.Day() == 29 && !isLeapYear(birthDate.Year()+16) {
			expiry91 = time.Date(birthDate.Year()+16, time.March, 1, 0, 0, 0, 0, time.UTC)
		}
		if ageAtDate(birthDate, today) >= 16 {
			return fmt.Errorf("91证件持有人已满16周岁，不得录入91类证件")
		}
		if !expiryDate.Equal(expiry91) {
			return fmt.Errorf("91类证件有效期应为持有人16周岁生日（%s）", expiry91.Format("2006-01-02"))
		}

	case "11", "21":
		expectedExpiry := addSameDayYears(issueDate, 5)
		if !expiryDate.Equal(expectedExpiry) {
			return fmt.Errorf("%s类证件有效期应为 %s（签发日期加5年）", docType, expectedExpiry.Format("2006-01-02"))
		}

	case "31":
		ageAtIssue := ageAtDate(birthDate, issueDate)
		years := 10
		if ageAtIssue < 18 {
			years = 5
		}
		expectedExpiry := addYearsMinusOneDay(issueDate, years)
		if !expiryDate.Equal(expectedExpiry) {
			return fmt.Errorf("31类证件有效期应为 %s（签发时年龄 %d 岁，%d年有效期）", expectedExpiry.Format("2006-01-02"), ageAtIssue, years)
		}

	case "02":
		ageAtIssue := ageAtDate(birthDate, issueDate)
		years := 10
		if ageAtIssue < 18 {
			years = 5
		}
		expectedExpiry := addYearsMinusOneDay(issueDate, years)
		if !expiryDate.Equal(expectedExpiry) {
			return fmt.Errorf("02类证件有效期应为 %s（签发时年龄 %d 岁，%d年有效期）", expectedExpiry.Format("2006-01-02"), ageAtIssue, years)
		}

	case "03", "52":
		expectedExpiry := addYearsMinusOneDay(issueDate, 5)
		if !expiryDate.Equal(expectedExpiry) {
			return fmt.Errorf("%s类证件有效期应为 %s（签发日期加5年减1天）", docType, expectedExpiry.Format("2006-01-02"))
		}

	case "04":
		ageAtIssue := ageAtDate(birthDate, issueDate)
		years := 10
		if ageAtIssue < 16 {
			years = 5
		}
		expectedExpiry := addYearsMinusOneDay(issueDate, years)
		if !expiryDate.Equal(expectedExpiry) {
			return fmt.Errorf("04类证件有效期应为 %s（签发时年龄 %d 岁，%d年有效期）", expectedExpiry.Format("2006-01-02"), ageAtIssue, years)
		}

	case "05":
		maxExpiry := addSameDayYears(issueDate, 10)
		if expiryDate.After(maxExpiry) {
			return fmt.Errorf("05类外国护照有效期最长10年（最晚应为 %s）", maxExpiry.Format("2006-01-02"))
		}
	}

	return nil
}

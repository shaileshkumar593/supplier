package param

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	// FilterRegexWord is a regex for letters, numbers and space
	FilterRegexWord = `([a-zA-Z0-9.]+[\s]?)+`
	// FilterRegexWordWithoutSpace is a regex for letters, numbers without space
	FilterRegexWordWithoutSpace = `[a-zA-Z0-9]+`
	// FilterRegexMatchAll match all
	FilterRegexMatchAll = `([^*]+)`
)

// ParseOrder converts a csv string (ex. id,-name,dateAdded) to an SQL ORDER BY list
func ParseOrder(s string, conv map[string]string) ([]string, error) {
	orderRule := "ASC"
	submatch := `(\-)?([a-zA-Z_0-9]+)`
	fullMatch := `^(` + submatch + `,)*` + submatch + `$`

	if s == "" {
		return make([]string, 0), nil
	}

	if m, err := regexp.MatchString(fullMatch, s); err != nil || !m {
		return nil, fmt.Errorf(`invalid string format or error (%v) encountered`, err)
	}

	rgx := regexp.MustCompile(submatch)
	matches := rgx.FindAllStringSubmatch(s, -1)
	var res []string

	for _, m := range matches {
		if _, ok := conv[m[2]]; !ok {
			return nil, fmt.Errorf(`unrecognized %s in allowed map`, m[2])
		}

		if m[1] == "-" {
			orderRule = "DESC"
		}
		res = append(res, conv[m[2]]+" "+orderRule)
	}
	return res, nil
}

// ParseFilter converts a string to an SQL LIKE parameter with
// an asterisk (*) wilcard which can be placed on either or both ends
func ParseFilter(s string, pat string) (string, error) {
	submatch := `(\*?(` + pat + `)\*?)`
	fullMatch := `^` + submatch + `$`

	if s == "" {
		return "", nil
	}

	if m, err := regexp.MatchString(fullMatch, s); err != nil || !m {
		return "", fmt.Errorf(`invalid string format or error (%v) encountered`, err)
	}

	rgx := regexp.MustCompile(submatch)
	m := rgx.FindAllStringSubmatch(s, 1)
	fm := m[0][0]
	rm := m[0][2]

	if fm[:1] == "*" {
		rm = "%" + rm
	}

	if fm[len(fm)-1:] == "*" {
		rm = rm + "%"
	}

	return rm, nil
}

// ParseDateRange converts a given date range to an valid SQL BETWEEN parameter list
func ParseDateRange(s string) (dates []string, err error) {
	var (
		tempDate = []string{
			"00:00:00",
			"23:59:59",
		}
	)

	dates = strings.Split(s, ",")
	if 2 != len(dates) {
		err = fmt.Errorf(`invalid dateRange format`)
		return
	}

	for i, date := range dates {
		if _, err = time.Parse("2006-01-02", date); nil != err {
			err = fmt.Errorf(`error (%v) encountered`, err)
			return
		}

		dates[i] = "" + date + " " + tempDate[i] + ""
	}

	return
}

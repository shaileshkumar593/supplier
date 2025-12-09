package pagination

import (
	"context"
	"errors"
	"math"
	"net/url"
	"strconv"
	"strings"

	"swallow-supplier/config"
	customContext "swallow-supplier/context"
	"swallow-supplier/utils/array"
)

// Request represents the common pagination request
type Request struct {
	Page           string
	RecordsPerPage string
	Condition      string
	SortDirection  string
}

// QueryFilter represents the search filter fields
type QueryFilter struct {
	Key   string
	Value string
}

// QueryDetail represents the condition and sorting details
type QueryDetail struct {
	Page           string
	RecordsPerPage string
	Condition      string
	SortOrder      string
	SortParams     []string
}

// Pagination represents the pagination field
type Pagination struct {
	TotalRecords   int64  `json:"total_records,omitempty"`
	RecordsPerPage int    `json:"records_per_page,omitempty"`
	TotalPages     int    `json:"total_pages,omitempty"`
	Links          *Links `json:"links,omitempty"`
}

// Links represents the pagination links
type Links struct {
	Self     string `json:"self,omitempty"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

// Reports represents reporting response
type Reports struct {
	Pagination *Pagination `json:"pagination"`
	Data       interface{} `json:"data"`
}

const (
	// QueryPage represents the page param
	QueryPage = "page"
	// QueryRecordsPerPage represents the records_per_page param
	QueryRecordsPerPage = "records_per_page"

	// DefaultRecordsPerPage default value of number of records per page
	DefaultRecordsPerPage = "10"
	// DefaultPageCount default value of page
	DefaultPageCount = "1"
	// MaxRecordsPerPage max number of records per page
	MaxRecordsPerPage = "50"

	// DefaultCondition default query condition
	DefaultCondition = "AND"
	// DefaultSortOrder default sorting
	DefaultSortOrder = "ASC"
)

var (
	// Conditions slice of supported query conditions
	Conditions = []string{
		"AND",
		"OR",
	}

	// SortDirection slice of supported sorting conditions
	SortDirection = []string{
		"ASC",
		"DESC",
	}

	// IsActiveMapping mappint of is_active values
	IsActiveMapping = map[string]string{
		"true":  "1",
		"false": "0",
		"all":   "",
	}
)

// ValidatePaginationRequest validates common request on reports endpoints
func (p *Request) ValidatePaginationRequest() map[string]error {
	var (
		err              error
		paginationErrors = make(map[string]error)
	)

	// Page validation
	if p.Page != "" {
		page, err := strconv.Atoi(p.Page)
		if err != nil {
			err = errors.New("should be a number")
			paginationErrors["page"] = err
		} else {
			if page < 1 {
				err = errors.New("should not be less than 1")
				paginationErrors["page"] = err
			}
		}
	}

	// Records Per Page validation
	if p.RecordsPerPage != "" {
		recordsPerPage, err := strconv.Atoi(p.RecordsPerPage)
		if err != nil {
			err = errors.New("should be a number")
			paginationErrors["records_per_page"] = err
		} else {
			if recordsPerPage < 1 {
				err = errors.New("should not be less than 1")
				paginationErrors["records_per_page"] = err
			}

			maxRecordsPerPage, _ := strconv.Atoi(MaxRecordsPerPage)
			if recordsPerPage > maxRecordsPerPage {
				err = errors.New("allowed max value is 50")
				paginationErrors["records_per_page"] = err
			}
		}
	}

	// Condition validation
	if p.Condition != "" {
		if exists, _ := array.InArray(strings.ToUpper(p.Condition), Conditions); !exists {
			err = errors.New("unsupported value")
			paginationErrors["condition"] = err
		}
	}

	// Sort Direction Validation
	if p.SortDirection != "" {
		if exists, _ := array.InArray(strings.ToUpper(p.SortDirection), SortDirection); !exists {
			err = errors.New("unsupported value")
			paginationErrors["sort_order"] = err
		}
	}

	return paginationErrors
}

// FormatResponse formats the pagination and data object on the final response for reports
func FormatResponse(ctx context.Context, data interface{}, count int64) Reports {
	var r Reports
	if ctx.Value(customContext.CtxLabelRequestURL) != nil {
		requestURL := ctx.Value(customContext.CtxLabelRequestURL).(*url.URL)
		r.Pagination = FormatPagination(requestURL, count)
	}

	if data != nil {
		r.Data = data
	}

	return r
}

// FormatPagination creates and formats the pagination object
func FormatPagination(url *url.URL, totalRecords int64) (p *Pagination) {
	var (
		currentPage    int
		recordsPerPage int
		err            error

		page = Pagination{}
	)

	// Retrieve the current page from query parameters
	currentPage, _ = strconv.Atoi(url.Query().Get(QueryPage))
	if currentPage == 0 {
		currentPage, _ = strconv.Atoi(DefaultPageCount)
	}

	// Retrieve the records per page from query parameters
	if recordsPerPage, err = strconv.Atoi(url.Query().Get(QueryRecordsPerPage)); err != nil || recordsPerPage == 0 {
		recordsPerPage, _ = strconv.Atoi(DefaultRecordsPerPage)
	}

	page.TotalRecords = totalRecords
	page.RecordsPerPage = recordsPerPage
	page.TotalPages = int(math.Ceil(float64(totalRecords) / float64(recordsPerPage)))

	path := config.Instance().AppDomain + url.Path + "?"
	query := url.Query()

	l := &Links{}
	if page.TotalRecords == 0 {
		l = nil
	} else {
		query.Set(QueryPage, strconv.Itoa(currentPage))
		query.Set(QueryRecordsPerPage, strconv.Itoa(recordsPerPage))
		l.Self = path + query.Encode()

		if page.TotalPages > 0 {
			if (currentPage + 1) <= page.TotalPages {
				query.Set(QueryPage, strconv.Itoa(currentPage+1))
				l.Next = path + query.Encode()
			}

			if (currentPage - 1) != 0 {
				query.Set(QueryPage, strconv.Itoa(currentPage-1))
				l.Previous = path + query.Encode()
			}
		}
	}

	page.Links = l
	p = &page

	return
}

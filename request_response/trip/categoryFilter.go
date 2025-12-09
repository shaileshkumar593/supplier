package trip

// Request to mongo repository to filter Trip based on CategoryCode, CategoryLevel
type CategoryFilter struct {
	CategoryCode  int
	CategoryLevel int
}

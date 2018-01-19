package analyze

// Match encapsulates an analyze result by containing the starting row and column, as well as the confidence score
type Match struct {
	Row        int
	Col        int
	Confidence float64
}

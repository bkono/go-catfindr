package analyze

import (
	"bufio"
	"io"
)

// Frame encapsulates a single txt matrix encoded as [][]byte for analysis
type Frame struct {
	rows    int
	cols    int
	encoded [][]byte
}

// Rows is the count of rows in the frame
func (f *Frame) Rows() int {
	return f.rows
}

// Cols is the count of columns in the Frame
func (f *Frame) Cols() int {
	return f.cols
}

// Pixels is the total number of points available in the Frame
func (f *Frame) Pixels() int {
	return f.rows * f.cols
}

// Encoded is the underlying [][]byte representing the entire Frame
func (f *Frame) Encoded() [][]byte {
	return f.encoded
}

// CharAt retrieves a single byte character from the encoded Frame
func (f *Frame) CharAt(row, col int) byte {
	return f.encoded[row][col]
}

// NewFrameFromReader reads an io.Reader into a new Frame struct
func NewFrameFromReader(r io.Reader) (*Frame, error) {
	res := [][]byte{}
	br := bufio.NewReader(r)
	maxCols := 0
	for {
		line, err := br.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		line = line[:len(line)-1]
		if maxCols < len(line) {
			maxCols = len(line)
		}
		res = append(res, line)
	}

	// normalize the column widths to ensure a true grid
	for i, v := range res {
		diff := maxCols - len(v)
		if diff > 0 {
			for j := 0; j < diff; j++ {
				res[i] = append(res[i], '\u0000')
			}
		}
	}

	return &Frame{cols: maxCols, rows: len(res), encoded: res}, nil
}

// FindMatches analyzes the input frame for matches to the known frame meeting the provided minimum confidence score
func FindMatches(in, known *Frame, minConfidence float64) []*Match {
	var matches []*Match
	for row := 0; row < in.Rows()-known.Rows(); row++ {
		for col := 0; col < in.Cols()-known.Cols(); col++ {
			chars := matchedChars(row, col, in, known)
			score := float64(chars) / float64(known.Pixels())
			if score >= minConfidence {
				matches = append(matches, &Match{Row: row, Col: col, Confidence: score})
			}
		}
	}
	return matches
}

func matchedChars(row, col int, in, known *Frame) int {
	matches := 0
	for trow := 0; trow < known.Rows(); trow++ {
		for tcol := 0; tcol < known.Cols(); tcol++ {
			if in.CharAt(trow+row, tcol+col) == known.CharAt(trow, tcol) {
				matches++
			}
		}
	}
	return matches
}

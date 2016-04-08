package mat

import (
	"errors"
)

// MatrixFloat64 is a naive matrix construct, good for small
// matric operations
type MatrixFloat64 struct {
	R           int
	C           int
	Values      []float64
	Determinant float64
}

func NewMatrixFloat64(r, c int, val []float64) *MatrixFloat64 {
	if r < 1 || c < 1 {
		return nil
	}

	size := r * c
	m := new(MatrixFloat64)
	m.R = r
	m.C = c
	m.Values = make([]float64, size, size)
	for k, v := range val {
		m.Values[k] = v
	}

	if r == c {
		m.calculateDeterminant()
	} else {
		m.Determinant = 0
	}
	return m
}

func (m *MatrixFloat64) calculateDeterminant() {
	if m.R == 1 {
		m.Determinant = m.Values[0]
	} else if m.R == 2 {
		m.Determinant = m.Values[0]*m.Values[3] - m.Values[1]*m.Values[2]
	} else if m.R == 3 {
		m.Determinant = (m.Values[0] * m.Values[4] * m.Values[8]) -
			(m.Values[0] * m.Values[5] * m.Values[7]) -
			(m.Values[1] * m.Values[3] * m.Values[8]) +
			(m.Values[1] * m.Values[5] * m.Values[6]) +
			(m.Values[2] * m.Values[3] * m.Values[7]) -
			(m.Values[2] * m.Values[4] * m.Values[6])
	} else {
		// TODO: Implement mechanism to calculate det for any size matrix
	}
}

func calculateDet2x2(val []float64) float64 {
	return val[0]*val[3] - val[1]*val[2]
}

func calculateDet(multiplier float64, submatrix []float64) {
	// TODO: Finish this
}

func (m *MatrixFloat64) Get(i, j int) (float64, error) {
	if i < 0 || i > m.R || j < 0 || j > m.C {
		return 0, errors.New("rows or columns out of bounds")
	}

	index := j*m.C + i
	return m.Values[index], nil
}

func (m *MatrixFloat64) Set(i, j int, v float64) error {
	if i < 0 || i > m.R || j < 0 || j > m.C {
		return errors.New("rows or columns out of bounds")
	}

	index := j*m.C + i
	m.Values[index] = v
	return nil
}

func (m *MatrixFloat64) Add(m2 *MatrixFloat64) error {
	if m.R != m2.R || m.C != m2.C {
		return errors.New("Matricies must be the same size")
	}

	for k := range m.Values {
		m.Values[k] += m2.Values[k]
	}
	return nil
}

func (m *MatrixFloat64) Subtract(m2 *MatrixFloat64) error {
	if m.R != m2.R || m.C != m2.C {
		return errors.New("Matricies must be the same size")
	}
	for k := range m.Values {
		m.Values[k] -= m2.Values[k]
	}
	return nil
}

func (m *MatrixFloat64) MultiplyScalar(s float64) {
	for k := range m.Values {
		m.Values[k] *= s
	}
}

func (m *MatrixFloat64) DivideScalar(s float64) {
	inv := 1 / s
	for k := range m.Values {
		m.Values[k] *= inv
	}
}

func Multiply(m1, m2 *MatrixFloat64) (*MatrixFloat64, error) {
	if m1.C != m2.R {
		return nil, errors.New("Columns size of matrix 1 one must match row size of matrix 2")
	}

	size := m1.R * m2.C
	values := make([]float64, 0, size)

	var index int
	var row_min int
	var row_max int
	var col_increment int

	for i := 0; i < m1.R; i++ {
		for j := 0; j < m2.C; j++ {
			index = i*m2.C + j
			row_min = i * m1.C
			row_max = (i + 1) * m1.C
			col_increment = m2.C
			values[index] = RowXCol(m1.Values, m2.Values, row_min, row_max, j, col_increment)
		}
	}

	m3 := NewMatrixFloat64(m1.R, m2.C, values)
	return m3, nil
}

func RowXCol(r, c []float64, row_min, row_max, col_min, col_iterator int) float64 {
	val := 0.0
	for k := range r {
		val += r[k] * c[k]
	}
	return val
}

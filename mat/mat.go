package mat

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/agdt3/goray/vec"
)

// MatrixFloat64 is a naive matrix construct,
// good for small matrix operations
type MatrixFloat64 struct {
	R           int
	C           int
	values      []float64
	Determinant float64
}

// NewMatrixFloat64 is a constructor for an r x c naive matrix
func NewMatrixFloat64(r, c int, val []float64) *MatrixFloat64 {
	if r < 1 || c < 1 {
		return nil
	}

	size := r * c
	m := new(MatrixFloat64)
	m.R = r
	m.C = c
	m.values = make([]float64, size, size)
	for k, v := range val {
		m.values[k] = v
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
		m.Determinant = m.values[0]
	} else if m.R == 2 {
		m.Determinant = m.values[0]*m.values[3] - m.values[1]*m.values[2]
	} else if m.R == 3 {
		m.Determinant = (m.values[0] * m.values[4] * m.values[8]) -
			(m.values[0] * m.values[5] * m.values[7]) -
			(m.values[1] * m.values[3] * m.values[8]) +
			(m.values[1] * m.values[5] * m.values[6]) +
			(m.values[2] * m.values[3] * m.values[7]) -
			(m.values[2] * m.values[4] * m.values[6])
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

// Get returns value of matrix at row i, column j (indexed at 0)
func (m *MatrixFloat64) Get(i, j int) (float64, error) {
	if i < 0 || i >= m.R || j < 0 || j >= m.C {
		return 0, errors.New("rows or columns out of bounds")
	}

	index := i*m.C + j
	return m.values[index], nil
}

// Set changes values of matrix at row i, column j (indexed at 0)
func (m *MatrixFloat64) Set(i, j int, v float64) error {
	if i < 0 || i > m.R || j < 0 || j > m.C {
		return errors.New("rows or columns out of bounds")
	}

	index := i*m.C + j
	m.values[index] = v
	m.calculateDeterminant()
	return nil
}

// Add adds one matrix to the given matrix
func (m *MatrixFloat64) Add(m2 *MatrixFloat64) error {
	if m.R != m2.R || m.C != m2.C {
		return errors.New("Matricies must be the same size")
	}

	for k := range m.values {
		m.values[k] += m2.values[k]
	}
	return nil
}

// Subtract convenience methods to subtract one matrix from the current one
func (m *MatrixFloat64) Subtract(m2 *MatrixFloat64) error {
	if m.R != m2.R || m.C != m2.C {
		return errors.New("Matricies must be the same size")
	}
	for k := range m.values {
		m.values[k] -= m2.values[k]
	}
	return nil
}

// MultiplyScalar multiplies the matrix by a scalar
func (m *MatrixFloat64) MultiplyScalar(s float64) {
	for k := range m.values {
		m.values[k] *= s
	}
}

// DivideScalar divides the matrix by a scalar
func (m *MatrixFloat64) DivideScalar(s float64) {
	inv := 1 / s
	m.MultiplyScalar(inv)
}

// Multiply multiplies two matricies and returns a third. This is an O(n^3) operation.
func Multiply(m1, m2 *MatrixFloat64) (*MatrixFloat64, error) {
	if m1.C != m2.R {
		return nil, errors.New("Columns size of matrix 1 one must match row size of matrix 2")
	}

	size := m1.R * m2.C
	values := make([]float64, size, size)

	// O(n^3)
	var sum float64
	for i := 0; i < m1.R; i++ {
		for j := 0; j < m2.C; j++ {
			sum = 0
			for k := 0; k < m2.R; k++ {
				vM1, err1 := m1.Get(i, k)
				vM2, err2 := m2.Get(k, j)
				if err1 == nil && err2 == nil {
					sum += vM1 * vM2
				}
			}
			index := i*m2.C + j
			values[index] = sum
		}
	}

	m3 := NewMatrixFloat64(m1.R, m2.C, values)
	return m3, nil
}

// MultiplyVec3 multiplies a matrix by a column vector
func MultiplyVec3(m *MatrixFloat64, v *vec.Vec3) (*vec.Vec3, error) {
	if m.R != 3 || m.C != 3 {
		return nil, errors.New("Must be a 3 x 3 matrix")
	}

	values := make([]float64, 3, 3)

	for i := 0; i < m.R; i++ {
		mX, err1 := m.Get(i, 0)
		mY, err2 := m.Get(i, 1)
		mZ, err3 := m.Get(i, 2)
		if err1 == nil && err2 == nil && err3 == nil {
			values[i] = mX*v.X + mY*v.Y + mZ*v.Z
		}
	}

	return vec.NewVec3(values[0], values[1], values[2]), nil
}

// String returns a semi-formatted representation of the matrix
func (m MatrixFloat64) String() string {
	var str string
	str += "R: %v C: %v Det: %v \n"
	for i := 0; i < m.R; i++ {
		for j := 0; j < m.C; j++ {
			v, _ := m.Get(i, j)
			str += strconv.FormatFloat(v, 'f', -1, 64)
			str += " "
		}
		str += "\n"
	}
	return fmt.Sprintf(str, m.R, m.C, m.Determinant)
}

// IsEqual compares Rows, Columns and values within tolerance
func IsEqual(m1, m2 *MatrixFloat64, tolerance float64) bool {
	if m1.R != m2.R || m1.C != m2.C {
		return false
	}

	for i := 0; i < m1.R; i++ {
		for j := 0; j < m1.C; j++ {
			v1, err1 := m1.Get(i, j)
			v2, err2 := m2.Get(i, j)
			if err1 != nil ||
				err2 != nil ||
				math.Abs(v1-v2) >= tolerance {
				return false
			}
		}
	}

	return true
}

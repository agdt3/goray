package files

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/agdt3/goray/obj"
)

// ReadMeshFiles reads a lot of mesh files and creates
// PolygonMesh objects out of each one
// TODO: Make this not flat
func ReadMeshFiles(directory string) error {
	dir, err := os.Open(directory)
	if err != nil {
		dir.Close()
		return err
	}

	//files
	return nil
}

// ReadMeshFile reads .mesh file (via path) into struct
func ReadMeshFile(path string, polygon *obj.PolygonMesh) error {
	file, err := os.Open(path)
	if err != nil {
		file.Close()
		return err
	}

	var numFaces []int
	var numVerticies []int
	var vertexIndecies []int
	var verticies []float64

	reader := bufio.NewReader(file)

	if line, err := ReadMeshLine(reader); err == nil {
		numFaces = StringToIntArray(line)
	} else {
		return err
	}

	if line, err := ReadMeshLine(reader); err == nil {
		numVerticies = StringToIntArray(line)
	} else {
		return err
	}

	if line, err := ReadMeshLine(reader); err == nil {
		vertexIndecies = StringToIntArray(line)
	} else {
		return err
	}

	if line, err := ReadMeshLine(reader); err == nil {
		verticies = StringToFloat64Array(line)
	} else {
		return err
	}
	//vertex_normals := StringToFloat64Array(ReadMeshLine(reader))

	polygon.NumFaces = numFaces
	polygon.NumVerticies = numVerticies
	polygon.VertexIndecies = vertexIndecies
	polygon.Verticies = verticies

	file.Close()
	return nil
}

// ReadMeshLine reads one line from a reader
// and separates the line on " ", returning a string array
func ReadMeshLine(reader *bufio.Reader) ([]string, error) {
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	lineArr := strings.Split(string(line), " ")
	return lineArr, nil
}

// StringToIntArray converts a string to an int array
func StringToIntArray(strArr []string) []int {
	intArr := make([]int, len(strArr), len(strArr))
	for i, v := range strArr {
		fv, _ := strconv.Atoi(v)
		intArr[i] = fv
	}
	return intArr
}

// StringToFloat64Array converts a string to a float64 array
func StringToFloat64Array(strArr []string) []float64 {
	floatArr := make([]float64, len(strArr), len(strArr))
	for i, v := range strArr {
		fv, _ := strconv.ParseFloat(v, 64)
		floatArr[i] = fv
	}
	return floatArr
}

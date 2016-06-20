package files

import (
	"bufio"
	"fmt"
	"io"
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

	if line, err := readMeshLine(reader); err == nil {
		numFaces = stringToIntArray(line)
	} else {
		return err
	}

	if line, err := readMeshLine(reader); err == nil {
		numVerticies = stringToIntArray(line)
	} else {
		return err
	}

	if line, err := readMeshLine(reader); err == nil {
		vertexIndecies = stringToIntArray(line)
	} else {
		return err
	}

	if line, err := readMeshLine(reader); err == nil {
		verticies = stringToFloat64Array(line)
	} else {
		return err
	}
	//vertex_normals := StringToFloat64Array(readMeshLine(reader))

	polygon.NumFaces = numFaces
	polygon.NumVerticies = numVerticies
	polygon.VertexIndecies = vertexIndecies
	polygon.Verticies = verticies

	file.Close()
	return nil
}

// readMeshLine reads one line from a reader
// and separates the line on " ", returning a string array
func readMeshLine(reader *bufio.Reader) ([]string, error) {
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	lineArr := strings.Split(string(line), " ")
	return lineArr, nil
}

// stringToIntArray converts a string to an int array
func stringToIntArray(strArr []string) []int {
	intArr := make([]int, len(strArr), len(strArr))
	for i, v := range strArr {
		v = strings.TrimSuffix(v, "\n")
		if v == "" {
			intArr[i] = 0
		} else {
			fv, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println(err)
			}
			intArr[i] = fv
		}
	}
	return intArr
}

// stringToFloat64Array converts a string to a float64 array
func stringToFloat64Array(strArr []string) []float64 {
	floatArr := make([]float64, len(strArr), len(strArr))
	for i, v := range strArr {
		v = strings.TrimSuffix(v, "\n")
		fv, _ := strconv.ParseFloat(v, 64)
		floatArr[i] = fv
	}
	return floatArr
}

// ReadWavFile reads in wavefront .obj file and processes it
func ReadWavFile(path string, polygon *obj.PolygonMesh) error {
	file, err := os.Open(path)
	if err != nil {
		file.Close()
		return err
	}

	reader := bufio.NewReader(file)

	for true {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		} else {
			parseWavLine(line, " ", polygon)
		}
	}

	return nil
}

func parseWavLine(line, sep string, poly *obj.PolygonMesh) {
	results := strings.Split(line, sep)
	prefix := results[0]
	values := results[1:]

	switch prefix {
	case "v":
		verticies := stringToFloat64Array(values)
		poly.Verticies = append(poly.Verticies, verticies...)
		//poly.NumVerticies = append(poly.NumVerticies, len(verticies))
	case "vt":
		fmt.Println("a texture vertex")
		fmt.Println(stringToFloat64Array(values))
	case "vn":
		fmt.Println("a vertex normal")
		fmt.Println(stringToFloat64Array(values))
	case "f":
		poly.NumFaces[0] += 1
		vertexIndecies, _, _ := extractIndecies(values)
		poly.VertexIndecies = append(poly.VertexIndecies, vertexIndecies...)
		poly.NumVerticies = append(poly.NumVerticies, len(vertexIndecies))
	default:
		//Do nothing
	}
}

func extractIndecies(faceValues []string) ([]int, []int, []int) {
	size := len(faceValues)
	vertexIndecies := make([]int, 0, size)
	textureIndecies := make([]int, 0, size)
	vertexNormalIndecies := make([]int, 0, size)

	// cases:
	// v/vt/vn
	// v//vn
	// v/vt
	// v
	for _, v := range faceValues {
		valuesString := strings.Split(v, "/")
		valuesInt := stringToIntArray(valuesString)
		if len(valuesString) == 3 && valuesString[1] != "" {
			// vertexIndecies are normalized to start at 0
			// blender export starts them at 1
			vertexIndecies = append(vertexIndecies, valuesInt[0]-1)
			textureIndecies = append(textureIndecies, valuesInt[1])
			vertexNormalIndecies = append(vertexNormalIndecies, valuesInt[2])
		} else if len(valuesString) == 3 && valuesString[1] == "" {
			vertexIndecies = append(vertexIndecies, valuesInt[0]-1)
			vertexNormalIndecies = append(vertexNormalIndecies, valuesInt[2])
		} else if len(valuesString) == 2 {
			vertexIndecies = append(vertexIndecies, valuesInt[0]-1)
			textureIndecies = append(textureIndecies, valuesInt[1])
		} else {
			vertexIndecies = append(vertexIndecies, valuesInt[0]-1)
		}
	}

	return vertexIndecies, textureIndecies, vertexNormalIndecies
}

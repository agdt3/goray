package read

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"github.com/agdt3/goray/obj"
	"github.com/agdt3/goray/vec"
)

func ReadMeshFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		file.Close()
		return err
	}

	var num_faces []int
	var num_verticies []int
	var vertex_indecies []int
	var verticies []float64

	reader := bufio.NewReader(file)

	if line, err := ReadMeshLine(reader); err != nil {
		return err
	} else {
		num_faces = StringToIntArray(line)
	}

	if line, err := ReadMeshLine(reader); err != nil {
		return err
	} else {
		num_verticies = StringToIntArray(line)
	}

	if line, err := ReadMeshLine(reader); err != nil {
		return err
	} else {
		vertex_indecies = StringToIntArray(line)
	}

	if line, err := ReadMeshLine(reader); err != nil {
		return err
	} else {
		verticies = StringToFloat64Array(line)
	}

	//vertex_normals := StringToFloat64Array(ReadMeshLine(reader))

	num_triangles := num_verticies[0] - 2
	for i := 0; i < num_faces[0]; i++ {
		for j := 0; j < num_triangles; j++ {
			vertex_index_offset := j * 3
			v0 := vec.NewVec3(
				verticies[vertex_indecies[vertex_index_offset]],
				verticies[vertex_indecies[vertex_index_offset]+1],
				verticies[vertex_indecies[vertex_index_offset]+2])

			v1 := vec.NewVec3(
				verticies[vertex_indecies[vertex_index_offset+1]],
				verticies[vertex_indecies[vertex_index_offset+1]+1],
				verticies[vertex_indecies[vertex_index_offset+1]+2])

			v2 := vec.NewVec3(
				verticies[vertex_indecies[vertex_index_offset+2]],
				verticies[vertex_indecies[vertex_index_offset+2]+1],
				verticies[vertex_indecies[vertex_index_offset+2]+2])

			triangle := obj.NewTriangle("", *v0, *v1, *v2, color.RGBA{255, 0, 0, 1}, 1.0, 1.0, true)
			fmt.Println(triangle)
		}
	}

	file.Close()
	return nil
}

func ReadMeshLine(reader *bufio.Reader) ([]string, error) {
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	line_arr := strings.Split(string(line), " ")
	return line_arr, nil
}

func StringToIntArray(str_arr []string) []int {
	int_arr := make([]int, len(str_arr), len(str_arr))
	for i, v := range str_arr {
		fv, _ := strconv.Atoi(v)
		int_arr[i] = fv
	}
	return int_arr
}

func StringToFloat64Array(str_arr []string) []float64 {
	float_arr := make([]float64, len(str_arr), len(str_arr))
	for i, v := range str_arr {
		fv, _ := strconv.ParseFloat(v, 64)
		float_arr[i] = fv
	}
	return float_arr
}

package scan

import (
	"fmt"
	"strconv"
)

type GraderRes struct {
	Msg   string
	Maso  string
	Made  string
	Anh   string
	Dapan []map[string]string
}

type DownloadRes struct {
	Link string
}

type SessionRes struct {
	Msg string
	Idx string
}

func (g *GraderRes) ToSlice(size int) []string {
	result := make([]string, size)
	result[0] = g.Maso
	result[1] = g.Made
	result[2] = g.Anh

	for _, d := range g.Dapan {
		cau := d["cau"]
		ans := d["answer"]

		i, err := strconv.Atoi(cau)
		if err != nil {
			fmt.Println(err)
			continue
		}

		result[i] = ans
	}

	return result
}

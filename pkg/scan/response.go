package scan

import "strconv"

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

func (g *GraderRes) ToSlice(size int) ([]string, error) {
	result := make([]string, size)
	result[0] = g.Maso
	result[1] = g.Made
	result[2] = g.Anh

	for _, d := range g.Dapan {
		cau := d["cau"]
		ans := d["answer"]

		i, err := strconv.Atoi(cau)
		if err != nil {
			return nil, err
		}

		result[i+2] = ans
	}

	return result, nil
}

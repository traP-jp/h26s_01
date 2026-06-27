package kanjipool

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"strconv"
)

//go:embed kanji.csv
var kanjiCSV []byte

type Kanji struct {
	Char   string
	Stroke int
}

func parseKanji() ([]Kanji, error) {
	r := csv.NewReader(bytes.NewReader(kanjiCSV))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	kanjis := make([]Kanji, 0, len(records)-1)
	for _, record := range records[1:] {
		kanji := record[1]
		stroke, err := strconv.Atoi(record[3])
		if err != nil {
			return nil, err
		}

		kanjis = append(kanjis, Kanji{kanji, stroke})
	}
	return kanjis, nil
}

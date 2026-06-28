package kanjipool

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"math/rand"
	"sort"
	"strconv"
)

const (
	MinStroke    = 11
	MaxStroke    = 17
	MaxCorrect   = 9
	MaxIncorrect = 3
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

func SelectKanjies(playerCount int) ([]Kanji, error) {
	kanjies, err := parseKanji()
	if err != nil {
		return nil, err
	}

	var filtered []Kanji
	for _, kanji := range kanjies {
		if kanji.Stroke >= MinStroke && kanji.Stroke <= MaxStroke {
			filtered = append(filtered, kanji)
		}
	}

	rand.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	requiredNum := MaxCorrect + MaxIncorrect + playerCount - 2

	requiredKanjies := filtered[:requiredNum]
	sort.SliceStable(requiredKanjies, func(i, j int) bool { return requiredKanjies[i].Stroke < requiredKanjies[j].Stroke })
	return requiredKanjies, nil
}

package utils

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

type CSVGradeRow struct {
	RollNumber string
	Grade      string
}

func ParseGradeCSV(r io.Reader) ([]CSVGradeRow, error) {
	reader := csv.NewReader(r)
	header, err := reader.Read()
	if err != nil {
		return nil, errors.New("could not read CSV header")
	}
	if len(header) < 2 ||
		strings.TrimSpace(strings.ToLower(header[0])) != "roll_number" ||
		strings.TrimSpace(strings.ToLower(header[1])) != "grade" {
		return nil, errors.New("CSV must have header: roll_number,grade")
	}

	var rows []CSVGradeRow
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // skip malformed rows; let caller handle
		}
		if len(record) < 2 {
			continue
		}
		rows = append(rows, CSVGradeRow{
			RollNumber: strings.TrimSpace(record[0]),
			Grade:      strings.TrimSpace(strings.ToUpper(record[1])),
		})
	}
	return rows, nil
}

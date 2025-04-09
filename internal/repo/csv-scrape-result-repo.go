package repo

import (
	"encoding/csv"
	"fmt"
	"hw_4_1/internal/entity"
	"os"
	"strconv"
	"time"
)

type CsvScrapeResultRepo struct {
	outputFile string
}

func NewCsvScrapeResultRepo(outputFile string) *CsvScrapeResultRepo {
	return &CsvScrapeResultRepo{outputFile: outputFile}
}

func (r *CsvScrapeResultRepo) SaveResults(results []entity.ScrapeResult) error {
	file, err := os.Create(r.outputFile)
	if err != nil {
		return fmt.Errorf("не удалось создать csv файл: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"Date",
		"Url",
		"StatusCode",
		"Title",
		"Description",
		"SuccessAttempt",
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("не удалось записать заголовки csv файла: %w", err)
	}

	for _, result := range results {
		row := []string{
			result.Date.Format(time.RFC3339),
			result.Url,
			strconv.Itoa(result.StatusCode),
			result.Title,
			result.Description,
			strconv.Itoa(result.SuccessAttempt),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("не удалось записать данные в csv файл: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("не удалось завершить запись данных в csv файл: %w", err)
	}

	return nil
}

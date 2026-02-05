package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"spends/internal/models"
)

// ExportToCSV экспортирует траты в CSV формат
func ExportToCSV(expenses []models.Expense) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ';' // Точка с запятой, как в твоём старом коде

	// Заголовки (без служебной информации)
	header := []string{"дата", "наименование", "сумма", "категория", "приоритет"}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// Данные
	for _, e := range expenses {
		record := []string{
			e.CreatedAt.Format("02.01.2006"), // DD.MM.YYYY как в твоём CSV
			e.Name,
			fmt.Sprintf("%.2f", e.Amount),
			e.Category,
			fmt.Sprintf("%d", e.Priority),
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

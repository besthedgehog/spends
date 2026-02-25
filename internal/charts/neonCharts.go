package charts

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// const url = "http://localhost:5005/plot"

const url = "http://flask-charts:5005/plot"

// SendDataForPlot sends data for plotting to the server
// and returnes the image and the error
func SendDataForPlot(urlName string, body []byte) ([]byte, error) {
	path := fmt.Sprintf("%v/%v", url, urlName)

	resp, err := http.Post(
		path,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Проверим статус
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Читаем всё тело ответа
	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return imgBytes, nil
}

// Столбчатая по дням
func NeonByDayBars(body []byte) ([]byte, error) {
	name := "by-day"
	img, err := SendDataForPlot(name, body)
	return img, err
}

// Столбики по приоритетам (всего 3 столбика)
func NeonPriorityBars(body []byte) ([]byte, error) {
	name := "priority"
	img, err := SendDataForPlot(name, body)
	return img, err
}

// Столбики по категориям
func NeonCategoriesBars(body []byte) ([]byte, error) {
	name := "category"
	img, err := SendDataForPlot(name, body)
	return img, err
}

// График с накоплением
func NeonCumulative(body []byte) ([]byte, error) {
	name := "cumulative"
	img, err := SendDataForPlot(name, body)
	return img, err
}

// Точечная диаграмма
func NeonScatter(body []byte) ([]byte, error) {
	name := "scatter"
	img, err := SendDataForPlot(name, body)
	return img, err
}

// Пирог по категорям
func NeonPie(body []byte) ([]byte, error) {
	name := "pie"
	img, err := SendDataForPlot(name, body)
	return img, err
}

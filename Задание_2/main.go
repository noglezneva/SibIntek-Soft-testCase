package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*
Задание 2. Напишите программу.

Напишите консольную программу на Go, которая выполняет следующие действия:
1. Читает из файла JSON с массивом чисел.
2. Считает сумму всех чисел в массиве.
3. Выполняет HTTP GET запрос на заданный URL и проверяет статус ответа (должен быть 200).
4. Логирует результаты каждого шага в файл.
Важные детали:
- Если во время выполнения программы возникает ошибка, программа должна её корректно обработать и залогировать.
- Логи должны содержать динамичные данные там, где это возможно (url, статус, текст ошибки и т.д.).
- Формат и структура JSON-файла должны быть описаны в документации к программе.
- Программа должна быть оформлена с соблюдением принципов чистого кода и хороших практик Go.
*/

type Config struct {
	Numbers []int  `json:"numbers"`
	URL     string `json:"url"`
	LogFile string `json:"log_file"`
}

func main() {
	// Загрузка конфигурации
	config := loadConfig()

	// Вычисление суммы чисел
	sum := calculateSum(config.Numbers)
	log.Println("Sum of numbers:", sum)

	// Выполнение HTTP GET запроса по указанному URL
	resp, err := http.Get(config.URL)
	if err != nil {
		log.Fatalf("Error making HTTP GET request: %s", err)
	}
	defer resp.Body.Close()

	// Проверка статуса HTTP ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected HTTP response status: %d", resp.StatusCode)
	}

	log.Println("HTTP request successful")

	// Запись результатов в указанный файл для логирования
	file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %s", err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println("Sum of numbers:", sum)
	log.Printf("HTTP request to %s successful\n", config.URL)
}

// Функция loadConfig загружает конфигурацию из JSON файла
func loadConfig() Config {
	configFile := "config.json"

	// Чтение данных конфигурации из файла
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error decoding config file: %s", err)
	}

	return config
}

// Функция calculateSum вычисляет сумму чисел в массиве
func calculateSum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

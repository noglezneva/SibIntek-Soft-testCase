package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
Задание 3. Доработайте программу.
Модифицируйте программу из задания 2 таким образом, чтобы она могла:

1. Читать данные не только из файла, но и из стандартного ввода.
2. Принимать через аргументы командной строки параметр для определения источника данных (файл или stdin).
3. Сохранять результат работы в указанный пользователем файл, а не только в стандартный вывод.
4. Поддерживать конфигурацию через файл настроек или переменные окружения (например, URL для HTTP запроса).
5. Результаты заданий следует представить в текстовом виде (не в виде ссылок).
*/
type Config struct {
	Numbers    []int  `json:"numbers"`     // Срез для хранения чисел
	URL        string `json:"url"`         // Строка для хранения URL
	LogFile    string `json:"log_file"`    // Строка для хранения имени файла журнала
	OutputFile string `json:"output_file"` // Строка для хранения имени выходного файла
	DataSource string `json:"data_source"` // Строка для хранения источника данных
	InputFile  string `json:"input_file"`  // Строка для хранения имени входного файла
}

func main() {
	config := loadConfig() // Загрузка конфигурации из файла

	var inputData []byte
	switch config.DataSource {
	case "file":
		inputData = readFromFile(config.InputFile) // Чтение данных из файла
	case "stdin":
		inputData = readFromStdin() // Чтение данных из стандартного ввода
	default:
		log.Fatal("Invalid data source") // Если указан недопустимый источник данных, программа завершается с ошибкой
	}

	sum := calculateSum(config.Numbers) // Вычисление суммы чисел
	log.Println("Sum of numbers:", sum) // Вывод суммы чисел в логи

	resp, err := http.Get(config.URL) // Выполнение HTTP GET запроса
	if err != nil {
		log.Fatalf("Error making HTTP GET request: %s", err) // Если произошла ошибка, программа завершается с ошибкой
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected HTTP response status: %d", resp.StatusCode) // Если получен непредвиденный статус ответа HTTP, программа завершается с ошибкой
	}

	log.Println("HTTP request successful")

	file, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // Открытие файла журнала
	if err != nil {
		log.Fatalf("Error opening log file: %s", err)
	}
	defer file.Close()

	log.SetOutput(file)                                       // Настройка вывода журнала в файл
	log.Println("Sum of numbers:", sum)                       // Вывод суммы чисел в журнал
	log.Printf("HTTP request to %s successful\n", config.URL) // Вывод успешного выполнения запроса в журнал

	saveToFile(config.OutputFile, sum, string(inputData)) // Сохранение результата в файл
}

// Функция для загрузки конфигурации из файла
func loadConfig() Config {
	configFile := "config.json"              // Путь к файлу конфигурации
	data, err := ioutil.ReadFile(configFile) // Чтение данных из файла
	if err == nil {
		var config Config
		err = json.Unmarshal(data, &config) // Декодирование данных JSON в структуру Config
		if err == nil {
			return config
		}
	}

	// Если не удалось загрузить конфигурацию из файла, используем значения по умолчанию из переменных окружения
	config := Config{
		Numbers:    getNumbersFromEnv(),    // Получение чисел из переменной окружения NUMBERS
		URL:        getURLFromEnv(),        // Получение URL из переменной окружения URL
		LogFile:    getLogFileFromEnv(),    // Получение имени файла журнала из переменной окружения LOG_FILE
		OutputFile: getOutputFileFromEnv(), // Получение имени выходного файла из переменной окружения OUTPUT_FILE
		DataSource: getDataSourceFromEnv(), // Получение источника данных из переменной окружения DATA_SOURCE
		InputFile:  getInputFileFromEnv(),  // Получение имени входного файла из переменной окружения INPUT_FILE
	}
	return config
}

// Функция для получения чисел из переменной окружения NUMBERS
func getNumbersFromEnv() []int {
	numbersEnv := os.Getenv("NUMBERS")           // Получение значения переменной окружения NUMBERS
	numbersStr := strings.Split(numbersEnv, ",") // Разделение значения на подстроки по запятой
	var numbers []int
	for _, numStr := range numbersStr {
		num, err := strconv.Atoi(numStr) // Преобразование подстроки в число
		if err == nil {
			numbers = append(numbers, num)
		}
	}
	return numbers
}

// Функция для получения URL из переменной окружения URL
func getURLFromEnv() string {
	return os.Getenv("URL") // Получение значения переменной окружения URL
}

// Функция для получения имени файла журнала из переменной окружения LOG_FILE
func getLogFileFromEnv() string {
	return os.Getenv("LOG_FILE") // Получение значения переменной окружения LOG_FILE
}

// Функция для получения имени выходного файла из переменной окружения OUTPUT_FILE
func getOutputFileFromEnv() string {
	return os.Getenv("OUTPUT_FILE") // Получение значения переменной окружения OUTPUT_FILE
}

// Функция для получения источника данных из переменной окружения DATA_SOURCE
func getDataSourceFromEnv() string {
	return os.Getenv("DATA_SOURCE") // Получение значения переменной окружения DATA_SOURCE
}

// Функция для получения имени входного файла из переменной окружения INPUT_FILE
func getInputFileFromEnv() string {
	return os.Getenv("INPUT_FILE") // Получение значения переменной окружения INPUT_FILE
}

// Функция для чтения данных из файла
func readFromFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename) // Чтение данных из файла
	if err != nil {
		log.Fatalf("Error reading file: %s", err) // Если произошла ошибка при чтении файла, программа завершается с ошибкой
	}
	return data
}

// Функция для чтения данных из стандартного ввода
func readFromStdin() []byte {
	data, err := ioutil.ReadAll(os.Stdin) // Чтение данных из стандартного ввода
	if err != nil {
		log.Fatalf("Error reading from stdin: %s", err) // Если произошла ошибка при чтении из стандартного ввода, программа завершается с ошибкой
	}
	return data
}

// Функция для вычисления суммы чисел
func calculateSum(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// Функция для сохранения результата в файл
func saveToFile(filename string, sum int, inputData string) {
	output := fmt.Sprintf("Sum: %d\nInputData: %s\n", sum, inputData) // Форматирование вывода
	err := ioutil.WriteFile(filename, []byte(output), 0644)           // Запись данных в файл
	if err != nil {
		log.Fatalf("Error saving file: %s", err) // Если произошла ошибка при сохранении файла, программа завершается с ошибкой
	}
}

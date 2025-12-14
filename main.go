package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <logfile.csv>")
	}
	
	filename := os.Args[1]
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Чтение логов из файла
	inputChan, err := readLogs(filename)
	if err != nil {
		log.Fatalf("Failed to read logs: %v", err)
	}
	
	// Параллельная обработка с использованием worker pool
	numWorkers := 3
	processedChan := processLogs(ctx, inputChan, numWorkers)
	
	// Фильтрация только ошибок (4xx и 5xx)
	filteredChan := filterLogs(processedChan, 400)
	
	// Подсчет статистики
	stats := calculateStats(filteredChan)
	
	// Вывод результатов
	printStatistics(stats)
}

func printStatistics(stats Statistics) {
	fmt.Println("\n=== Анализ логов веб-сервера ===")
	fmt.Printf("Общее количество запросов: %d\n", stats.TotalRequests)
	fmt.Printf("Количество ошибок (4xx/5xx): %d\n", stats.ErrorCount)
	fmt.Printf("Среднее время ответа: %.2f мс\n", stats.AverageRespTime)
	fmt.Printf("Количество уникальных IP-адресов: %d\n", len(stats.RequestsByIP))
	
	if len(stats.RequestsByIP) > 0 {
		fmt.Println("\nТоп 5 IP-адресов по количеству запросов:")
		printTopIPs(stats.RequestsByIP, 5)
	}
	fmt.Println("==================================")
}
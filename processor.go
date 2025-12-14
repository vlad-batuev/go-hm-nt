package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type LogEntry struct {
	Timestamp    string // время в формате "2024-01-15 10:30:00"
	IP           string // IP адрес клиента
	Method       string // HTTP метод (GET, POST и т.д.)
	URL          string // путь запроса
	StatusCode   int    // HTTP статус код
	ResponseTime int    // время ответа в миллисекундах
}

type Statistics struct {
	TotalRequests   int            // общее количество запросов
	ErrorCount      int            // количество ошибок (статус >= 400)
	RequestsByIP    map[string]int // количество запросов с каждого IP
	AverageRespTime float64        // среднее время ответа
}

func readLogs(filename string) (<-chan LogEntry, error) {
	out := make(chan LogEntry)
	
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	
	go func() {
		defer func() {
			file.Close()
			close(out)
			log.Println("Чтение логов завершено")
		}()
		
		scanner := bufio.NewScanner(file)
		lineNumber := 0
		
		for scanner.Scan() {
			line := scanner.Text()
			lineNumber++
			
			// Пропускаем заголовок
			if lineNumber == 1 && strings.Contains(line, "timestamp") {
				continue
			}
			
			entry, err := parseLogLine(line)
			if err != nil {
				log.Printf("Ошибка парсинга строки %d: %v", lineNumber, err)
				continue
			}
			
			out <- entry
		}
		
		if err := scanner.Err(); err != nil {
			log.Printf("Ошибка при чтении файла: %v", err)
		}
	}()
	
	return out, nil
}

func parseLogLine(line string) (LogEntry, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 6 {
		return LogEntry{}, fmt.Errorf("неверное количество полей: %d", len(parts))
	}
	
	statusCode, err := strconv.Atoi(strings.TrimSpace(parts[4]))
	if err != nil {
		return LogEntry{}, fmt.Errorf("ошибка парсинга статус-кода: %v", err)
	}
	
	responseTime, err := strconv.Atoi(strings.TrimSpace(parts[5]))
	if err != nil {
		return LogEntry{}, fmt.Errorf("ошибка парсинга времени ответа: %v", err)
	}
	
	return LogEntry{
		Timestamp:    strings.TrimSpace(parts[0]),
		IP:           strings.TrimSpace(parts[1]),
		Method:       strings.TrimSpace(parts[2]),
		URL:          strings.TrimSpace(parts[3]),
		StatusCode:   statusCode,
		ResponseTime: responseTime,
	}, nil
}

func processLogs(ctx context.Context, input <-chan LogEntry, numWorkers int) <-chan LogEntry {
	out := make(chan LogEntry)
	var wg sync.WaitGroup
	
	// Запускаем воркеров
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for entry := range input {
				select {
				case <-ctx.Done():
					log.Printf("Воркер %d остановлен по контексту", workerID)
					return
				default:
					// Валидация и обогащение данных
					if entry.StatusCode < 100 || entry.StatusCode > 599 {
						log.Printf("Воркер %d: некорректный статус-код %d", workerID, entry.StatusCode)
						continue
					}
					
					if entry.ResponseTime < 0 {
						log.Printf("Воркер %d: некорректное время ответа %d", workerID, entry.ResponseTime)
						entry.ResponseTime = 0
					}
					
					out <- entry
				}
			}
		}(i)
	}
	
	// Закрываем выходной канал после завершения всех воркеров
	go func() {
		wg.Wait()
		close(out)
		log.Println("Обработка логов завершена")
	}()
	
	return out
}

func filterLogs(input <-chan LogEntry, minStatus int) <-chan LogEntry {
	out := make(chan LogEntry)
	
	go func() {
		defer close(out)
		
		for entry := range input {
			if entry.StatusCode >= minStatus {
				out <- entry
			}
		}
		
		log.Println("Фильтрация логов завершена")
	}()
	
	return out
}

func calculateStats(input <-chan LogEntry) Statistics {
	stats := Statistics{
		RequestsByIP: make(map[string]int),
	}
	
	var totalResponseTime int
	var count int
	
	for entry := range input {
		stats.TotalRequests++
		totalResponseTime += entry.ResponseTime
		count++
		
		if entry.StatusCode >= 400 {
			stats.ErrorCount++
		}
		
		stats.RequestsByIP[entry.IP]++
	}
	
	if count > 0 {
		stats.AverageRespTime = float64(totalResponseTime) / float64(count)
	}
	
	return stats
}

func printTopIPs(requestsByIP map[string]int, n int) {
	type ipCount struct {
		ip    string
		count int
	}
	
	var ipCounts []ipCount
	for ip, count := range requestsByIP {
		ipCounts = append(ipCounts, ipCount{ip, count})
	}
	
	// Сортировка по убыванию количества запросов
	sort.Slice(ipCounts, func(i, j int) bool {
		return ipCounts[i].count > ipCounts[j].count
	})
	
	// Вывод топ N IP-адресов
	for i := 0; i < len(ipCounts) && i < n; i++ {
		fmt.Printf("%d. %s - %d запросов\n", i+1, ipCounts[i].ip, ipCounts[i].count)
	}
}
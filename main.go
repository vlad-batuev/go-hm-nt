package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// Job представляет задание для обработки
type Job struct {
	ID  int
	URL string
}

// Result представляет результат обработки задания
type Result struct {
	Job      Job
	Status   string
	Duration time.Duration
}

// generateJobs создает задания на основе списка URL
func generateJobs(urls []string) []Job {
	jobs := make([]Job, len(urls))
	for i, url := range urls {
		jobs[i] = Job{
			ID:  i + 1,
			URL: url,
		}
	}
	return jobs
}

// worker обрабатывает задания из канала jobs и отправляет результаты в канал results
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d started processing job %d: %s\n", id, job.ID, job.URL)

		// Имитация HTTP-запроса
		start := time.Now()
		processingTime := time.Millisecond * time.Duration(rand.Intn(1000)+500) // Случайная задержка 500-1500 мс
		time.Sleep(processingTime)

		// Формирование результата
		result := Result{
			Job:      job,
			Status:   "processed",
			Duration: time.Since(start),
		}

		// Отправка результата
		results <- result
		fmt.Printf("Worker %d finished processing job %d: %s (took %v)\n",
			id, job.ID, job.URL, result.Duration)
	}
}

// collectResults собирает результаты из канала в слайс
func collectResults(results <-chan Result, done chan<- bool) []Result {
	var allResults []Result

	for result := range results {
		allResults = append(allResults, result)
	}

	done <- true
	return allResults
}

func main() {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Список URL для обработки
	urls := []string{
		"https://example.com",
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
		"https://golang.org",
		"https://medium.com",
		"https://reddit.com",
		"https://twitter.com",
		"https://youtube.com",
		"https://linkedin.com",
	}

	// Создание заданий
	jobs := generateJobs(urls)

	// Параметры worker pool
	numWorkers := 3
	numJobs := len(jobs)

	// Создание каналов
	jobsChan := make(chan Job, numJobs)
	resultsChan := make(chan Result, numJobs)
	done := make(chan bool)

	// WaitGroup для синхронизации воркеров
	var wg sync.WaitGroup

	// Fan-out: запуск воркеров
	fmt.Printf("Starting %d workers to process %d jobs...\n\n", numWorkers, numJobs)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobsChan, resultsChan, &wg)
	}

	// Отправка заданий в канал
	go func() {
		for _, job := range jobs {
			jobsChan <- job
		}
		close(jobsChan)
	}()

	// Fan-in: сбор результатов
	var allResults []Result
	go func() {
		allResults = collectResults(resultsChan, done)
	}()

	// Ожидание завершения всех воркеров и закрытие канала результатов
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Ожидание завершения сбора результатов
	<-done

	// Вывод отчета
	printReport(allResults)
}

// printReport выводит агрегированный отчет по результатам
func printReport(results []Result) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("FINAL REPORT")
	fmt.Println(strings.Repeat("=", 60))

	var totalTime time.Duration
	var fastest, slowest time.Duration = time.Hour, 0
	var fastestURL, slowestURL string

	fmt.Println("\nDetailed Results:")
	fmt.Println(strings.Repeat("-", 60))

	for _, result := range results {
		fmt.Printf("Job %2d: %-30s -> %-10s (Time: %v)\n",
			result.Job.ID, result.Job.URL, result.Status, result.Duration)

		totalTime += result.Duration

		// Поиск самого быстрого
		if result.Duration < fastest {
			fastest = result.Duration
			fastestURL = result.Job.URL
		}

		// Поиск самого медленного
		if result.Duration > slowest {
			slowest = result.Duration
			slowestURL = result.Job.URL
		}
	}

	// Вывод статистики
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("Statistics:")
	fmt.Println(strings.Repeat("-", 60))

	avgTime := totalTime / time.Duration(len(results))

	fmt.Printf("Total jobs processed: %d\n", len(results))
	fmt.Printf("Total processing time: %v\n", totalTime)
	fmt.Printf("Average processing time: %v\n", avgTime)
	fmt.Printf("Fastest: %v (%s)\n", fastest, fastestURL)
	fmt.Printf("Slowest: %v (%s)\n", slowest, slowestURL)

	// Гистограмма времени выполнения
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("Time Distribution:")
	fmt.Println(strings.Repeat("-", 60))

	// Группировка по диапазонам времени
	fastCount, mediumCount, slowCount := 0, 0, 0
	for _, result := range results {
		switch {
		case result.Duration < 800*time.Millisecond:
			fastCount++
		case result.Duration < 1200*time.Millisecond:
			mediumCount++
		default:
			slowCount++
		}
	}

	fmt.Printf("Fast (<800ms):   %d jobs\n", fastCount)
	fmt.Printf("Medium (800-1200ms): %d jobs\n", mediumCount)
	fmt.Printf("Slow (>1200ms):  %d jobs\n", slowCount)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("All jobs completed successfully!")
	fmt.Println(strings.Repeat("=", 60))
}

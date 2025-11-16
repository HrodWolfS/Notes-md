package main

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

type SearchResult struct {
	FilePath string
	LineNum  int
	Line     string
	Match    string
}

type searchStartedMsg struct{}

type searchCompletedMsg struct {
	results []SearchResult
	err     error
}

func performContentSearch(rootDir, query string) tea.Cmd {
	return func() tea.Msg {
		results := searchFiles(rootDir, query)
		return searchCompletedMsg{results: results}
	}
}

func searchFiles(rootDir, query string) []SearchResult {
	var results []SearchResult
	var resultsMutex sync.Mutex

	jobs := make(chan string, 100)
	resultsChan := make(chan []SearchResult, 100)

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go searchWorker(jobs, resultsChan, query, &wg)
	}

	go func() {
		filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if !info.IsDir() && filepath.Ext(path) == ".md" {
				jobs <- path
			}
			return nil
		})
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for res := range resultsChan {
		resultsMutex.Lock()
		results = append(results, res...)
		resultsMutex.Unlock()
	}

	return results
}

func searchWorker(jobs <-chan string, results chan<- []SearchResult, query string, wg *sync.WaitGroup) {
	defer wg.Done()

	queryLower := strings.ToLower(query)

	for path := range jobs {
		matches := searchInFile(path, queryLower)
		if len(matches) > 0 {
			results <- matches
		}
	}
}

func searchInFile(path, queryLower string) []SearchResult {
	var results []SearchResult

	file, err := os.Open(path)
	if err != nil {
		return results
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		lineLower := strings.ToLower(line)

		if strings.Contains(lineLower, queryLower) {
			results = append(results, SearchResult{
				FilePath: path,
				LineNum:  lineNum,
				Line:     strings.TrimSpace(line),
				Match:    queryLower,
			})

			if len(results) >= 5 {
				break
			}
		}
	}

	return results
}

type searchResultItem struct {
	result SearchResult
}

func (s searchResultItem) Title() string {
	fileName := filepath.Base(s.result.FilePath)
	return "📝 " + fileName + ":" + string(rune(s.result.LineNum))
}

func (s searchResultItem) Description() string {
	if len(s.result.Line) > 80 {
		return s.result.Line[:80] + "..."
	}
	return s.result.Line
}

func (s searchResultItem) FilterValue() string {
	return s.result.Line
}

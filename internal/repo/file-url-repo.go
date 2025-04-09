package repo

import (
	"bufio"
	"fmt"
	"os"
)

type FileUrlRepo struct {
	filename string
}

func NewFileUrlRepo(filename string) *FileUrlRepo {
	return &FileUrlRepo{filename: filename}
}

func (r *FileUrlRepo) GetUrls() ([]string, error) {
	file, err := os.Open(r.filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла со ссылками: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла со ссылками: %w", err)
	}

	return lines, nil
}

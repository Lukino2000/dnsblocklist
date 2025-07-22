package main

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

// Legge gli URL da un file
func ReadUrlsFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}
	return urls, scanner.Err()
}

// Scarica tutti i file e restituisce una slice di stringhe (contenuto dei file)
func DownloadFiles(urls []string) ([]string, error) {
	var results []string
	for _, url := range urls {
		content, err := DownloadFile(url)
		if err != nil {
			return nil, err
		}
		results = append(results, content)
	}
	return results, nil
}

// Scarica un singolo file da URL
func DownloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Errore HTTP: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Legge la whitelist da file (uno per riga)
func ReadWhitelistFile(filename string) (map[string]struct{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	whitelist := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			whitelist[line] = struct{}{}
		}
	}
	return whitelist, scanner.Err()
}

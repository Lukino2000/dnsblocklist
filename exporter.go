package main

import (
	"fmt"
	"os"
	"sort"
	"time"
)

// Esporta la mappa domini in formato hosts, includendo header con data e URL
func ExportHostsFile(filename string, domains map[string]struct{}, urls []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Data e ora attuali
	now := time.Now().Format("2006-01-02 15:04:05 MST")

	// Scrivi header
	file.WriteString("# Hosts file generato automaticamente\n")
	file.WriteString(fmt.Sprintf("# Data di generazione: %s\n", now))
	file.WriteString("# Lista degli URL sorgente:\n")
	for _, url := range urls {
		file.WriteString(fmt.Sprintf("#   %s\n", url))
	}

	// Ordina per dominio
	domainList := make([]string, 0, len(domains))
	for d := range domains {
		domainList = append(domainList, d)
	}
	sort.Strings(domainList)

	for _, domain := range domainList {
		// Scrivi in formato hosts
		_, err := fmt.Fprintf(file, "0.0.0.0 %s\n", domain)
		if err != nil {
			return err
		}
	}
	return nil
}

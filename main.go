package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Avvio programma...")

	// Leggi gli URL dal file
	urls, err := ReadUrlsFile("urls.txt")
	if err != nil {
		log.Fatalf("Errore lettura urls.txt: %v", err)
	}

	log.Printf("Trovati %d URL, scarico i file...", len(urls))

	// Scarica i file in memoria
	files, err := DownloadFiles(urls)
	if err != nil {
		log.Fatalf("Errore nel download: %v", err)
	}

	log.Printf("Download completato. Inizio parsing dei file...")

	// Parsing e deduplica
	domains := make(map[string]struct{})
	for i, content := range files {
		parsed := ParseBlocklist(content, urls[i])
		for _, domain := range parsed {
			domains[domain] = struct{}{}
		}
		log.Printf("File %d: %d domini trovati - %s", i+1, len(parsed), urls[i])
	}

	log.Printf("Totale domini unici prima della whitelist: %d", len(domains))

	// Carica la whitelist se esiste
	whitelist, err := ReadWhitelistFile("whitelist.txt")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Errore lettura whitelist.txt: %v", err)
	}
	if len(whitelist) > 0 {
		log.Printf("Whitelist caricata: %d domini", len(whitelist))
		// Rimuovi i domini presenti in whitelist
		for d := range whitelist {
			delete(domains, d)
		}
		log.Printf("Totale domini dopo whitelist: %d", len(domains))
	} else {
		log.Println("Nessuna whitelist trovata o vuota.")
	}

	// Esporta in formato hosts, includendo data/ora e lista URL
	err = ExportHostsFile("hosts", domains, urls)
	if err != nil {
		log.Fatalf("Errore esportazione hosts: %v", err)
	}

	log.Println("File hosts generato con successo!")
}

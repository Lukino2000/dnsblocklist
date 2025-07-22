package main

import (
	"bufio"
	"regexp"
	"strings"
)

// Rileva il formato del file e chiama il parser corretto
func ParseBlocklist(content string, source string) []string {
	// Heuristic: se contiene "[Adblock Plus]", è ABP
	if strings.Contains(content, "[Adblock Plus]") || strings.Contains(content, "||") {
		return ParseABP(content)
	}
	// Se contiene "0.0.0.0" o "127.0.0.1", è hosts
	if strings.Contains(content, "0.0.0.0") || strings.Contains(content, "127.0.0.1") {
		return ParseHosts(content)
	}
	// Altrimenti, fallback: cerca righe tipo "||dominio^"
	return ParseABP(content)
}

// Parsing formato hosts (StevenBlack)
func ParseHosts(content string) []string {
	var domains []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		ip := fields[0]
		// Solo righe che iniziano con 0.0.0.0 o 127.0.0.1
		if ip == "0.0.0.0" || ip == "127.0.0.1" {
			domain := fields[1]
			if IsValidDomain(domain) {
				domains = append(domains, domain)
			}
		}
	}
	return domains
}

// Parsing formato ABP / AdGuard / filter
func ParseABP(content string) []string {
	var domains []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	re := regexp.MustCompile(`\|\|([a-zA-Z0-9\.\-\_]+)\^`)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "!") || strings.HasPrefix(line, "[") {
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) == 2 {
			domain := matches[1]
			if IsValidDomain(domain) {
				domains = append(domains, domain)
			}
		}
	}
	return domains
}

// Verifica se una stringa è un dominio valido
func IsValidDomain(domain string) bool {
	// Semplice check: almeno un punto e nessuno spazio
	return strings.Contains(domain, ".") && !strings.Contains(domain, " ")
}

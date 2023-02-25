package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	// wordlist dosyasını aç
	file, err := os.Open("wordlist.txt")
	if err != nil {
		fmt.Println("Dosya açılırken hata oluştu:", err)
		return
	}
	defer file.Close()

	// dosyayı satır satır oku
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		packageName := strings.TrimSpace(scanner.Text())
		if packageName == "" {
			continue
		}
		fmt.Printf("paket ismi >> %s ", packageName)
		resp, err := http.Get("https://registry.npmjs.org/" + packageName)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if resp.StatusCode == 200 {
			fmt.Println("Bu paket npmjs'de tanımlıdır")
		} else {
			fmt.Println("Bu paket npmjs'de tanımlı değildir")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Dosya okunurken hata oluştu:", err)
	}
}

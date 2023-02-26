package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	// Command line flag'ları oluştur
	tokenPtr := flag.String("t", "", "GitHub token")
	orgPtr := flag.String("org", "", "GitHub organization")

	flag.Parse()

	// Zorunlu flag'ların kontrolü
	if *tokenPtr == "" || *orgPtr == "" {
		flag.Usage()
		os.Exit(1)
	}

	// GitHub istemcisine bağlanmak için bir HTTP istemci oluşturulur
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *tokenPtr},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// GitHub organization'ın repoları için arama sorgusu oluşturulur
	query := fmt.Sprintf("org:%s", *orgPtr)
	opts := &github.SearchOptions{
		Sort:        "updated",
		Order:       "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// İlk sayfadaki repoların sayısını al
	result, _, err := client.Search.Repositories(ctx, query, opts)
	if err != nil {
		log.Fatal(err)
	}
	totalCount := *result.Total
	fmt.Printf("Found %d repositories at %s\n", totalCount, *orgPtr)

	// İşlem başlamadan önce 2 saniye bekle
	time.Sleep(2 * time.Second)

	for {
		// GitHub arama sorgusu çalıştırılır
		result, resp, err := client.Search.Repositories(ctx, query, opts)
		if err != nil {
			log.Fatal(err)
		}

		// Bulunan sonuçlar işlenir
		for _, repo := range result.Repositories {
			cloneURL := *repo.CloneURL
			dirName := strings.Replace(*repo.FullName, "/", "_", -1)
			dirPath := filepath.Join(".", *orgPtr, dirName)

			// Repoyu indirme işlemi yapılır
			cmd := exec.Command("git", "clone", cloneURL, dirPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Run()
			if err != nil {
				log.Printf("Error cloning repository %s: %s", cloneURL, err)
			} else {
				fmt.Printf("Repository cloned: %s\n", cloneURL)
			}
		}

		// Sonraki sayfa var mı diye kontrol edilir
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
}

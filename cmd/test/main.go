// this file is for testing purposes only
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"velo-server/internal/config"
	"velo-server/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while loading the .env or it is missing")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r2 := storage.NewR2Storage(cfg)

	// Replace with an actual object key in your bucket.
	objectKey := "Anyma, Volkoder - Other Dimension.mp3"

	url, err := r2.GeneratePresignedURL(
		context.Background(),
		objectKey,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Streaming URL:", url)

	// Optional: save URL for testing.
	_ = os.WriteFile("stream-url.txt", []byte(url), 0600)
}

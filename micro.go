package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"microservices/modules" // modules paketini içe aktar
)

func main() {
	// PostgreSQL bağlantısını ayarlıyoruz
	connStr := "DB CONNECTION"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Veritabanı bağlantı hatası: %v", err)
	}
	defer db.Close()

	// Bağlantıyı doğrula
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Bağlantı doğrulama hatası: %v", err)
	}

	fmt.Println("Bağlantı başarılı!")

	// Ticker'ları ayarlıyoruz
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Null kontrol işlemini ayrı bir goroutine içinde başlatıyoruz
	go func() {
		for {
			fmt.Println("Null Token Kontrol Ediliyor...")
			modules.NullControl(db) // modules paketinden NullControl fonksiyonunu çağırıyoruz
			<-ticker.C
		}
	}()
	//
	//// Auto price işlemini ayrı bir goroutine içinde başlatıyoruz
	go func() {
		ticker2 := time.NewTicker(2 * time.Minute)
		defer ticker2.Stop()
		for {
			fmt.Println("Token Fiyat Kontrol Ediliyor...")
			modules.AutoPrice(db) // modules paketinden AutoPrice fonksiyonunu çağırıyoruz
			<-ticker2.C
		}
	}()
	go func() {
		ticker2 := time.NewTicker(2 * time.Minute)
		defer ticker2.Stop()
		for {
			fmt.Println("Cüzdan Usd Değeri Kontrol Ediliyor...")
			modules.UpdateWallets(db) // modules paketinden AutoPrice fonksiyonunu çağırıyoruz
			<-ticker2.C
		}
	}()

	select {} // Sonsuz bekleme
}

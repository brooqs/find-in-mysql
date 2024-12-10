package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

func main() {
	// Parametreleri tanımlayın
	tableName := flag.String("t", "", "Arama yapılacak tablo adı")
	searchWord := flag.String("w", "", "Aranacak kelime")
	flag.Parse()

	// Parametre kontrolü
	if *tableName == "" || *searchWord == "" {
		log.Fatalf("Hata: -t ve -w parametreleri zorunludur. Örnek: -t tabloAdi -w arananKelime")
	}

	// Config dosyasını yükleyin veya oluşturun
	configFile := "config.ini"
	var cfg *ini.File
	var err error

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("config.ini dosyası bulunamadı. Lütfen veritabanı bilgilerini girin:")
		cfg = ini.Empty()

		host := askInput("Host (default: 127.0.0.1): ", "127.0.0.1")
		port := askInput("Port (default: 3306): ", "3306")
		user := askInput("Kullanıcı adı: ", "")
		password := askInput("Şifre: ", "")
		database := askInput("Veritabanı adı: ", "")

		// Config dosyasına yaz
		cfg.Section("mysql").Key("host").SetValue(host)
		cfg.Section("mysql").Key("port").SetValue(port)
		cfg.Section("mysql").Key("user").SetValue(user)
		cfg.Section("mysql").Key("password").SetValue(password)
		cfg.Section("mysql").Key("database").SetValue(database)

		if err = cfg.SaveTo(configFile); err != nil {
			log.Fatalf("Config dosyası kaydedilirken hata oluştu: %v", err)
		}
		fmt.Println("Config bilgileri config.ini dosyasına kaydedildi.")
	} else {
		cfg, err = ini.Load(configFile)
		if err != nil {
			log.Fatalf("Config dosyası yüklenirken hata oluştu: %v", err)
		}
	}

	// Veritabanı bilgilerini oku
	host := cfg.Section("mysql").Key("host").String()
	port := cfg.Section("mysql").Key("port").String()
	user := cfg.Section("mysql").Key("user").String()
	password := cfg.Section("mysql").Key("password").String()
	database := cfg.Section("mysql").Key("database").String()

	// MySQL bağlantı dizesi oluştur
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

	// Veritabanına bağlan
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Veritabanına bağlanırken hata: %v", err)
	}
	defer db.Close()

	// Arama yap
	fmt.Printf("Tablo: %s, Aranan kelime: %s\n", *tableName, *searchWord)
	query := fmt.Sprintf("SELECT * FROM %s WHERE CONCAT_WS(' ', %s) LIKE ?", *tableName, getAllColumnNames(db, *tableName))
	rows, err := db.Query(query, "%"+*searchWord+"%")
	if err != nil {
		log.Fatalf("Arama sırasında hata oluştu: %v", err)
	}
	defer rows.Close()

	// Sonuçları göster
	columns, _ := rows.Columns()
	fmt.Printf("Sonuçlar (%d sütun):\n", len(columns))
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		for i, col := range columns {
			fmt.Printf("%s: %v\t", col, values[i])
		}
		fmt.Println()
	}
}

// Kullanıcıdan giriş alır ve varsayılan değer önerir
func askInput(prompt string, defaultValue string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	if input == "" {
		return defaultValue
	}
	return input
}

// Tablo içindeki tüm sütunları listele
func getAllColumnNames(db *sql.DB, tableName string) string {
	query := fmt.Sprintf("SHOW COLUMNS FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Sütunlar alınırken hata oluştu (Tablo: %s): %v", tableName, err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column, dataType, isNull, key, defaultValue, extra string
		if err := rows.Scan(&column, &dataType, &isNull, &key, &defaultValue, &extra); err != nil {
			log.Fatalf("Sütun okuma hatası: %v", err)
		}
		columns = append(columns, column)
	}
	return strings.Join(columns, ", ")
}

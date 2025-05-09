package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Greeter interface {
	Greet() string
}

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() string {
	return fmt.Sprintf("Hello, my name is %s and I am %d years old.", p.Name, p.Age)
}

func printValue[T any](value T) {
	fmt.Printf("Value: %v\n", value)
}

func printGoVersion() {
	fmt.Println("Go version:", runtime.Version())
}

func checkAge(age int) string {
	if age < 18 {
		return "Underage"
	} else if age >= 18 && age <= 65 {
		return "Adult"
	}
	return "Senior"
}

func sayGreeting(g Greeter) string {
	return g.Greet()
}

func main() {
	dsn := "root:password@tcp(db:3306)/testdb"
	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil && db.Ping() == nil {
			break
		}
		log.Println("Waiting for MySQL...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to MySQL: %v", err)
	}
	defer db.Close()

	fmt.Println("MySQL connection established!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server is running!\n")

		printGoVersion()
		printValue(31)
		printValue("Hello!")

		p := Person{Name: "Son", Age: 31}
		fmt.Fprintf(w, sayGreeting(p)+"\n")
		fmt.Fprintf(w, checkAge(p.Age)+"\n")
	})

	fmt.Println("HTTP server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

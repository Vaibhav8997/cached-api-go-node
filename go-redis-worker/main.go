package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type UserData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	//postgres connection setup
	dbAddr := "host=postgres port=5432 user=admin password=123456 dbname=testDb sslmode=disable"

	db, err := sql.Open("postgres", dbAddr)
	if err != nil {
		log.Fatal("DB Connection Error:", err)
	}
	defer db.Close()

	//Redis connection setup
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	fmt.Println("Redis is running.....waiting for data!")

	// infinite loop
	for {
		//BRPop -> block and pop    ->Blocking right pop
		result, err := rdb.BRPop(ctx, 0, "user_queue").Result()
		if err != nil {
			fmt.Println("Error extracting data from redis...", err)
			continue
		}
		//convert json string into struct
		var user UserData
		json.Unmarshal([]byte(result[1]), &user)
		// if err != nil {
		// 	fmt.Println("Incorrect JSON: ", err)
		// 	continue
		// }

		//insert data into postgres
		query := "INSERT INTO test_data(id, name, email) VALUES($1, $2, $3)"

		_, err = db.Exec(query, user.ID, user.Name, user.Email)

		if err != nil {
			fmt.Printf("Error saving id %d: %v\n", user.ID, err)
		} else {
			fmt.Printf("Data saved successfully: %s (ID: %d)\n", user.Name, user.ID)
		}
	}
}

package main

import (
	"encoding/json"
	"fmt"

	"github.com/SowinskiBraeden/local-go-db/driver"
)

type User struct {
	Name string
	Age  json.Number
}

func main() {
	database, err := driver.NewConnection("mytestdb", nil)
	// database, err := driver.Connect("mytestdb", nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// users := []User{
	// 	{"John Doe", "25"},
	// 	{"Mike Tyson", "50"},
	// }

	// for _, value := range users {
	// 	database.InsertOne("users", User{
	// 		Name: value.Name,
	// 		Age:  value.Age,
	// 	})
	// }

	// err = database.InsertOne("users", User{
	// 	Name: "Test User",
	// 	Age:  "19",
	// })
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }

	record, findErr := database.FindOne("users", "54d65722-1a2e-455c-a75d-38f9026e8f97") // <-- Pass in the ID of the object generated
	if findErr != nil {
		fmt.Println("Error: ", findErr)
	}
	fmt.Println(record)

	// records, err := database.ReadAll("users")
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }
	// fmt.Println(records)

	// allUsers := []User{}

	// for _, f := range records {
	// 	userFound := User{}
	// 	if err := json.Unmarshal([]byte(f), &userFound); err != nil {

	// 	}
	// }

	// delErr := database.Delete("users", "44516cb3-eacd-49d6-9816-30f317e2dd5d")
	// if delErr != nil {
	// 	fmt.Println("Error: ", err)
	// }
}

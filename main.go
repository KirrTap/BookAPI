package main

import (
	"fmt"
	"flag"
)

const BookAPI = "https://openlibrary.org/works/%s.json"


func main() {
	bookKey := flag.String("key", "", "Book key")     //OL27448W 
	sortOrder := flag.String("sort", "", "Sort order (asc or desc)")
	

	flag.Parse()
	if *bookKey == "" {
		fmt.Println("Error: Specify book key with -key.")
		return
	}

	if *sortOrder != "" && *sortOrder != "asc" && *sortOrder != "desc" {
		fmt.Println("Error: You need to specify a valid sort order -sort (asc or desc).")
		return
	}

	bookAPI := fmt.Sprintf(BookAPI, *bookKey)

	fmt.Println(*bookKey)
	fmt.Println(*sortOrder)
	fmt.Println(bookAPI)
}


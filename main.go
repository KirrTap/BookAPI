package main

import (
	"fmt"
	"flag"
	"net/http"
	"io/ioutil"
	"encoding/json"
	// "gopkg.in/yaml.v2"
)

type AuthorKey struct {
	Authors []struct {
		Author struct {
			Key string `json:"key"`
		} `json:"author"`
	} `json:"authors"`
}

const BookAPI = "https://openlibrary.org/works/%s.json"


func main() {
	bookKey := flag.String("key", "", "Book key")     //OL27448W   OL18146933W
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

	
	author_key := getAuthorKey(*bookKey)


	for _,author := range author_key.Authors {
		fmt.Printf(author.Author.Key + "\n")
	}



	// bookAPI := fmt.Sprintf(BookAPI, *bookKey)

	// fmt.Println(*bookKey)
	// fmt.Println(bookAPI)
}


func getAuthorKey(bookKey string) AuthorKey {
	resp, err := http.Get(fmt.Sprintf(BookAPI, bookKey))
	if err != nil {
		fmt.Println("No response from request")
		return AuthorKey{}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result AuthorKey
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error unmarshalling JSON: ", err)
		return AuthorKey{}
	}

	return result


}




// for _,author := range author_key.Authors {
// 	yamlData, err := yaml.Marshal(author.Author)
// 	if err != nil {
// 		fmt.Println("Error converting to YAML:", err)
// 		return
// 	}
// 	fmt.Println(string(yamlData))
// }


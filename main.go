package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
	"sort"
)

const BookAPI = "https://openlibrary.org/search.json?q=%s"

type AuthorBook struct {
	Docs []struct {
		AuthorName 					     []string `json:"author_name"`
		Title                            string   `json:"title"`
		PublishDate                      []string `json:"publish_date,omitempty"`
		EditionCount                     int      `json:"edition_count"`
		Isbn                             []string `json:"isbn,omitempty"`
		
	} `json:"docs"`
}


func main() {
	bookTitle := "1984"
	sortOrder := "desc"
	
	resp, err := http.Get(fmt.Sprintf(BookAPI, strings.Replace(bookTitle, " ", "+", -1)))

	if err != nil {
		fmt.Println("No response from request")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result AuthorBook
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}


	AuthorNameSort(sortOrder, result)
	PublishDateSort(sortOrder, result)

	yamlData, err := yaml.Marshal(result)
	if err != nil {
		fmt.Println("Error converting to YAML:", err)
		return
	}

	fmt.Println(string(yamlData))

}


func AuthorNameSort(sortOrder string, result AuthorBook) {
	for _, sorted := range result.Docs {
		if sortOrder == "asc" {
			if len(sorted.AuthorName) > 0 {
				sort.Strings(sorted.AuthorName)
			}
		} else if sortOrder == "desc" {
			if len(sorted.AuthorName) > 0 {
				sort.Sort(sort.Reverse(sort.StringSlice(sorted.AuthorName)))
			}
		}
	}
	
}

func PublishDateSort(sortOrder string, result AuthorBook) {
	for _, sorted := range result.Docs {
		if sortOrder == "asc" {
			if len(sorted.PublishDate) > 0 {
				sort.Strings(sorted.PublishDate)
			}
		} else if sortOrder == "desc" {
			if len(sorted.PublishDate) > 0 {
				sort.Sort(sort.Reverse(sort.StringSlice(sorted.PublishDate)))
			}
		}
	}
}
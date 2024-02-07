package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

//APIs
const (
	findAuthorKeyPI = "https://openlibrary.org/works/%s.json" //APi na najdenie kluca autora
	findAuthorNameAPI = "https://openlibrary.org/authors/%s.json" //API na najdenie mena autora
	findAuthorBooksInfoAPI = "https://openlibrary.org/search.json?author_key=%s&fields=author_name,title,first_publish_year,isbn,edition_count" //API na najdenie informacii o knihach autora
		
)

// struktura pre ziskanie kluca autora
type AuthorKey struct {
	Authors []struct {
		Author struct {
			Key string `json:"key"`
		} `json:"author"`
	} `json:"authors"`
}

// struct pre ziskanie mena autora
type AuthorName struct {
	Name string `json:"name"`
}

// struktura pre ziskanie informacii o knihach autora
type BookInfo struct {
	Title           string   `json:"title"`
	FirstPublishYear int      `json:"first_publish_year"`
	Isbn             []string `json:"isbn"`
	EditionCount     int      `json:"edition_count"`
}

// moja struktura na vypis informacii
type MyOutput struct {
	Author struct{
		AuthorName string `json:"author_name"`
	} `json:"author"`
	Books struct {
		Title           string   `json:"title"`
		FirstPublishYear int      `json:"first_publish_year"`
		Isbn             []string `json:"isbn"`
		EditionCount     int      `json:"edition_count"`
	} `json:"books"`
}


func main(){
	//nacitanie parametrov
	bookKey := flag.String("key", "", "Book key")  // napr OL27448W  OL18146933W OL27370133W
	sortOrder := flag.String("sort", "", "Sort order (asc or desc)")
	flag.Parse()

	//kontrola ci je zadany parameter -key
	if *bookKey == "" {
		fmt.Println("Error: Specify book key with -key.")
		return
	}
	
	//kontrola ci je zadany parameter -sort
	if *sortOrder != "" && *sortOrder != "asc" && *sortOrder != "desc" {
		fmt.Println("Error: You need to specify a valid sort order -sort (asc or desc).")
		return
	}

	//ziskanie kluca autora
	authorKey, errorKey := getAuthorKey(*bookKey)
	if errorKey != nil {
		fmt.Println("Error getting author name", errorKey)
		return
	}

	var authorNames []string
	var keyField []string
	
	for _, author := range authorKey.Authors {
		authorKey := strings.TrimPrefix(author.Author.Key, "/authors/")
		keyField = append(keyField, authorKey) //pridanie kluca autora do pola
		name, errorKey := getAuthorName(authorKey)
		if errorKey != nil {
			fmt.Println("Error getting author name:", errorKey)
			continue
		}
		authorNames = append(authorNames, name)
		
	}

	//zoradenie mien autorov podla abecedy asc alebo desc
	if *sortOrder == "asc" {
		sort.SliceStable(authorNames, func(i, j int) bool {
			return authorNames[i] < authorNames[j]
		})
	} else if *sortOrder == "desc" {
		sort.SliceStable(authorNames, func(i, j int) bool {
			return authorNames[i] > authorNames[j]
		})
	}
	
	for i, name := range authorNames {
		authorKey := keyField[i]
		myStruct := MyOutput{}
		myStruct.Author.AuthorName = name

		yamlAuthorName, errorAuthorName := yaml.Marshal(myStruct.Author)
		if errorAuthorName != nil {
			fmt.Println("Error marshaling MyStruct to YAML ", errorAuthorName)
			continue
		}
		fmt.Println(string(yamlAuthorName))

		bookInfo, errorBookInfo := getBookInfo(authorKey, *sortOrder)
		if errorBookInfo != nil {
			fmt.Println("Error getting book info for author", name, ":", errorBookInfo)
			continue
		}

		for _, book := range bookInfo {
			if len(book.Isbn) == 0 {
				book.Isbn = []string{"ISBN not found"}
			}

			myStruct.Books = struct {
				Title           string   `json:"title"`
				FirstPublishYear int      `json:"first_publish_year"`
				Isbn             []string `json:"isbn"`
				EditionCount     int      `json:"edition_count"`
			}{
			Title:          book.Title,
			FirstPublishYear: book.FirstPublishYear,
			Isbn:            book.Isbn,
			EditionCount:    book.EditionCount,
			}

			yamlBookInfo, errorBookInfo := yaml.Marshal(myStruct.Books)
			if errorBookInfo != nil {
				fmt.Println("Error marshaling MyStruct to YAML:", errorBookInfo)
				continue
			}
			fmt.Println(string(yamlBookInfo))

		}

	}
}	


//funcia na ziskanie kluca autora
func getAuthorKey(bookKey string) (AuthorKey, error) {
	responseKey, errorKey := http.Get(fmt.Sprintf(findAuthorKeyPI, bookKey))
	if errorKey != nil {
		return AuthorKey{}, errorKey
	}
	defer responseKey.Body.Close() 
	bodyKey, errorKey := ioutil.ReadAll(responseKey.Body) //nacitanie dat z response_key
	if errorKey != nil {
		return AuthorKey{}, errorKey
	}

	var authorKey AuthorKey
	errorKey = json.Unmarshal(bodyKey, &authorKey)
	if errorKey != nil {
		return AuthorKey{}, errorKey
	}
	
	return authorKey, nil
}

//funcia na ziskanie mena autora
func getAuthorName(authorKey string) (string, error) {
	responseAuthorName, errorAuthorName := http.Get(fmt.Sprintf(findAuthorNameAPI, authorKey))
	if errorAuthorName != nil {
		return "", errorAuthorName
	}
	defer responseAuthorName.Body.Close()
	bodyAuthorName, errorAuthorName := ioutil.ReadAll(responseAuthorName.Body)
	if errorAuthorName != nil {
		return "", errorAuthorName
	}

	var authorName AuthorName
	errorAuthorName = json.Unmarshal(bodyAuthorName, &authorName)
	if errorAuthorName != nil {
		return "", errorAuthorName
	}
	return authorName.Name, nil
}

//funcia na ziskanie informacii o knihach autora
func getBookInfo(authorKey string, sortOrder string) ([]BookInfo, error) {
	responseBookInfo, errorBookInfo := http.Get(fmt.Sprintf(findAuthorBooksInfoAPI, authorKey))
	if errorBookInfo != nil {
		return nil, errorBookInfo
	}
	defer responseBookInfo.Body.Close()
	bodyBookInfo, errorBookInfo := ioutil.ReadAll(responseBookInfo.Body)
	if errorBookInfo != nil {
		return nil, errorBookInfo
	}
	var BookInfo struct {
		Docs []BookInfo `json:"docs"`
	}
	errorBookInfo = json.Unmarshal(bodyBookInfo, &BookInfo)
	if errorBookInfo != nil {
		return nil, errorBookInfo
	}


	//zoradenie knih podla roku prveho vydania asc alebo desc
	if sortOrder == "asc" {
		sort.Slice(BookInfo.Docs, func(i, j int) bool {
			return BookInfo.Docs[i].FirstPublishYear < BookInfo.Docs[j].FirstPublishYear
		})
	} else if sortOrder == "desc" {
		sort.Slice(BookInfo .Docs, func(i, j int) bool {
			return BookInfo.Docs[i].FirstPublishYear > BookInfo.Docs[j].FirstPublishYear
		})
	}
	return BookInfo.Docs, nil
}



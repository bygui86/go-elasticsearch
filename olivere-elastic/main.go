package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
)

const (
	esUrl = "http://localhost:9200"
	index = "students"
)

func main() {
	ctx := context.Background()
	esClient, err := GetESClient()
	if err != nil {
		fmt.Printf("Error initializing: %s \n", err.Error())
		panic("Client fail ")
	}

	IndexNewDoc(esClient, ctx)

	QueryDoc(esClient, ctx)
}

type Student struct {
	Name         string  `json:"name"`
	Age          int64   `json:"age"`
	AverageScore float64 `json:"average_score"`
}

func GetESClient() (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(esUrl),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Println("ES initialized...")

	return client, err
}

func IndexNewDoc(esClient *elastic.Client, ctx context.Context) {
	newStudent := Student{
		Name:         "Gopher doe",
		Age:          10,
		AverageScore: 99.9,
	}

	dataJSON, jsonErr := json.Marshal(newStudent)
	if jsonErr != nil {
		panic(jsonErr)
	}

	js := string(dataJSON)
	response, indexErr := esClient.Index().
		Index(index).
		BodyJson(js).
		Do(ctx)
	if indexErr != nil {
		panic(indexErr)
	}

	fmt.Printf("Doc insertion response: %+v \n", response)
	fmt.Println("Doc insertion successful")
}

func QueryDoc(esClient *elastic.Client, ctx context.Context) {
	var students []Student

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("name", "Doe"))

	/* this block will basically print out the es query */
	queryStr, err1 := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)

	if err1 != nil || err2 != nil {
		fmt.Println("err during query marshal=", err1, err2)
	}
	fmt.Println("Final ESQuery=\n", string(queryJs))
	/* until this block */

	searchService := esClient.Search().Index("students").SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return
	}

	for _, hit := range searchResult.Hits.Hits {
		var student Student
		err := json.Unmarshal(hit.Source, &student)
		if err != nil {
			fmt.Println("[Getting Students][Unmarshal] Err=", err)
		}

		students = append(students, student)
	}

	if err != nil {
		fmt.Println("Fetching student fail: ", err)
	} else {
		for _, s := range students {
			fmt.Printf("Student found Name: %s, Age: %d, Score: %f \n", s.Name, s.Age, s.AverageScore)
		}
	}
}

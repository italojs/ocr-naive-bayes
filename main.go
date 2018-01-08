package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

func distinct(elements *[]string) (result []string) {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}

	for v := range *elements {
		if encountered[(*elements)[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[(*elements)[v]] = true
			// Append to result slice.
			result = append(result, (*elements)[v])
		}
	}
	// Return the new slice.
	return result
}

func getCollum(elements *[][]string, ColumnIndex int) (column []string) {
	for i := range *elements {
		if i != 0 {
			column = append(column, (*elements)[i][ColumnIndex])
		}
	}
	return
}

func readFile(filePath string) (records [][]string) {
	//open the file train.csv
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Error opening the file train.csv: ", err)
	}
	defer file.Close()

	//read the opened file train.csv
	reader := csv.NewReader(file)
	records, err = reader.ReadAll()
	if err != nil {
		log.Fatalln("Error read file train.csv: ", err)
	}
	return
}

func sumColumn(column *[]string) (sum int) {
	var wgl sync.WaitGroup
	wgl.Add(len(*column))
	for i := range *column {
		go func(i int) {
			defer wgl.Done()
			if (*column)[i] != "" {
				value, err := strconv.Atoi((*column)[i])
				if err != nil {
					log.Fatalln("Erro:", err)
				}
				sum += value
			}
		}(i)
	}
	wgl.Wait()
	return
}

func getRecorderByClass(records [][]string, class string) (newRecords [][]string, err error) {
	newRecords = make([][]string, 1)
	for l := range records {
		if records[l][0] == class {
			newRecords = append(newRecords, records[l])
		}
	}
	return
}

func sumMatrix(matrix *[][]string) (value int) {

	for l := range *matrix {
		for c := range (*matrix)[l] {
			number, _ := strconv.Atoi((*matrix)[l][c])
			value += number
		}
	}
	return
}

func main() {
	records := readFile("data.csv")
	classes := getCollum(&records, 0)
	classes = distinct(&classes)

	//choose some csv line
	//obs: the first colum is the class(the number that's written on picture)
	line := records[63]

	lenC := len(records[0])
	//lenL := len(records)

	var wgCl sync.WaitGroup
	wgCl.Add(len(classes))

	for cl := range classes {
		go func(cl int) {
			defer wgCl.Done()

			recordsByClass, err := getRecorderByClass(records, classes[cl])
			if err != nil {
				fmt.Println("erro:", err)
			}

			denominator := (sumMatrix(&recordsByClass) + lenC)

			//prob := float64(len(recordsByClass)) / float64(lenL)
			prob := 1.0

			for c := range line {
				if c != 0 {
					if line[c] != "0" {
						colum := getCollum(&recordsByClass, c)

						numerator := sumColumn(&colum) + 1
						train := float64(numerator) / float64(denominator)
						//time, _ := strconv.ParseFloat(line[c], 64)
						//prob *= ((train + 1) * time)
						prob *= (train + 1)
					}
				}
			}
			mgn := fmt.Sprintf("Class: %s | prob: %v | answer: %s", classes[cl], prob, line[0])
			fmt.Println(mgn)
		}(cl)
	}
	wgCl.Wait()
}

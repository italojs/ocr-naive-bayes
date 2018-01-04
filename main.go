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
		column = append(column, (*elements)[i][ColumnIndex])
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

//TODO: simplyfy this method
func getValuesByClass(records *[][]string, class string) (newRecords [][]string, err error) {
	var lAux, cAux int
	lenC := len((*records)[0])
	newRecords = make([][]string, 1)
	for l := range *records {
		for c := range (*records)[l] {
			if (*records)[l][0] == class {
				if cAux == lenC {
					lAux++
					newRecords = append(newRecords, []string{})
					newRecords[lAux] = append(newRecords[lAux], (*records)[l][c])
				} else {
					cAux++
					newRecords[lAux] = append(newRecords[lAux], (*records)[l][c])
				}

			}
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

func multiplyArray(elements *[]float32) (result float64) {
	result = 1.0
	for i := range *elements {
		result *= float64((*elements)[i])
	}
	return
}

func main() {
	records := readFile("data.csv")
	columClassifier := getCollum(&records, 0)
	classes := distinct(&columClassifier)

	train := make([][]float32, len(classes))
	lenC := len(records[0])

	var wgCl sync.WaitGroup
	wgCl.Add(len(classes))

	for cl := range classes {
		go func(cl int) {
			defer wgCl.Done()

			train[cl] = make([]float32, lenC)

			recordsByClass, err := getValuesByClass(&records, classes[cl])
			if err != nil {
				fmt.Println("erro:", err)
			}

			denominator := (sumMatrix(&recordsByClass) + lenC)

			fmt.Println("Colunas: ", lenC)
			for c := range records[0] {

				colum := getCollum(&records, c)

				numerator := sumColumn(&colum) + 1

				trainResult := float32(numerator) / float32(denominator)
				// fmt.Sprintf("numerator: %g  | denominador: %g | result: %g", float32(numerator), float32(denominator), trainResult)
				train[cl][c] = trainResult
			}

		}(cl)
	}
	wgCl.Wait()

	for cl := range classes {
		var result float64
		result = 0.0
		for c := range records[2] {
			result += float64(train[cl][c])
		}
		fmt.Sprintf("Classe: %s  | Pro: %g", classes[cl], result)
	}

}

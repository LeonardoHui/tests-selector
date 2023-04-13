package filesystem

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func GetFileList(dir string) []fs.FileInfo {
	// open the directory and get a list of file names
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return files
}

// specify the directory path
func FeaturesToSingleLine(dir string, file fs.FileInfo) []string {

	// create a 2D slice to hold the results
	results := make([]string, 0)

	if file.Mode().IsRegular() {
		// read the file contents
		content, err := ioutil.ReadFile(dir + "/" + file.Name())
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		// create a scanner to read the file line by line
		scanner := bufio.NewScanner(strings.NewReader(string(content)))

		// create a variable to hold the current scenario
		var scenario string

		// iterate over the lines in the file
		for scanner.Scan() {
			line := scanner.Text()

			// check if the line contains a Gherkin scenario keyword
			if strings.HasPrefix(strings.TrimSpace(line), "Scenario:") {
				// if there was a previous scenario, add it to the results slice
				if scenario != "" {
					results = append(results, scenario)
					scenario = ""
				}

				// extract the scenario name
				scenarioName := strings.TrimSpace(strings.TrimPrefix(line, "Scenario:"))

				// add the scenario name to the scenario string
				scenario += fmt.Sprintf("%s: ", scenarioName)
			} else if strings.HasPrefix(strings.TrimSpace(line), "Given ") ||
				strings.HasPrefix(strings.TrimSpace(line), "When ") ||
				strings.HasPrefix(strings.TrimSpace(line), "Then ") {
				// if the line contains a Gherkin step keyword, add it to the scenario string
				scenario += fmt.Sprintf("%s ", strings.TrimSpace(line))
			} else if strings.TrimSpace(line) == "" {
				// if the line is empty, skip it
				continue
			} else {
				// if the line is not a step or scenario name, it's the end of the scenario
				// add the scenario to the results slice and reset the scenario variable
				if scenario != "" {
					results = append(results, scenario)
					scenario = ""
				}
			}
		}

		// if there was a previous scenario, add it to the results slice
		if scenario != "" {
			results = append(results, scenario)
		}
	}

	return results
}

func AppendToCSV(filename string, tag string, scenario string, values []float64) error {
	// open the CSV file for appending
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// create a CSV writer and write the scenario and values to a new row
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// convert the float values to strings
	valueStrings := make([]string, len(values))
	for i, v := range values {
		valueStrings[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}

	// create a new row with the filename, scenario string, and comma-separated float values
	row := []string{tag, scenario, strings.Join(valueStrings, ",")}
	if err := writer.Write(row); err != nil {
		return err
	}

	return nil
}

func SaveToCSV(results [][]string) {
	// create a new CSV file
	file, err := os.Create("results.csv")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// create a CSV writer
	writer := csv.NewWriter(file)

	// write the results to the CSV file
	for _, result := range results {
		err := writer.Write(result)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	// flush the CSV writer
	writer.Flush()

	// close the CSV file
	file.Close()
}

func ReadFloatArrayFromCSV(filename string) ([][]float64, error) {
	// open the CSV file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a CSV reader
	reader := csv.NewReader(file)

	// create a slice to hold the float arrays
	var arrays [][]float64

	// loop through each row and convert the third column to a float array
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				// end of file
				break
			}
			return nil, err
		}

		// parse the third column as a float array
		arrayStr := strings.Trim(row[2], "[]")
		if arrayStr == "" {
			arrays = append(arrays, []float64{})
			continue
		}
		arrayParts := strings.Split(arrayStr, ",")
		var array []float64
		for _, part := range arrayParts {
			value, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
			if err != nil {
				return nil, err
			}
			array = append(array, value)
		}
		arrays = append(arrays, array)
	}

	if len(arrays) == 0 {
		return nil, errors.New("no float arrays found in CSV file")
	}

	return arrays, nil
}

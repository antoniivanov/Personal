package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	funk "github.com/thoas/go-funk"
)

func printCsvHeader(w *bufio.Writer, numTasks int) {

	sb := strings.Builder{}
	sb.WriteString("FN,Name,Group,") // група
	for i := 1; i <= numTasks; i++ {
		sb.WriteString(fmt.Sprintf("Task%[1]d,%[1]d Results,%[1]d Passed,%[1]d Test Count,%[1]d Tests Score,", i))
	}
	sb.WriteString("TotalScore")
	fmt.Fprintln(w, sb.String())
}

func printToCsv(resultFolder string, studentHws []*StudentHomeWork, numTasks int) {
	os.MkdirAll(resultFolder, 0777)
	csvFile := filepath.Join(resultFolder, "result.csv")

	file, err := os.Create(csvFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	printCsvHeader(w, numTasks)

	for _, shw := range studentHws {
		logsBuilder := strings.Builder{}

		sb := strings.Builder{}
		sb.WriteString(shw.student.FacultyNumber)
		sb.WriteString(",")
		sb.WriteString(shw.student.FullName())
		sb.WriteString(",")
		sb.WriteString(shw.student.Group)
		sb.WriteString(",")
		var totalScore float64
		for i := 1; i <= numTasks; i++ {
			hwSol := shw.homeWorkSolutions[i]
			sb.WriteString(strconv.Itoa(i))
			sb.WriteString(",")
			sb.WriteString(strings.Join(hwSol.TestResults, ";"))
			sb.WriteString(",")
			passed := 0
			funk.ForEach(hwSol.TestResults, func(res string) {
				if res == "OK" || res == "POK" {
					passed++
				}
			})
			sb.WriteString(strconv.Itoa(passed))
			sb.WriteString(",")
			sb.WriteString(strconv.Itoa(len(hwSol.TestResults)))
			sb.WriteString(",")
			var score float64
			if len(hwSol.TestResults) > 0 {
				score = float64(passed) / float64(len(hwSol.TestResults))
			}
			sb.WriteString(strconv.FormatFloat(score, 'f', 3, 64))
			sb.WriteString(",")
			totalScore += score

			funk.ForEach(hwSol.TestLogs, func(s string) { logsBuilder.WriteString(s) })
		}
		sb.WriteString(strconv.FormatFloat(totalScore, 'f', 3, 64))
		fmt.Fprintln(w, sb.String())

		if logsBuilder.Len() > 0 {
			errorsFile := filepath.Join(resultFolder, shw.student.FacultyNumber+".errors.txt")
			ioutil.WriteFile(errorsFile, []byte(logsBuilder.String()), 0644)
		}

	}
	w.Flush()
}

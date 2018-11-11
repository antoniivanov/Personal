package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
)

// Student is a student
type Student struct {
	FacultyNumber string
	Name          string
	Group         string
}

// FullName returns all names
func (student Student) FullName() string {
	return student.Name
}

func parseStudentsInfo(fileName string) []Student {
	csvFile, _ := os.Open(fileName)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var people []Student
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		}
		if error != nil {
			log.Fatalf("Cannot read the file %s", fileName)
		}
		student := Student{FacultyNumber: line[0]}
		student.Name = line[0]
		if len(line) > 5 {
			student.Group = line[5]
		}

		people = append(people, student)
	}
	return people
}

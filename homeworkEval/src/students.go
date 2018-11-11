package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

// Student is a student
type Student struct {
	FacultyNumber string
	FirstName     string
	SecondName    string
	FamilyName    string
	Group         string
}

// FullName returns all names
func (student Student) FullName() string {
	return student.FirstName + " " + student.SecondName + " " + student.FamilyName
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
		names := strings.Split(line[1], " ")
		student.FirstName = names[0]
		student.SecondName = names[1]
		student.FamilyName = names[2]
		if len(names) > 5 {
			student.Group = names[5]
		}
		people = append(people, student)
	}
	return people
}

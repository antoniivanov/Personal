package main

import (
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	solutionsRootFolder := getEnv("SOLUTIONS_DIR", "/Users/aivanov/Google Drive/FMI/up-2018-kn-8")
	fileName := getEnv("STUDENTS_FILE", "../test-students.csv")
	testsRootFolder := getEnv("TEST_DIR", "/Users/aivanov/Google Drive/FMI/up-2018-kn-8/tests")
	homeWorkNumTasks := 5
	students := parseStudentsInfo(fileName)

	log.Print("Number of students are : ", len(students))

	studentsHomework := mapToStudentHomework(students, homeWorkNumTasks, solutionsRootFolder)

	for _, shw := range studentsHomework {
		checkStudentHomework(shw, testsRootFolder)
	}
	log.Println("Generate csvFile with results")
	//funk.ForEach(studentsHomework, func(shw *StudentHomeWork) { log.Printf("%+v", *shw) })
	printToCsv("../result", studentsHomework, homeWorkNumTasks)
	log.Println("The End.")
}

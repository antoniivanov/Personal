package main

import (
	"log"
	"sync"
)

func startGradingStudents() {
	solutionsRootFolder := getEnv("SOLUTIONS_DIR", "/Users/aivanov/Google Drive/FMI/up-2018-kn-8/hw-1-solutions")
	fileName := getEnv("STUDENTS_FILE", "../all-students.csv")
	testsRootFolder := getEnv("TEST_DIR", "/Users/aivanov/Google Drive/FMI/up-2018-kn-8/hw-1-solutions/tests")
	homeWorkNumTasks := 5
	students := parseStudentsInfo(fileName)

	log.Print("Number of students are : ", len(students))

	studentsHomework := mapToStudentHomework(students, homeWorkNumTasks, solutionsRootFolder)

	maxGoroutines := 10
	guard := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup
	for _, shw := range studentsHomework {
		guard <- struct{}{} // would block if guard channel is already filled
		wg.Add(1)
		go func(shw *StudentHomeWork) {
			defer wg.Done()
			checkStudentHomework(shw, testsRootFolder)
			<-guard
		}(shw)
	}
	wg.Wait()
	log.Println("Generate csvFile with results")
	//funk.ForEach(studentsHomework, func(shw *StudentHomeWork) { log.Printf("%+v", *shw) })
	printToCsv("../result", studentsHomework, homeWorkNumTasks)
	log.Println("The End.")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	startGradingStudents()
}

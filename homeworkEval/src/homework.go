package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	funk "github.com/thoas/go-funk"
)

// HomeWork stores
type HomeWork struct {
	Index        int
	TestResults  []string
	TestLogs     []string
	solutionFile string
}

// StudentHomeWork contains info about the homework tasks
type StudentHomeWork struct {
	student           *Student
	homeWorkSolutions map[int]*HomeWork
}

func (shw StudentHomeWork) String() string {
	return fmt.Sprintf("StudentHomeWork { student: %+v, homeWorkSolutions: %v}", *shw.student, shw.homeWorkSolutions)
}

// NewStudentHomeWork creates new student homework
func NewStudentHomeWork(student *Student) *StudentHomeWork {
	return &StudentHomeWork{student: student, homeWorkSolutions: make(map[int]*HomeWork)}
}

func listStudentDirectories(taskSolutionsFolder string) []os.FileInfo {
	files, _ := ioutil.ReadDir(taskSolutionsFolder)
	return files
}

func updateStudentHomeworkStatus(homeWorkIndex int, taskSolutionsFolder string, students []*StudentHomeWork) {
	studentDirs := listStudentDirectories(taskSolutionsFolder)
	for _, student := range students {
		solutionDir :=
			funk.Find(studentDirs, func(d os.FileInfo) bool { return student.student.FacultyNumber == d.Name() })
		if solutionDir != nil {
			//log.Printf("Has HW: %+v\n", student)
			solutionFilePath := filepath.Join(taskSolutionsFolder, solutionDir.(os.FileInfo).Name(), "solution.cpp")
			student.homeWorkSolutions[homeWorkIndex] = &HomeWork{Index: homeWorkIndex, solutionFile: solutionFilePath}
		} else {
			//log.Printf("NO HW: %+v\n", student)
			student.homeWorkSolutions[homeWorkIndex] = &HomeWork{Index: homeWorkIndex, solutionFile: ""}
		}
	}
}

func mapToStudentHomework(students []Student, homeWorkNumTasks int, solutionsRootFolder string) []*StudentHomeWork {
	homeWorks := funk.Map(students, func(s Student) *StudentHomeWork { return NewStudentHomeWork(&s) }).([]*StudentHomeWork)
	for hwIndex := 1; hwIndex <= homeWorkNumTasks; hwIndex++ {
		log.Printf("Looking for Homework Ex %d\n\n", hwIndex)
		taskSolutionsFolder := filepath.Join(solutionsRootFolder, "task-"+strconv.Itoa(hwIndex+1)+"-solutions")
		updateStudentHomeworkStatus(hwIndex, taskSolutionsFolder, homeWorks)
	}
	return homeWorks
}

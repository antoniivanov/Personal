package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func compileSolution(solutionFile string, outputFile string) ([]byte, error) {

	// g++ -o "$EXE" $CPPOPTS "$SRC"
	cmd := exec.Command("g++", "-o", outputFile, solutionFile)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(cmd, stdout, err)
		return stdout, err
	}
	//log.Println("Compiled successfully", cmd, stdout)
	return nil, nil
}

func formatLog(executableSolutionFile string, message string) string {
	return fmt.Sprintf("\n\nERROR IN  %s\n%s\n", filepath.Base(executableSolutionFile), message)
}

func normalizeAnswer(ans string) string {
	trimedChars := "\n\t "
	ans = strings.Trim(ans, trimedChars)
	return ans
}

func runSingleTest(homeWork *HomeWork, testInputFile string, answerFile string, executableSolutionFile string) {
	testTimeoutSeconds, err := time.ParseDuration(getEnv("TEST_TIMEOUT", "10s"))
	if err != nil {
		log.Fatalln("Invalid TEST_TIMEOUT_SECONDS provided", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeoutSeconds)
	defer cancel()

	cmd := exec.CommandContext(ctx, executableSolutionFile)
	testInputFileHandle, err := os.Open(testInputFile)
	if err != nil {
		log.Fatalln("Cannot open test input file ", testInputFile)
	}
	defer testInputFileHandle.Close()
	cmd.Stdin = bufio.NewReader(testInputFileHandle)

	stdout, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		//log.Println("Timed out")
		homeWork.TestResults = append(homeWork.TestResults, "TO") // TIME_OUT
		return
	}

	if err != nil {
		log.Println(stdout, err)
		homeWork.TestResults = append(homeWork.TestResults, "RE") // RUNTIME_ERROR
		homeWork.TestLogs = append(homeWork.TestLogs, formatLog(executableSolutionFile, string(stdout)))
	} else {
		actualAnswer := normalizeAnswer(string(stdout))
		expectedAnswer := normalizeAnswer(readFileToString(answerFile))
		if actualAnswer == expectedAnswer {
			homeWork.TestResults = append(homeWork.TestResults, "OK")
		} else {
			// our results are floating point or int numbers ignore any input chars before
			re := regexp.MustCompile("[-.0-9]*$")
			newActualAnswer := re.FindString(actualAnswer)
			if newActualAnswer == expectedAnswer {
				homeWork.TestResults = append(homeWork.TestResults, "POK") // POSSIBLY_OK
				msg := fmt.Sprintf("Input :%s, Actual: %s ; Expected: %s", readFileToString(testInputFile), actualAnswer, expectedAnswer)
				log.Printf("Possibly OK: %s", msg)
				homeWork.TestLogs = append(homeWork.TestLogs, formatLog(executableSolutionFile, msg))
			} else {
				homeWork.TestResults = append(homeWork.TestResults, "WA") // WRONG_ANSWER
				msg := fmt.Sprintf("Input :%s, Actual: %s ; Expected: %s", readFileToString(testInputFile), actualAnswer, expectedAnswer)
				log.Printf("Result mismatch: %s", msg)
				homeWork.TestLogs = append(homeWork.TestLogs, formatLog(executableSolutionFile, msg))
			}
		}
	}
}

func runTests(hw *HomeWork, taskTestFolder string, executableSolutionFile string) {

	files, err := ioutil.ReadDir(taskTestFolder)
	if err != nil {
		log.Fatalln("Cannot read taskTestFolder", taskTestFolder, err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".test") {
			testInputFile := filepath.Join(taskTestFolder, f.Name())
			testAnswerFile := testInputFile + ".ans"
			runSingleTest(hw, testInputFile, testAnswerFile, executableSolutionFile)
		}
	}

}

func checkStudentHomework(studentsHw *StudentHomeWork, testsRootFolder string) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "hw_workspace")
	if err != nil {
		log.Fatalln("Cannot create workspace temp directory", err)
	}
	defer os.RemoveAll(tempDir)

	for _, hw := range studentsHw.homeWorkSolutions {
		if hw.solutionFile != "" {
			taskTestFolder := filepath.Join(testsRootFolder, "task"+strconv.Itoa(hw.Index))

			outputFile := "out-" + studentsHw.student.FacultyNumber + "-ex-" + strconv.Itoa(hw.Index)
			outputFile = filepath.Join(tempDir, outputFile)
			if stdout, err := compileSolution(hw.solutionFile, outputFile); err == nil {
				//log.Println(outputFile, err)
				runTests(hw, taskTestFolder, outputFile)
			} else {
				hw.TestResults = append(hw.TestResults, "CE") // COMPILE_ERROR
				hw.TestLogs = append(hw.TestLogs, formatLog(hw.solutionFile, string(stdout)))
			}
		} else {
			hw.TestResults = []string{"NA"} // NOT_AVAILABLE
		}
	}
}

package ioworker

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"log"
	"os"
)

type Students struct {
	Studs []Student
}
type Student struct {
	Index string
	Data  []string
	Grade string
}

type Header struct {
	FeatureName []string
}

type Pair struct {
	Val string
	Key string
}

func (students Students) AttributesById(index int) (result []string) {
	for _, val := range students.Studs {
		result = append(result, val.Data[index])
	}
	return
}

func (students Students) GetWhereEq(index int, targetVal string) (result Students) {
	for _, val := range students.Studs {
		if val.Data[index] == targetVal {
			result.Studs = append(result.Studs, val)
		}
	}
	return
}

func (students Students) GetGrades() (result []string) {
	for _, val := range students.Studs {
		result = append(result, val.Grade)
	}
	return
}

func ParseHeader(record [][]string) Header {
	var header = make([]string, 0, len(record[0]))
	for i := 0; i < len(record[0]); i++ {
		header = append(header, record[0][i])
	}
	return Header{FeatureName: header}
}

func ParseStudents(records [][]string) Students {
	var cntOfStudents = 0
	for i := 0; i < len(records); i++ {
		if len(records[i]) == len(records[0]) {
			cntOfStudents++
		}
	}
	var students = make([]Student, 0, cntOfStudents)
	for i := 1; i < cntOfStudents; i++ {
		var data = make([]string, 0, len(records[i])-1)
		for j := 1; j < len(records[i])-1; j++ {
			data = append(data, records[i][j])
		}
		grade := records[i][len(records[i])-1]
		var student = Student{Data: data, Grade: grade, Index: records[i][0]}
		students = append(students, student)
	}
	return Students{Studs: students}

}

func ReadCSVFile(filePath string) (Header, Students) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)

	records, err := csvReader.ReadAll()
	var header = ParseHeader(records)
	var students = ParseStudents(records)
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	return header, students
}

func Write(filename, sep string, data [][]string, columns int) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Can't create file")
	}
	nl := "\n"
	toWrite := "actual" + sep + "predicted" + nl
	_, err = file.WriteString(toWrite)
	if err != nil {
		return errors.Wrap(err, "Can't write header")
	}

	for _, value := range data {
		if len(value) != columns {
			return errors.New("Length of data is not equal to columns size")
		}
		toWrite = value[0]
		for i := 1; i < len(value); i++ {
			toWrite += sep + value[i]
		}
		toWrite += nl
		_, err = file.WriteString(toWrite)
		if err != nil {
			return errors.Wrap(err, "Can't write line")
		}
	}
	return nil
}

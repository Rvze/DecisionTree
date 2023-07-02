package main

import (
	"AISlab3/src/decision_tree"
	"AISlab3/src/ioworker"
	"fmt"
)

func main() {
	header, students := ioworker.ReadCSVFile("docs/train.csv")

	tree := decision_tree.DecisionTree{Header: header, Students: students}
	root := tree.Start(40)
	_, testStudents := ioworker.ReadCSVFile("docs/test.csv")

	rows := float64(len(testStudents.Studs))
	p := 0.0
	tp := 0.0
	fp := 0.0
	fn := 0.0

	toWrite := make([][]string, len(testStudents.Studs))
	for i, value := range testStudents.Studs {
		actual := value.Grade
		predicted := tree.Predict(root, value)

		if good(actual) == good(predicted) {
			p += 1.0
		}
		if good(actual) && good(predicted) {
			tp += 1.0
		}
		if !good(actual) && good(predicted) {
			fp += 1.0
		}
		if good(actual) && !good(predicted) {
			fn += 1.0
		}
		toWrite[i] = make([]string, 2)
		toWrite[i][0] = goodString(actual)
		toWrite[i][1] = goodString(predicted)
	}
	err := ioworker.Write("docs/result.txt", ",", toWrite, 2)
	if err != nil {
		return
	}

	tree.PrintSelected()

	fmt.Println("Accuracy:", p/rows)
	fmt.Println("Precision:", tp/(tp+fp))
	fmt.Println("Recall:", tp/(tp+fn))
}

func good(grade string) bool {
	return grade < "4"
}

func goodString(grade string) string {
	if good(grade) {
		return "1"
	} else {
		return "0"
	}
}

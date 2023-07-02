package decision_tree

import (
	"AISlab3/src/ioworker"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

type TreeNode struct {
	class    string
	children map[string]*TreeNode
	index    int
}

type DecisionTree struct {
	Header   ioworker.Header
	Students ioworker.Students
	root     *TreeNode
}

func newNode() *TreeNode {
	node := TreeNode{children: make(map[string]*TreeNode)}
	return &node
}
func newLeaf(class string) *TreeNode {
	leaf := TreeNode{class: class}
	return &leaf
}

var selectedAttributes []int
var maxDepth int

func (tree DecisionTree) Start(maxDepthIn int) *TreeNode {
	selectedAttributes = selectAttributes(tree)
	maxDepth = maxDepthIn
	tree.root = newNode()
	tree.root.children = build(tree.Students, tree.root, 0)
	return tree.root
}

func build(students ioworker.Students, parent *TreeNode, depth int) map[string]*TreeNode {
	if parent == nil {
		return nil
	}
	children := make(map[string]*TreeNode)
	index := findBestAttribute(students)
	parent.index = index

	attrVals := students.AttributesById(index)
	uniqueVals := getUniqElements(attrVals)

	for _, val := range uniqueVals {
		selectedStuds := students.GetWhereEq(index, val)
		selectedGrades := selectedStuds.GetGrades()
		selectedUniqGrades := getUniqElements(selectedGrades)

		if len(selectedUniqGrades) == 1 {
			children[val] = newLeaf(selectedUniqGrades[0])
		} else if depth > maxDepth {
			maxGrades := 0
			grade := ""
			for _, attrVal := range selectedUniqGrades {
				freq := freq(attrVal, selectedGrades)
				if freq > maxGrades {
					maxGrades = freq
					grade = attrVal
				}
			}
			children[val] = newLeaf(grade)
		} else {
			newNode := newNode()
			newNode.children = build(selectedStuds, newNode, depth+1)
			children[val] = newNode
		}
	}
	return children
}

func (tree DecisionTree) Predict(root *TreeNode, student ioworker.Student) (grade string) {
	current := root
	for current.children != nil {
		idx := current.index
		value := student.Data[idx]
		if _, ok := current.children[value]; ok {
			current = current.children[value]
		} else {
			min := math.MaxFloat32
			var foundedKey string
			for key := range current.children {
				kInt, _ := strconv.Atoi(key)
				valInt, _ := strconv.Atoi(value)
				diff := math.Abs(float64(kInt - valInt))

				if diff < min {
					foundedKey = key
					min = diff
				}
			}
			current = current.children[foundedKey]
		}
	}
	return current.class
}

func (tree DecisionTree) PrintSelected() {
	fmt.Println("Выбранные аттрибуты:")
	for _, val := range selectedAttributes {
		fmt.Print(tree.Header.FeatureName[val], " ")
	}
	fmt.Println()
}

func findBestAttribute(students ioworker.Students) (res int) {
	grades := students.GetGrades()
	uniqueClasses := getUniqElements(grades)

	info := info(grades, uniqueClasses)
	maxGainRatio := 0.0
	for index := range selectedAttributes {
		infoX := infoX(students, index)
		gainRatio := (info - infoX) / split(students, index)
		if gainRatio > maxGainRatio {
			maxGainRatio = gainRatio
			res = index
		}
	}
	return
}

type void struct {
}

var member void

func selectAttributes(root DecisionTree) (attributes []int) {
	min := 0
	max := len(root.Students.Studs[0].Data)
	n := int(math.Sqrt(float64(max)))
	set := make(map[int]void)
	for len(set) < n {
		set[random(min, max)] = member
	}

	for k := range set {
		attributes = append(attributes, k)
	}
	sort.Slice(attributes, func(i, j int) bool {
		return attributes[i] < attributes[j]
	})
	return
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	if min > max {
		return min
	} else {
		return rand.Intn(max-min) + min
	}
}

func infoX(students ioworker.Students, index int) (res float64) {
	res = 0
	attrVals := students.AttributesById(index)
	uniqVals := getUniqElements(attrVals)

	for _, val := range uniqVals {
		selectedStuds := students.GetWhereEq(index, val)
		selectedGrades := selectedStuds.GetGrades()
		selectedClassesUniq := getUniqElements(selectedGrades)

		res += float64(len(selectedStuds.Studs)) / float64(len(students.Studs)) * info(selectedGrades, selectedClassesUniq)
	}
	return

}

func info(entropy []string, classExample []string) (res float64) {
	for _, class := range classExample {
		freq := freq(class, entropy)
		div := float64(freq) / float64(len(entropy))
		res -= div * (math.Log2(div))
	}
	return
}

func split(students ioworker.Students, index int) (result float64) {
	attrVals := students.AttributesById(index)
	uniqAttrVals := getUniqElements(attrVals)

	for _, val := range uniqAttrVals {
		selectedStuds := students.GetWhereEq(index, val)
		div := float64(len(selectedStuds.Studs)) / float64(len(students.Studs))
		result -= div * math.Log2(div)
	}
	return
}

func freq(feature string, featuresList []string) int {
	cnt := 0
	for _, val := range featuresList {
		if feature == val {
			cnt++
		}
	}
	return cnt
}

func getUniqElements(list []string) (result []string) {
	for _, val := range list {
		if !contains(result, val) {
			result = append(result, val)
		}
	}
	return
}

func contains(list []string, target string) bool {
	for _, val := range list {
		if target == val {
			return true
		}
	}
	return false
}

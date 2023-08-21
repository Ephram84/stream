package stream

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Employee struct {
	ID     int
	Name   string
	Salary float64
}

var sliceEmployee = []Employee{
	{
		ID:     1,
		Name:   "Jeff Bezos",
		Salary: 100000.0,
	},
	{
		ID:     2,
		Name:   "Bill Gates",
		Salary: 200000.0,
	},
	{
		ID:     3,
		Name:   "Mark Zuckerberg",
		Salary: 300000.0,
	},
}

func TestWords(t *testing.T) {
	words, err := FromFile("../assets/words.txt").ForEach(changer).Filter(filter).ToArray()
	if err != nil {
		t.Fatal(err)
	}

	for _, word := range words {
		if len(word) == 1 {
			t.Fatal(word)
		}
		fmt.Println(word)
	}
}

var isWord = regexp.MustCompile(`[A-Za-z]+`)

func changer(word string) string {
	return isWord.FindString(word)
}

func filter(word string) (bool, error) {
	return len(word) > 1, nil
}

func TestSlice(t *testing.T) {
	ints, err := From([][]int{{0, 1, 2}, {3, 4, 5}, {6}}).
		Filter(func(ints []int) (bool, error) { return len(ints) > 1, nil }).
		ToArray()
	if err != nil {
		t.Fatal(err)
	}

	for _, slice := range ints {
		fmt.Println(slice)
	}
}

func TestMatch(t *testing.T) {
	numbers := []int{2, 4, 5, 6, 8}

	allEven, err := From(numbers).AllMatch(isEven)
	assert.NoError(t, err)
	assert.False(t, allEven)

	oneEven, err := From(numbers).AnyMatch(isEven)
	assert.NoError(t, err)
	assert.True(t, oneEven)

	noneMultipleOfThree, err := From(numbers).NoneMatch(func(elem int) (bool, error) {
		return elem%3 == 0, nil
	})
	assert.NoError(t, err)
	assert.False(t, noneMultipleOfThree)
}

func isEven(elem int) (bool, error) {
	return elem%2 == 0, nil
}

func TestWrite(t *testing.T) {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	_, err := From(numbers).Write(os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindFirst(t *testing.T) {
	employee, err := From(sliceEmployee).Filter(func(elem Employee) (bool, error) {
		return elem.Salary > 100000.0, nil
	}).FindFirst()
	assert.NoError(t, err)

	assert.NotNil(t, employee)
	assert.Equal(t, employee.Salary, 200000.0)
	assert.Equal(t, employee.Name, "Bill Gates")
}

func TestFindFirstOrElse(t *testing.T) {
	employee, err := From(sliceEmployee).Filter(func(elem Employee) (bool, error) {
		return elem.Salary > 1000000.0, nil
	}).FindFirst(Employee{Name: "John Doe"})
	assert.NoError(t, err)

	assert.Equal(t, employee.Name, "John Doe")
	assert.Equal(t, employee.Salary, 0.0)
}

func TestMapToFloat64(t *testing.T) {
	salaries, err := From(sliceEmployee).MapToFloat(func(elem Employee) (float64, error) {
		return elem.Salary, nil
	}).ToArray()
	assert.NoError(t, err)

	assert.Equal(t, []float64{100000.0, 200000.0, 300000.0}, salaries)
}

func TestPartitionBy(t *testing.T) {
	numbers := []int{2, 4, 5, 6, 8}

	isEven, err := From(numbers).PartitioningBy(isEven).ToMap()
	assert.NoError(t, err)
	assert.Len(t, isEven[true], 4)
	assert.Len(t, isEven[false], 1)
}

func TestGroupBy(t *testing.T) {
	groupByAlphabet, err := From(sliceEmployee).GroupByString(func(elem Employee) string {
		return elem.Name[0:1]
	}).ToMap()
	assert.NoError(t, err)
	assert.Equal(t, groupByAlphabet["B"][0].Name, "Bill Gates")
	assert.Equal(t, groupByAlphabet["J"][0].Name, "Jeff Bezos")
	assert.Equal(t, groupByAlphabet["M"][0].Name, "Mark Zuckerberg")
}

func TestMapToNewType(t *testing.T) {
	lengthOfNames, err := Map(From(sliceEmployee), func(elem Employee) (int, error) {
		return len(elem.Name), nil
	}).ToArray()
	assert.NoError(t, err)
	assert.Equal(t, lengthOfNames, []int{10, 10, 15})
}

func TestMapToInt(t *testing.T) {
	lengthOfNames, err := From(sliceEmployee).MapToInt(func(elem Employee) (int, error) {
		return len(elem.Name), nil
	}).ToArray()
	assert.NoError(t, err)
	assert.Equal(t, lengthOfNames, []int{10, 10, 15})
}

func TestWordCount(t *testing.T) {
	wordCount, err := FromFile("../assets/words.txt").GroupByString(mapper).Reducing(reducer)
	assert.NoError(t, err)

	for word, count := range wordCount {
		fmt.Println(word, count)
	}
	fmt.Println(len(wordCount))
}

func mapper(word string) string {
	return isWord.FindString(strings.ToLower(word))
}

func reducer(key string, values []string) (string, string) {
	return key, strconv.Itoa(len(values))
}

func TestMax(t *testing.T) {
	maxEmployer, err := From(sliceEmployee).Max(func(max, elem Employee) bool {
		return max.Salary < elem.Salary
	})
	assert.NoError(t, err)
	assert.NotNil(t, maxEmployer)
	assert.Equal(t, 3, maxEmployer.ID)

	max, err := From([]int{}).Max(MaxInt)
	assert.NoError(t, err)
	assert.Nil(t, max)
}

func TestMapSlice(t *testing.T) {
	accounts := []Account{
		{
			Transactions: []Transaction{
				{
					ID:     "1",
					Amount: 120.0,
				},
				{
					ID:     "2",
					Amount: -120.0,
				},
			},
		},
		{
			Transactions: []Transaction{
				{
					ID:     "3",
					Amount: 20.0,
				},
				{
					ID:     "4",
					Amount: 10.0,
				},
			},
		},
	}

	numberOfTransactions, err := MapSlice[Account, Transaction](From(accounts, len(accounts)), func(elem Account) ([]Transaction, error) {
		return elem.Transactions, nil
	}).Filter(func(elem Transaction) (bool, error) {
		return elem.Amount > 0.0, nil
	}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 3, numberOfTransactions)
}

type Account struct {
	ID           string
	Transactions []Transaction
}

type Transaction struct {
	ID     string
	Amount float64
}

func TestSort(t *testing.T) {
	numbers := []int{5, 3, 1, 2, 4}
	result, err := From(numbers).ToSortedArray(SortInts)
	assert.NoError(t, err)
	assert.True(t, sort.IntsAreSorted(result))
}

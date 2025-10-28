package stream

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/Ephram84/stream/stream/common"
	"github.com/Ephram84/stream/stream/tupel"
	"github.com/stretchr/testify/assert"
)

// help functions, structures and variables for tests

func isEven(elem int) (bool, error) {
	return elem%2 == 0, nil
}

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

func TestError(t *testing.T) {
	numbers := []string{"1", "2", "a", "4"}
	result, err := From(numbers).MapToInt(common.StringToInt).Filter(isEven).ToSlice()
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestWords(t *testing.T) {
	words, err := FromFile("../assets/words.txt").GroupByString(mapper).CountValues().ToMap()
	assert.NoError(t, err)

	assert.Equal(t, 127, words["a"])
	assert.Equal(t, 141, words["ac"])
	assert.Equal(t, 68, words["luctus"])
}

var isWord = regexp.MustCompile(`[A-Za-z]+`)

func TestSlice(t *testing.T) {
	ints, err := From([][]int{{0, 1, 2}, {3, 4, 5}, {6}}).
		Filter(func(ints []int) (bool, error) { return len(ints) > 1, nil }).
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, [][]int{{0, 1, 2}, {3, 4, 5}}, ints)
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

func TestWrite(t *testing.T) {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	_, err := From(numbers).Write(os.Stdout)
	assert.NoError(t, err)
}

func TestFindFirst(t *testing.T) {
	employee, err := From(sliceEmployee).Filter(func(elem Employee) (bool, error) {
		return elem.Salary > 100000.0, nil
	}).First()
	assert.NoError(t, err)

	assert.NotNil(t, employee)
	assert.Equal(t, employee.Salary, 200000.0)
	assert.Equal(t, employee.Name, "Bill Gates")
}

func TestFindFirstOrElse(t *testing.T) {
	employee, err := From(sliceEmployee).Filter(func(elem Employee) (bool, error) {
		return elem.Salary > 1000000.0, nil
	}).First(Employee{Name: "John Doe"})
	assert.NoError(t, err)

	assert.Equal(t, employee.Name, "John Doe")
	assert.Equal(t, employee.Salary, 0.0)
}

func TestMapToFloat(t *testing.T) {
	salaries, err := From(sliceEmployee).MapToFloat(func(elem Employee) (float64, error) {
		return elem.Salary, nil
	}).ToSlice()
	assert.NoError(t, err)

	assert.Equal(t, []float64{100000.0, 200000.0, 300000.0}, salaries)
}

func TestPartitionBy(t *testing.T) {
	numbers := []int{2, 4, 5, 6, 8}

	isEven, err := From(numbers).PartitioningBy(isEven).ToMap()
	assert.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6, 8}, isEven[true])
	assert.Equal(t, []int{5}, isEven[false])
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

func TestMapToInt(t *testing.T) {
	lengthOfNames, err := From(sliceEmployee).MapToInt(func(elem Employee) (int, error) {
		return len(elem.Name), nil
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, lengthOfNames, []int{10, 10, 15})
}

func TestMaxWithInts(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	maxInt, err := From(numbers).Reduce(0, accumulator.Max(common.MaxInt))
	assert.NoError(t, err)
	assert.Equal(t, 9, maxInt)
}

func TestMax(t *testing.T) {
	maxEmployer, err := From(sliceEmployee).Reduce(Employee{}, accumulator.Max(func(max, elem Employee) bool {
		return max.Salary < elem.Salary
	}))
	assert.NoError(t, err)
	assert.NotEqual(t, Employee{}, maxEmployer)
	assert.Equal(t, 3, maxEmployer.ID)
}

func TestMin(t *testing.T) {
	minEmployer, err := From(sliceEmployee).Reduce(sliceEmployee[2], accumulator.Min(func(min, elem Employee) bool {
		return min.Salary > elem.Salary
	}))
	assert.NoError(t, err)
	assert.Equal(t, 1, minEmployer.ID)
}

func TestFlatMapInts(t *testing.T) {
	intSlice := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	ints, err := FlatMapSlice(From(intSlice), func(elem []int) ([]int, error) {
		return elem, nil
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, ints)
}

func TestFlatMap(t *testing.T) {
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

	numberOfTransactions, err := FlatMapSlice(From(accounts), func(elem Account) ([]Transaction, error) {
		return elem.Transactions, nil
	}).Filter(func(elem Transaction) (bool, error) {
		return elem.Amount > 0.0, nil
	}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 3, numberOfTransactions)
}

func TestFlatMapWithEmptyAccounts(t *testing.T) {
	accounts := []Account{}

	numberOfTransactions, err := FlatMapSlice(From(accounts), func(elem Account) ([]Transaction, error) {
		return elem.Transactions, nil
	}).Filter(func(elem Transaction) (bool, error) {
		return elem.Amount > 0.0, nil
	}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 0, numberOfTransactions)
}

type Account struct {
	ID           string
	Transactions []Transaction
}

type Transaction struct {
	ID          string
	Amount      float64
	BookingDate int64
	Tags        []string
}

func TestSort(t *testing.T) {
	numbers := []int{5, 3, 1, 2, 4}
	result, err := From(numbers).Sort(SortInts).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestDistinct(t *testing.T) {
	numbers := []int{5, 5, 3, 1, 1, 2, 4}
	result, err := From(numbers).Distinct(common.Eq).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{5, 3, 1, 2, 4}, result)
}

func TestDistinctObjects(t *testing.T) {
	transactions := []Transaction{
		{
			ID:     "T1",
			Amount: 12.5,
		},
		{
			ID:     "T1",
			Amount: 14.5,
		},
		{
			ID:     "T2",
			Amount: -50.0,
		},
	}

	result, err := From(transactions).Distinct(func(a, b Transaction) bool {
		return a.ID == b.ID
	}).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []Transaction{
		{
			ID:     "T1",
			Amount: 12.5,
		},
		{
			ID:     "T2",
			Amount: -50.0,
		},
	}, result)
}

func TestWithNils(t *testing.T) {
	accounts := []*Account{
		{},
		nil,
		{},
	}
	n, err := From(accounts).Filter(func(elem *Account) (bool, error) {
		return elem != nil, nil
	}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
}

func TestSum(t *testing.T) {
	numbers := []int{5, 3, 1, 2, 4}
	result, err := From(numbers).Reduce(0, accumulator.Sum[int]())
	assert.NoError(t, err)
	assert.Equal(t, 15, result)
}

func TestAvg(t *testing.T) {
	numbers := []int{5, 3, 1, 2, 4}
	result, err := From(numbers).MapToFloat(common.IntToFloat64).Reduce(0.0, accumulator.Avg())
	assert.NoError(t, err)
	assert.Equal(t, 3.0, result)
}

func TestReducingToFloat(t *testing.T) {
	transactions := []Transaction{
		{
			Amount:      10.0,
			BookingDate: time.Date(2023, time.September, 9, 0, 0, 0, 0, time.UTC).Unix(),
		},
		{
			Amount:      15.0,
			BookingDate: time.Date(2023, time.September, 9, 0, 0, 0, 0, time.UTC).Unix(),
		},
		{
			Amount:      150.0,
			BookingDate: time.Date(2023, time.August, 20, 0, 0, 0, 0, time.UTC).Unix(),
		},
	}

	result, err := From(transactions).GroupByString(func(elem Transaction) string {
		date := time.Unix(elem.BookingDate, 0)
		return fmt.Sprintf("%d-%d", date.Year(), int(date.Month()))
	}).MapToFloat64(func(elem Transaction) (float64, error) { return elem.Amount, nil }).Reduce(0.0, accumulator.Sum[float64]()).ToMap()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{
		"2023-9": 25.0,
		"2023-8": 150.0,
	}, result)
}

func TestAssociateByString(t *testing.T) {
	transactions := []Transaction{
		{
			ID:     "T01",
			Amount: 10.0,
		},
		{
			ID:     "T02",
			Amount: 20.0,
		},
		{
			ID:     "T03",
			Amount: 30.0,
		},
	}

	result, err := From(transactions).AssociateByString(func(elem Transaction) (string, error) {
		return elem.ID, nil
	}).ToMap()
	assert.NoError(t, err)
	assert.Equal(t, map[string]Transaction{
		"T01": {
			ID:     "T01",
			Amount: 10.0,
		},
		"T02": {
			ID:     "T02",
			Amount: 20.0,
		},
		"T03": {
			ID:     "T03",
			Amount: 30.0,
		},
	}, result)
}

func TestConcatSlices(t *testing.T) {
	result, err := From([]int{4, 5}).Concat(From([]int{6, 7, 8})).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{4, 5, 6, 7, 8}, result)
}

func TestZipSlices(t *testing.T) {
	result, err := ZipSlices(From([]int{1, 2, 3}), From([]string{"a", "b", "c", "d"})).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []tupel.Tupel[int, string]{
		{First: 1, Second: "a"},
		{First: 2, Second: "b"},
		{First: 3, Second: "c"},
	}, result)
}

func TestGroupBySlice(t *testing.T) {
	numbers := From([]int{1, 2, 3, 4, 5, 6})
	result := GroupBySlice(numbers, func(elem int) (string, error) {
		if elem%2 == 0 {
			return "even", nil
		}
		return "odd", nil
	})
	assert.NotNil(t, result)
	m, err := result.ToMap()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 3, 5}, m["odd"])
	assert.Equal(t, []int{2, 4, 6}, m["even"])
}

func TestTakeSlice(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).
		Take(3).
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestSkipSlice(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).
		Skip(2).
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestReverseSlice(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).
		Reverse().
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{5, 4, 3, 2, 1}, result)
}

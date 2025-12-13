package slice

import (
	"os"
	"testing"

	"github.com/Ephram84/stream/stream/accumulator"
	"github.com/Ephram84/stream/stream/common"
	"github.com/Ephram84/stream/stream/tupel"
	"github.com/stretchr/testify/assert"
)

// help functions, structures and variables for tests

func isEven(elem int) (bool, error) {
	return elem%2 == 0, nil
}

func mapper(s string) (string, error) {
	return s, nil
}

type Employee struct {
	ID     int
	Name   string
	Salary float64
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
	groupByAlphabet, err := From(sliceEmployee).GroupByString(func(elem Employee) (string, error) {
		return elem.Name[0:1], nil
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
	maxInt, err := From(numbers).Reduce(accumulator.Max)
	assert.NoError(t, err)
	assert.Equal(t, 9, maxInt)
}

func TestCustomMax(t *testing.T) {
	maxEmployer, err := From(sliceEmployee).Reduce(func(values []Employee) (Employee, error) {
		max := values[0]
		maxSalary := values[0].Salary
		for _, employee := range values {
			if employee.Salary > maxSalary {
				max = employee
				maxSalary = employee.Salary
			}
		}
		return max, nil
	})
	assert.NoError(t, err)
	assert.NotEqual(t, Employee{}, maxEmployer)
	assert.Equal(t, 3, maxEmployer.ID)
}

func TestMin(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	minInt, err := From(numbers).Reduce(accumulator.Min)
	assert.NoError(t, err)
	assert.Equal(t, 1, minInt)
}

func TestFlatMapInts(t *testing.T) {
	intSlice := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	ints, err := FlatMap(From(intSlice), func(elem []int) ([]int, error) {
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

	numberOfTransactions, err := FlatMap(From(accounts), func(elem Account) ([]Transaction, error) {
		return elem.Transactions, nil
	}).Filter(func(elem Transaction) (bool, error) {
		return elem.Amount > 0.0, nil
	}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 3, numberOfTransactions)
}

func TestFlatMapWithEmptyAccounts(t *testing.T) {
	accounts := []Account{}

	numberOfTransactions, err := FlatMap(From(accounts), func(elem Account) ([]Transaction, error) {
		return elem.Transactions, nil
	}).Filter(func(elem Transaction) (bool, error) {
		return elem.Amount > 0.0, nil
	}).Count()
	assert.NoError(t, err)
	assert.Equal(t, 0, numberOfTransactions)
}

func TestSort(t *testing.T) {
	numbers := []int{5, 3, 1, 2, 4}
	result, err := From(numbers).Sort(common.Sort).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestSortDesc(t *testing.T) {
	numbers := []int{5, 3, 1, 2, 4}
	result, err := From(numbers).Sort(common.SortDesc).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{5, 4, 3, 2, 1}, result)
}

func TestCustomSortFunc(t *testing.T) {
	transactions := []Transaction{
		{
			ID:          "T02",
			Amount:      20.0,
			BookingDate: 1759356000, // 2025-10-02
		},
		{
			ID:          "T03",
			Amount:      30.0,
			BookingDate: 1756764000, // 2025-09-02
		},
		{
			ID:          "T01",
			Amount:      10.0,
			BookingDate: 1762038000, // 2025-11-02
		},
	}

	sortTransaction, err := From(transactions).Sort(func(tnxs []Transaction) func(i, j int) bool {
		return func(i, j int) bool {
			return tnxs[i].BookingDate < tnxs[j].BookingDate
		}
	}).ToSlice()
	assert.NoError(t, err)
	assert.ElementsMatch(t, []Transaction{
		{
			ID:          "T01",
			Amount:      10.0,
			BookingDate: 1762038000, // 2025-11-02
		},
		{
			ID:          "T02",
			Amount:      20.0,
			BookingDate: 1759356000, // 2025-10-02
		},
		{
			ID:          "T03",
			Amount:      30.0,
			BookingDate: 1756764000, // 2025-09-02
		},
	}, sortTransaction)
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
	result, err := From(numbers).Reduce(accumulator.Sum)
	assert.NoError(t, err)
	assert.Equal(t, 15, result)
}

func TestAvg(t *testing.T) {
	numbers := []int{5, 3, 1, 2, 3}
	result, err := From(numbers).MapToFloat(common.IntToFloat64).Reduce(accumulator.Avg)
	assert.NoError(t, err)
	assert.Equal(t, 2.8, result)

	floats := []float64{5.0, 3.5, 1.0, 2.0, 4.0}
	resultF, err := From(floats).Reduce(accumulator.Avg)
	assert.NoError(t, err)
	assert.Equal(t, 3.1, resultF)
}

func TestCustomAccumulator(t *testing.T) {
	transactions := []Transaction{
		{
			ID:          "T01",
			Amount:      10.0,
			BookingDate: 1762038000, // 2025-11-02
		},
		{
			ID:          "T02",
			Amount:      20.0,
			BookingDate: 1759356000, // 2025-10-02
		},
		{
			ID:          "T03",
			Amount:      30.0,
			BookingDate: 1756764000, // 2025-09-02
		},
	}

	oldestTransaction, err := From(transactions).Reduce(func(values []Transaction) (Transaction, error) {
		oldest := values[0]
		for _, transaction := range values[1:] {
			if transaction.BookingDate < oldest.BookingDate {
				oldest = transaction
			}
		}

		return oldest, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "T03", oldestTransaction.ID)
	assert.Equal(t, 30.0, oldestTransaction.Amount)
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

func TestConcat(t *testing.T) {
	result, err := From([]int{4, 5}).Concat(From([]int{6, 7, 8})).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{4, 5, 6, 7, 8}, result)
}

func TestZip(t *testing.T) {
	result, err := Zip(From([]int{1, 2, 3}), From([]string{"a", "b", "c", "d"})).ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []tupel.Tupel[int, string]{
		{First: 1, Second: "a"},
		{First: 2, Second: "b"},
		{First: 3, Second: "c"},
	}, result)
}

func TestGroupBy2(t *testing.T) {
	numbers := From([]int{1, 2, 3, 4, 5, 6})
	result := GroupBy(numbers, func(elem int) (string, error) {
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

func TestTake(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).
		Take(3).
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestSkip(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).
		Skip(2).
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestReverse(t *testing.T) {
	result, err := From([]int{1, 2, 3, 4, 5}).
		Reverse().
		ToSlice()
	assert.NoError(t, err)
	assert.Equal(t, []int{5, 4, 3, 2, 1}, result)
}

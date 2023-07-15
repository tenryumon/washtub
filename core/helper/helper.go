package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"
)

// ChunkSliceInt64 creates an array of int64 split into groups with the length of size.
// If array can't be split evenly, the final chunk will be the remaining element.
func ChunkSliceInt64(slice []int64, chunkSize int) [][]int64 {
	var chunks [][]int64

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// Get sql common clause from param that already validated
// order and limit are common clause and always have last position in query
func GetSqlCommonClauses(param map[string]interface{}) string {
	queryOrderBy := " ORDER BY %s %s"
	queryOffsetLimit := " LIMIT %d,%d"

	clauses := ""
	if _, ok := param["sort_by"]; ok {
		if _, ok := param["sort_type"]; ok {
			clauses += fmt.Sprintf(queryOrderBy, param["sort_by"], param["sort_type"])
		}
	}

	if _, ok := param["limit"]; ok {
		if _, ok := param["offset"]; ok {
			clauses += fmt.Sprintf(queryOffsetLimit, param["offset"], param["limit"])
		} else {
			clauses += fmt.Sprintf(queryOffsetLimit, 0, param["limit"])
		}

	}

	return clauses
}

// Adapted from https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int, letters string) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func AgesToDateBetweenClause(field string, ages []int) (string, map[string]interface{}) {
	clauses := []string{}
	param := map[string]interface{}{}
	sort.Ints(ages)
	// groupAge to group age that have 1 year different
	// ex: age 2,4,5,6 will have 2 group -> 2 and 4,5,6
	groupAges := [][]int{}
	groupAge := []int{}
	lastIdx := len(ages) - 1
	for idx, age := range ages {
		lastIdxGroup := len(groupAge) - 1
		isGroupEmpty := len(groupAge) == 0
		if isGroupEmpty || age-groupAge[lastIdxGroup] == 1 {
			groupAge = append(groupAge, age)
		} else {
			groupAges = append(groupAges, groupAge)
			groupAge = []int{}
			groupAge = append(groupAge, age)
		}
		if idx == lastIdx {
			groupAges = append(groupAges, groupAge)
		}
	}
	// begin to convert groupAges to clause and param
	for idx, groupAge := range groupAges {
		age_start_str := fmt.Sprintf("age_start_%d", idx)
		age_end_str := fmt.Sprintf("age_end_%d", idx)
		clauses = append(clauses, fmt.Sprintf(":%s <= %s AND %s < :%s", age_start_str, field, field, age_end_str))

		// group example is 4,5,6
		// baseAge used as age start of the group, from the example is 4
		// rangeAge use as how many age range on group, from example is range from 4 to 6, which is 3 year
		baseAge := groupAge[0]
		rangeAge := len(groupAge)
		// sample: current date is 13 Jan 2023, age group is 4,5,6, so we are looking date from
		// 14 Jan 2016 until 13 Jan 2019 -> range is 3 year
		ld := time.Now().AddDate(-(baseAge + rangeAge), 0, 1)
		param[age_start_str] = time.Date(ld.Year(), ld.Month(), ld.Day(), 0, 0, 0, 0, time.UTC)
		ld = time.Now().AddDate(-(baseAge), 0, 0)
		param[age_end_str] = time.Date(ld.Year(), ld.Month(), ld.Day(), 0, 0, 0, 0, time.UTC)
	}

	clause := strings.Join(clauses, " OR ")
	return clause, param
}

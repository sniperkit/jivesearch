package instant

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/language"
)

// Stats is an instant answer that
// returns the average, median, etc.
type Stats struct {
	Answer
}

var reStats *regexp.Regexp

func (s *Stats) setQuery(r *http.Request, qv string) answerer {
	s.Answer.setQuery(r, qv)
	return s
}

func (s *Stats) setUserAgent(r *http.Request) answerer {
	return s
}

func (s *Stats) setLanguage(lang language.Tag) answerer {
	s.language = lang
	return s
}

func (s *Stats) setType() answerer {
	s.Type = "stats"
	return s
}

func (s *Stats) setRegex() answerer {
	triggers := []string{
		"avg", "average", "mean", "median", "sum", "total",
	}

	t := strings.Join(triggers, "|")
	s.regex = append(s.regex, regexp.MustCompile(fmt.Sprintf(`^(?P<trigger>%s) (?P<remainder>.*)$`, t)))
	s.regex = append(s.regex, regexp.MustCompile(fmt.Sprintf(`^(?P<remainder>.*) (?P<trigger>%s)$`, t)))

	return s
}

func (s *Stats) solve(r *http.Request) answerer {
	// get all the numbers..this regexp will correctly grab e notation
	numbersStrings := reStats.FindAllString(s.remainder, -1)
	numbers := []float64{}

	for _, n := range numbersStrings {
		if i, err := strconv.ParseFloat(n, 64); err == nil {
			numbers = append(numbers, i)
		}
	}

	var txt string
	var ans float64

	switch s.triggerWord {
	case "avg", "average", "mean":
		txt = "Average: "
		ans = average(numbers)
	case "median":
		txt = "Median: "
		ans = median(numbers)
	case "sum", "total":
		txt = "Sum: "
		ans = sum(numbers)
	}

	s.Solution = txt + strconv.FormatFloat(ans, 'f', -1, 64)

	return s
}

func (s *Stats) setCache() answerer {
	s.Cache = true
	return s
}

func (s *Stats) tests() []test {
	typ := "stats"

	tests := []test{
		{
			query: "avg 3 4e6",
			expected: []Data{
				{
					Type:      typ,
					Triggered: true,
					Solution:  "Average: 2000001.5",
					Cache:     true,
				},
			},
		},
		{
			query: "11 18 -142 Average",
			expected: []Data{
				{
					Type:      typ,
					Triggered: true,
					Solution:  "Average: -37.666666666666664",
					Cache:     true,
				},
			},
		},
		{
			query: "6 3 -5 23 Median",
			expected: []Data{
				{
					Type:      typ,
					Triggered: true,
					Solution:  "Median: 4.5",
					Cache:     true,
				},
			},
		},
		{
			query: "median 17 12 -18",
			expected: []Data{
				{
					Type:      typ,
					Triggered: true,
					Solution:  "Median: 12",
					Cache:     true,
				},
			},
		},
		{
			query: "58 96 -41 sum",
			expected: []Data{
				{
					Type:      typ,
					Triggered: true,
					Solution:  "Sum: 113",
					Cache:     true,
				},
			},
		},
		{
			query: "Total -17 3 87 -476",
			expected: []Data{
				{
					Type:      typ,
					Triggered: true,
					Solution:  "Sum: -403",
					Cache:     true,
				},
			},
		},
	}

	return tests
}

func average(numbers []float64) float64 {
	total := sum(numbers)
	return total / float64(len(numbers))
}

func median(numbers []float64) float64 {
	sort.Float64s(numbers)
	middle := len(numbers) / 2
	result := numbers[middle]
	if len(numbers)%2 == 0 {
		result = (result + numbers[middle-1]) / 2
	}
	return result
}

func sum(numbers []float64) float64 {
	var total float64
	for _, value := range numbers {
		total += value
	}
	return total
}

func init() {
	reStats = regexp.MustCompile(`[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?`)
}

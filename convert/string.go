package convert

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func StringToTime(s string) *time.Time {
	t := new(time.Time)
	if s == "" {
		return t
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	*t = time.Unix(i, 0)

	return t
}

func StringToIntSlice(s string) []int {
	split := strings.Split(s, ",") // Or another delimiter if not ","
	ints := make([]int, 0, len(split))
	for _, str := range split {
		if i, err := strconv.Atoi(str); err == nil {
			ints = append(ints, i)
		}
	}
	return ints
}

func StringToBool(s string) bool {
	return s == "1"
}

func StringToInt(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return i
}

func StringToFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, fmt.Errorf("couldn't convert %q to float: %q", s, err)
	}
	return f, nil
}

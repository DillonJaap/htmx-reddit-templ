package convert

import "strconv"

func Int(str string) (int, error) {
	id, err := strconv.ParseInt(str, 10, 64)
	return int(id), err
}

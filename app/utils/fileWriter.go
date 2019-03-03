package utils

import (
	"io/ioutil"
)

/*
WriteFile func(filename string, content string)
*/
func WriteFile(filename string, content string) error {
	data := []byte(content)
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

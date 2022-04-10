package File

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func fileHandleErr(err *error) bool {
	if *err != nil {
		log.Println(*err)
		return true
	}

	return false
}

func Exists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func Read(filename string) (string, error) {
	file, err := os.Open(filename)

	if fileHandleErr(&err) {
		message := fmt.Sprintf("Cannot open file: %s", filename)
		err := errors.New(message)
		return "", err
	}

	defer file.Close()

	info, err := file.Stat()

	if fileHandleErr(&err) {
		message := fmt.Sprintf("Cannot get file size: %s", filename)
		err := errors.New(message)
		return "", err
	}

	size := info.Size()

	buffer := make([]byte, size)

	_, err = file.Read(buffer)

	fileHandleErr(&err)

	return string(buffer), nil
}

func Write(filename string, contents string) error {
	is := Exists(filename)

	if !is {
		Create(filename)
	}

	file, err := os.OpenFile(filename, os.O_RDWR, os.ModeAppend)

	if fileHandleErr(&err) {
		message := fmt.Sprintf("Cannot open file: %s", filename)
		err := errors.New(message)
		return err
	}

	defer file.Close()

	buffer := []byte(contents)

	_, err = file.Write(buffer)

	if err != nil {
		return err
	}

	return nil
}

func Create(filename string) error {
	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func Delete(filename string) error {
	err := os.Remove(filename)

	if fileHandleErr(&err) {
		message := fmt.Sprintf("Cannot delete file: %s", filename)
		err := errors.New(message)
		return err
	}

	return nil
}

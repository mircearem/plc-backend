package Shm

import (
	"errors"
	"plc-backend/File"
)

func Write(file string, lockfile string, contents string) error {
	is := File.Exists(lockfile)

	if is {
		err := errors.New("ERR: Previous read operation not finished")
		return err
	}

	// Write data to shared memory file
	err := File.Write(file, contents)

	if err != nil {
		return err
	}

	// Create lock file
	err = File.Create(lockfile)

	if err != nil {
		return err
	}

	return nil
}

func Read(file string, lockFile string) (string, error) {
	// Check codesys.write.done
	var sharedMemoryContents string

	is := File.Exists(lockFile)

	if !is {
		err := errors.New("ERR: Writing not finished")
		return "", err
	}

	// Read codesys.json
	sharedMemoryContents, err := File.Read(file)

	if err != nil {
		return "", err
	}

	// Delete codesys.write.done
	err = File.Delete(lockFile)

	if err != nil {
		return "", err
	}

	return sharedMemoryContents, nil
}

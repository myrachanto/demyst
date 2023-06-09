package support

import "os"

func RenameImage(oldName, newName string) string {
	err := os.Rename(oldName, newName)
	if err != nil {
		panic(err)
	}
	return newName
}

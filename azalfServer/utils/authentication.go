/*package utils

import (
	"io/ioutil"
	"log"
	"os"
)

func authenticate(hash string) bool {
	allowed := false

	// load the encrypted key from the /home/user/.local/azalf/auth/auth.lck file
	// or in the case of windows, from the %APPDATA%/azalf/auth/auth.lck file
	// decrypt the key and compare it to the hash after unhashing it
	// if they match, return true
	// if they don't match, return false

	// load the key from the file
	if _, err := os.Stat(authFile); os.IsNotExist(err) {
		// read the key
		key, err := ioutil.ReadFile(authFile)
		if err != nil {
			log.Fatal(err)
		}
		// decrypt the key
	}
}

func encrypt(pass string) string {
	// using the passphrase, and the name of the current user, encrypt the passphrase
	// and return the key
	key := pass + os.Getenv("USER")

	// encrypt the key
	// encrypt using SHA256
	// then return the encrypted key
	var encryptedKey string

}
*/
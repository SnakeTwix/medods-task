package env

import (
	"log"
	"os"
)

func Get(envName Type) string {
	envValue, ok := os.LookupEnv(string(envName))
	if !ok {
		log.Fatalf("%s is not set \n", envName)
	}

	return envValue
}

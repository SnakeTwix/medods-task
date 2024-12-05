package env

import (
	"log"
	"os"
)

func Get(envName string) string {
	envValue, ok := os.LookupEnv(envName)
	if !ok {
		log.Fatalf("%s is not set \n", envName)
	}

	return envValue
}

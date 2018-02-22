package helpers

import "os"

//FetchEnv is fetching env variables.
//If run with "false" argument, non existing key will not cause panic.
func FetchEnv(key string, optionalPanicArgs ...bool) string {
	shouldPanic := true
	if len(optionalPanicArgs) > 0 && !optionalPanicArgs[0] {
		shouldPanic = false
	}

	value := os.Getenv(key)
	if value == "" && shouldPanic {
		panic("No env with key: " + key)
	}

	return value
}

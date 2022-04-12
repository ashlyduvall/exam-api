package main

import "os"

// General env var lookup with default function
func GetConfig(key string, default_value string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	} else {
		return default_value
	}
}

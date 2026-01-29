package utils

import "fmt"

// Greet returns a greeting message
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s! Welcome to home-ctrl!", name)
}
// Package main is the entry point for the spot CLI tool.
package main

import "github.com/crnvl96/spot/internal"

// main runs the spot CLI application.
//
// spot is a cli tool to check if you have uncommited changes in any of your repos

func main() {
	internal.Execute()
}

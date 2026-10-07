package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/HaeengIn/Image-Tweaks-Go/internal/multiple"
	"github.com/HaeengIn/Image-Tweaks-Go/internal/single"
)

var reader = bufio.NewReader(os.Stdin)

// Clear console
func clearConsole() {
	fmt.Print("\033[H\033[2J")
}

// Get input from user
func getInput() string {
	fmt.Print("> ")
	input, _ := reader.ReadString('\n')

	return strings.TrimSpace(input)
}

// select which process to do
func selectProcess() string {
	fmt.Println("[ Optimize / Convert ]")
	for {
		input := strings.ToLower(getInput())
		if input == "optimize" || input == "convert" {
			return input
		}
		fmt.Println("\nInvalid Input.")
	}
}

// Run single process
func runSingle() {
	clearConsole()

	fmt.Println("[ Single Process ]")
	fmt.Println("Select process from the menu below.")
	process := selectProcess()

	clearConsole()

	for {
		switch process {
		case "optimize":
			single.RunOptimize()
		case "convert":
			single.RunConvert()
		}

		// Check for another process
		fmt.Println("\nDo you want to process another image? (Y/N)")

		for {
			answer := strings.ToLower(getInput())
			if answer == "y" {
				clearConsole()
				break
			} else if answer == "n" {
				clearConsole()
				main()
				return
			}
			fmt.Println("\nInvalid Input. Please enter Y or N.")
		}
	}
}

// Run multiple process
func runMultiple() {
	clearConsole()

	fmt.Println("[ Multiple Process ]")
	fmt.Println("Select process from the menu below.")
	process := selectProcess()

	clearConsole()

	for {
		switch process {
		case "optimize":
			multiple.RunOptimize()
		case "convert":
			multiple.RunConvert()
		}

		// Check for another process
		fmt.Println("\nDo you want to process another folder? (Y/N)")

		for {
			answer := strings.ToLower(getInput())
			if answer == "y" {
				clearConsole()
				break
			} else if answer == "n" {
				clearConsole()
				main()
				return
			}
			fmt.Println("\nInvalid Input. Please enter Y or N.")
		}
	}
}

// Print program information
func openInfo() {
	fmt.Println("[ Image Tweaks for Go ]\n")
	fmt.Println("Currently supported process:")
	fmt.Println("[ Optimize, Convert ]\n")
	fmt.Println("Supported formats:")
	fmt.Println("[ JPG, JPEG, PNG, WEBP ]\n")
	fmt.Println("Version: Alpha 1.0.0")
	fmt.Println("\nPress Enter to return to the main menu...")

	reader.ReadString('\n')
	clearConsole()

	main()
}

// Exit app
func exitApp() {
	fmt.Println("Do you want to really exit the app? (Y/N)")

	for {
		answer := strings.ToLower(getInput())
		if answer == "y" {
			clearConsole()
			os.Exit(0)
		} else if answer == "n" {
			clearConsole()
			main()
			return
		}
		fmt.Println("\nInvalid Input. Please enter Y or N.")
	}
}

// Main app
func main() {
	clearConsole()

	fmt.Println("[ Image Tweaks for Go ]")
	fmt.Println("Type and enter the menu option what you want to do.")
	fmt.Println("[ Single / Multiple / Info / Exit ]")

	for {
		input := strings.ToLower(getInput())
		switch input {
		case "single":
			runSingle()
			return
		case "multiple":
			runMultiple()
			return
		case "info":
			clearConsole()
			openInfo()
			return
		case "exit":
			clearConsole()
			exitApp()
			return
		default:
			fmt.Println("\nInvalid Input.")
		}
	}
}

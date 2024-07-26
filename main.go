package main

import (
	"bufio"     // Package bufio implements buffered Input/Output.
	"fmt"       // Package fmt implements formatted I/O with functions analogous to C's printf and scanf.
	"math/rand" // Package rand implements pseudo-random number generators.
	"os"        // Package os provides a platform-independent interface to operating system functionality.
	"strconv"   // Package strconv implements conversions to and from string representations of basic data types.
	"strings"   // Package strings implements simple functions to manipulate UTF-8 encoded strings.
)

const maxPlayerNumbers byte = 6 // The maximum number of numbers a player can select
const maxDrawNumbers byte = 7   // 6 + bonus
const lowestNumber byte = 1     // The lowest number in the range
const highestNumber byte = 59   // The highest number in the range

type GameState int // custom type for enum

const ( // Enum values
	home GameState = iota
	luckyDip
	pickNumbers
	results
)

var gameState GameState // Enum variable

var prizePots = [...]float32{25, 100, 1000, 10000, 100000} // Prize pots for 3, 4, 5, 6, 6+bonus
var playerNumbers []byte                                   // Numbers which the player has selected or been assigned
var remainingNumbers [59]byte                              // Pool of numbers that contain only numbers that have not been selected or drawn

func main() {
	gameState = home // Set the initial game state to home

	for /*ever*/ {
		switch gameState {
		case home:
			Home()
		case luckyDip:
			LuckyDip()
		case pickNumbers:
			PickNumbers()
		case results:
			Results()
		}
	}
}

func Home() {
	fmt.Println("Welcome to the lottery game\n")
	fmt.Println("1.\tLucky Dip\n")
	fmt.Println("2.\tPick Numbers\n")

	key := ReadInput() // Read the user's input

	switch key {
	case "1":
		gameState = luckyDip
	case "2":
		gameState = pickNumbers
	default:
		fmt.Printf("%s Invalid option\n\n", key)
		Home()
	}
}

func DrawNumbers() []byte {
	PopulateRemainingNumbers()   // Populate the remainingNumbers array with numbers 1-59
	slice := remainingNumbers[:] // Create a slice from the remainingNumbers array (take them all)

	var arr []byte
	for i := 0; i < int(maxDrawNumbers); i++ {
		randIdx := rand.Intn(len(slice))      // Generate a random index
		winningNumber := slice[randIdx]       // Get the value at the random index
		arr = append(arr, winningNumber)      // Add the value to the playerNumbers slice
		slice = deleteElement(slice, randIdx) // Delete the element at the random index
	}

	return arr
}

func LuckyDip() {
	playerNumbers = DrawNumbers()

	fmt.Println("Your numbers are: ", playerNumbers)
	fmt.Println("Would you like to keep these numbers? (Y/N)")

	key := ReadInput()

	if strings.ToUpper(key) == "Y" {
		gameState = results
	} else {
		playerNumbers = nil
	}
}

func PickNumbers() {
	PopulateRemainingNumbers()
	clearConsole()

	for byte(len(playerNumbers)) < maxPlayerNumbers+1 {
		fmt.Print("Please select a numbers between 1 and 59: ")
		val := ReadInput()           // Read the user's input
		iVal, _ := strconv.Atoi(val) // Convert the user's input to an integer

		if iVal < 1 || iVal > 59 {
			fmt.Println("Invalid number, please try again")
		} else if !(contains(playerNumbers, byte(iVal))) {
			playerNumbers = append(playerNumbers, byte(iVal))
		} else {
			fmt.Println("You have already selected that number")
		}
	}

	fmt.Println("Your numbers are: ", playerNumbers)

	gameState = results
}

func Results() {
	clearConsole()

	winningNumbers := DrawNumbers()

	fmt.Println("Your numbers are: ", playerNumbers)
	fmt.Println("The winning numbers are: ", winningNumbers)

	matches := 0

	for _, playerNumber := range playerNumbers { // discard index keep value
		if contains(winningNumbers, playerNumber) { // if the player number is in the winning numbers
			matches++ // increment the matches
		}
	}

	if matches == 6 {
		fmt.Println("Congratulations! You have won the jackpot!")
	} else {
		fmt.Printf("You matched %d numbers\n", matches)
	}

	fmt.Println("Would you like to play again? (Y/N)")

	key := ReadInput()

	if strings.ToUpper(key) == "Y" {
		playerNumbers = nil
		gameState = home
	} else {
		os.Exit(0)
	}
}

func PopulateRemainingNumbers() {
	for i := lowestNumber; i <= highestNumber; i++ {
		remainingNumbers[i-1] = i
	}
}

func ReadInput() string {
	// Initiate user input reader
	reader := bufio.NewReader(os.Stdin)

	// Call the reader to read user's input
	key, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return key[:len(key)-1]
} //https://www.mizouzie.dev/articles/3-ways-to-read-input-with-go-cli/

func deleteElement(slice []byte, index int) []byte {
	return append(slice[:index], slice[index+1:]...)
} // https://www.tutorialspoint.com/delete-elements-in-a-slice-in-golang#:~:text=In%20this%20example%2C%20we%20create,we%20assign%20back%20to%20slice.

func clearConsole() {
	fmt.Print("\033[H\033[2J")
} // Copilot

func contains(s []byte, e byte) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
} //https://stackoverflow.com/questions/10485743/contains-method-for-a-slice

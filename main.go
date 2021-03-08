package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	handSize  int = 4
	maxPower  int = 12
	maxHealth int = 12
	maxRank   int = 9
)

type card struct {
	value    int
	name     string
	playable bool
}
type hand struct {
	health int
	power  int
	cards  []card
}

func initHand() hand {
	var tmpHand hand
	var tmpCard card
	tmpHand.health = maxHealth
	tmpHand.power = maxPower
	tmpHand.cards = make([]card, handSize)
	for i := range tmpHand.cards {
		tmpCard.name = "Card " + strconv.Itoa(i)
		tmpCard.playable = true
		tmpCard.value = rand.Intn(maxRank) + 1
		tmpHand.cards[i] = tmpCard
	}

	return tmpHand
}

func drawHand(hand hand) {
	fmt.Println("--------------------------")
	stringHand := "| \033[34m" + strconv.Itoa(hand.health) + "\033[0m |"
	for _, v := range hand.cards {
		if v.playable {
			stringHand += "  \033[32m" + strconv.Itoa(v.value) + "\033[0m"
		} else {
			stringHand += "  " + strconv.Itoa(v.value)
		}

	}
	stringHand += "  | \033[31m" + strconv.Itoa(hand.power) + "\033[0m |"
	fmt.Println(stringHand)
	fmt.Println("--------------------------")
}

func drawTable(firstHand hand, secondHand hand) {
	drawHand(firstHand)
	fmt.Println("++++++++++++++++++++++++++")
	drawHand(secondHand)

}

func main() {

	rand.Seed(time.Now().UnixNano())

	userHand := initHand()
	compHand := initHand()

	//Bool value to select current player
	currentPlayer := rand.Intn(2) == 0

	scanner := bufio.NewScanner(os.Stdin)

	var cardNumber, cardPower int
	var err error

	for {
		drawTable(userHand, compHand)
		if currentPlayer {
			for {
				fmt.Println("Enter card number: ")
				scanner.Scan()
				cardNumber, err = strconv.Atoi(scanner.Text())
				if err != nil {
					fmt.Println("Unrecognized character")
					continue
				} else {
					if cardNumber > handSize || cardNumber < 1 {
						fmt.Println("Card number is too big. Card number range is 1 ..", handSize)
						continue
					}
					if !userHand.cards[cardNumber-1].playable {
						fmt.Println("Card already played. Please choose another")
						continue
					}
				}
				userHand.cards[cardNumber-1].playable = false
				break
			}
			fmt.Println("Enter power: ")
			scanner.Scan()
			cardPower, err = strconv.Atoi(scanner.Text())
			fmt.Printf("Playing card %d with power %d\n", cardNumber, cardPower)
		} else {
			fmt.Println("My turn!")
		}
		currentPlayer = !currentPlayer
	}
	//Loop User Hand
	//askForCard()
	//askForPower()

	//Loop Comp Hand
	//calcBestPlay()

	//Loop until run out of cards or out of HP

	/* 	total := 0
	   	for a := 0; a <= 12; a++ {
	   		for b := 0; b <= 12-a; b++ {
	   			for c := 0; c <= 12-a-b; c++ {
	   				d := 12 - a - b - c
	   				fmt.Println(a, b, c, d)
	   				total++
	   			}
	   		}
	   	}
	   	fmt.Println(total)
	*/
}

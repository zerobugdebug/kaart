package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	health        int
	power         int
	selectedCard  int
	selectedPower int
	cards         []card
}

func initHand() hand {
	var tmpHand hand
	var tmpCard card
	tmpHand.health = maxHealth
	tmpHand.power = maxPower
	tmpHand.selectedCard = -1
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
	fmt.Println("┌────┬" + strings.Repeat("─", 15) + "┬────┐")
	fmt.Printf("│\033[34;1m%3d\033[0m │", hand.health)
	for _, v := range hand.cards {
		if v.playable {
			fmt.Printf("\033[32m%2d\033[0m │", v.value)
		} else {
			fmt.Printf("\033[37m%2d\033[0m │", v.value)
		}

	}
	fmt.Printf("\033[31m%3d\033[0m │\n", hand.power)
	fmt.Println("└────┴" + strings.Repeat("─", 15) + "┴────┘")
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func drawTable(firstHand hand, secondHand hand) {
	//clearScreen()
	fmt.Print("\033[H\033[2J")

	drawHand(firstHand)

	fmt.Println("╔" + strings.Repeat("═", 25) + "╗")
	fmt.Println("║" + strings.Repeat(" ", 25) + "║")

	if firstHand.selectedCard != -1 {
		fmt.Printf("║\033[32m%4d\033[0m+\033[31m%-20d\033[0m║\n", firstHand.cards[firstHand.selectedCard].value, firstHand.selectedPower)
	} else {
		fmt.Println("║" + strings.Repeat(" ", 25) + "║")
	}

	fmt.Println("║" + strings.Repeat(" ", 25) + "║")

	if secondHand.selectedCard != -1 {

		fmt.Printf("║\033[32m%20d\033[0m+\033[31m%-4d\033[0m║\n", secondHand.cards[secondHand.selectedCard].value, secondHand.selectedPower)

		//selectedCards += "\033[32m" + strconv.Itoa(secondHand.cards[secondHand.selectedCard].value) + "\033[0m+"
		//selectedCards += "\033[31m" + strconv.Itoa(secondHand.selectedPower) + "\033[0m    ++"
	} else {
		fmt.Println("║" + strings.Repeat(" ", 25) + "║")
	}

	fmt.Println("║" + strings.Repeat(" ", 25) + "║")
	fmt.Println("╚" + strings.Repeat("═", 25) + "╝")

	drawHand(secondHand)

}

func drawBattle(firstHand hand, secondHand hand) {

	compTotalPower := firstHand.cards[firstHand.selectedCard].value + firstHand.selectedPower
	userTotalPower := secondHand.cards[secondHand.selectedCard].value + secondHand.selectedPower

	fmt.Print("\033[H\033[2J")

	drawHand(firstHand)

	fmt.Println("╔" + strings.Repeat("═", 25) + "╗")
	fmt.Println("║" + strings.Repeat(" ", 25) + "║")
	fmt.Printf("║\033[32m%11d\033[0m vs \033[32m%-10d\033[0m║\n", compTotalPower, userTotalPower)
	if userTotalPower > compTotalPower {
		fmt.Println("║" + strings.Repeat(" ", 8) + "\033[32mUSER WINS\033[0m" + strings.Repeat(" ", 8) + "║")
		fmt.Printf("║"+strings.Repeat(" ", 8)+"\033[32mDAMAGE:%3d\033[0m"+strings.Repeat(" ", 7)+"║\n", userTotalPower-compTotalPower)
	} else if userTotalPower < compTotalPower {
		fmt.Println("║" + strings.Repeat(" ", 8) + "\033[31mCOMP WINS\033[0m" + strings.Repeat(" ", 8) + "║")
		fmt.Printf("║"+strings.Repeat(" ", 8)+"\033[31mDAMAGE:%3d\033[0m"+strings.Repeat(" ", 7)+"║\n", compTotalPower-userTotalPower)
	} else {
		fmt.Println("║" + strings.Repeat(" ", 25) + "║")
		fmt.Println("║" + strings.Repeat(" ", 11) + "\033[31mDRAW\033[0m" + strings.Repeat(" ", 10) + "║")
	}

	fmt.Println("║" + strings.Repeat(" ", 25) + "║")
	fmt.Println("╚" + strings.Repeat("═", 25) + "╝")

	drawHand(secondHand)

}

func processUserTurn(userHand hand) hand {
	var cardNumber, cardPower int
	var err error

	scanner := bufio.NewScanner(os.Stdin)

	userHand.selectedCard = -1
	for {
		fmt.Print("Enter card number: ")
		scanner.Scan()
		cardNumber, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Unrecognized character")
			continue
		} else {
			if cardNumber > handSize || cardNumber < 1 {
				fmt.Println("Incorrect card number. Card number range is 1 ..", handSize)
				continue
			}
			if !userHand.cards[cardNumber-1].playable {
				fmt.Println("Card already played. Please choose another")
				continue
			}
		}
		userHand.cards[cardNumber-1].playable = false
		userHand.selectedCard = cardNumber - 1
		break
	}
	if userHand.power > 0 {
		for {
			fmt.Print("Enter power: ")
			scanner.Scan()
			cardPower, err = strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Unrecognized character")
				continue
			} else {
				if cardPower > userHand.power || cardPower < 0 {
					fmt.Println("Incorrect power value. Power value range is 0 ..", userHand.power)
					continue
				}
			}
			break
		}
	} else {
		cardPower = 0
	}
	userHand.selectedPower = cardPower
	userHand.power -= cardPower
	return userHand
}

func processCompTurn(compHand hand) hand {
	var cardNumber, cardPower int
	compHand.selectedCard = -1
	for {
		cardNumber = rand.Intn(handSize)
		if compHand.cards[cardNumber].playable {
			break
		}
	}
	cardPower = rand.Intn(compHand.power + 1)
	compHand.cards[cardNumber].playable = false
	compHand.selectedCard = cardNumber
	compHand.selectedPower = cardPower
	compHand.power -= cardPower
	return compHand
}

func main() {

	rand.Seed(time.Now().UnixNano())

	userHand := initHand()
	compHand := initHand()

	//Bool value to select current player
	isUserTurn := rand.Intn(2) == 0
	//isTableDirty := true
	var totalUserPower, totalCompPower, totalUserDamage int
	//var cardNumber, cardPower int
	//var err error

	for i := 0; i < handSize; i++ {
		userHand.selectedCard = -1
		compHand.selectedCard = -1
		drawTable(compHand, userHand)
		if isUserTurn {
			fmt.Println("USER TURN")
			userHand = processUserTurn(userHand)
			drawTable(compHand, userHand)
			compHand = processCompTurn(compHand)
		} else {
			fmt.Println("COMP TURN")
			compHand = processCompTurn(compHand)
			drawTable(compHand, userHand)
			userHand = processUserTurn(userHand)
		}
		drawTable(compHand, userHand)
		fmt.Print("Press 'Enter' for the turn results...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		drawBattle(compHand, userHand)
		totalUserPower = userHand.cards[userHand.selectedCard].value + userHand.selectedPower
		totalCompPower = compHand.cards[compHand.selectedCard].value + compHand.selectedPower
		totalUserDamage = totalUserPower - totalCompPower
		if totalUserDamage > 0 {
			compHand.health -= totalUserDamage
		} else {
			userHand.health += totalUserDamage
		}

		if userHand.health < 1 || compHand.health < 1 {
			break
		}
		if i < handSize-1 {
			fmt.Print("Press 'Enter' for the next turn...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}

		isUserTurn = !isUserTurn
	}

	fmt.Print("Press 'Enter' for the game results...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	userHand.selectedCard = -1
	compHand.selectedCard = -1
	drawTable(compHand, userHand)
	if userHand.health == compHand.health {
		fmt.Print("DRAW")
	} else if userHand.health < compHand.health {
		fmt.Print("COMP WINS")
	} else {
		fmt.Print("USER WINS")
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

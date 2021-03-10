package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zerobugdebug/go-log"
)

//Game constants
const (
	//HandSize is amount of cards every player should have
	HandSize int = 4
	//MaxPower is the maximum power available to the player
	MaxPower  int = 12
	maxHealth int = 12
	maxRank   int = 8
	minRank   int = 2
	maxDamage int = 8
	minDamage int = 1
)

//Color constants
const (
	clrReset string = "\033[0m"
	//clrPlayableCard    string = "\033[40;1m\033[37;1m"
	//clrNonPlayableCard string = "\033[40m\033[30;1m"
	clrSelectedCard          string = "\033[40;1m\033[38;5;10m"
	clrPlayableCard          string = "\033[40;1m\033[38;5;15m"
	clrNonPlayableCard       string = "\033[40m\033[38;5;239m"
	clrPlayableCardDamage    string = "\033[31m"
	clrPlayableCardPower     string = "\033[35m"
	clrNonPlayableCardDamage string = "\033[38;5;52m"
	clrNonPlayableCardPower  string = "\033[38;5;53m"
	clrHealth                string = "\033[34;1m\033[1m"
	clrPower                 string = "\033[35;1m\033[1m"
	clrGoodMessage           string = "\033[32m"
	clrBadMessage            string = "\033[31m"
)

type card struct {
	value    int
	damage   int
	name     string
	playable bool
}

//Hand is a struct to store player or comp hand
type Hand struct {
	health        int
	power         int
	selectedCard  int
	selectedPower int
	active        bool
	cards         []card
}

//Logger is a default log adapter
var Logger = log.New(os.Stdout).WithoutDebug()

func initHand() Hand {
	var tmpHand Hand
	var tmpCard card
	var totalValue int
	maxTotalValue := (HandSize-1)*maxRank + 1
	tmpHand.health = maxHealth
	tmpHand.power = MaxPower
	tmpHand.selectedCard = -1
	tmpHand.cards = make([]card, HandSize-1)
	for i := range tmpHand.cards {
		tmpCard.name = "Card " + strconv.Itoa(i)
		tmpCard.playable = true
		tmpCard.value = rand.Intn(maxRank-minRank+1) + minRank
		tmpCard.damage = rand.Intn(maxDamage-tmpCard.value+1) + minDamage
		totalValue += tmpCard.value
		tmpHand.cards[i] = tmpCard
	}
	//Add last card with remaining power to balance total power on both players
	tmpCard.name = "Card 3"
	tmpCard.playable = true
	tmpCard.value = (maxTotalValue - totalValue) + rand.Intn(minRank)
	if tmpCard.value > maxRank {
		tmpCard.value = maxRank
	}
	tmpCard.damage = rand.Intn(maxDamage-tmpCard.value+1) + minDamage
	tmpHand.cards = append(tmpHand.cards, tmpCard)
	return tmpHand
}

func drawHand(hand Hand) {
	//fmt.Println("┌────┬" + strings.Repeat("─", 15) + "┬────┐")

	fmt.Println("┌────────────────────────────┐")
	fmt.Print("│      ")
	for i, v := range hand.cards {
		if hand.selectedCard == i {
			fmt.Print(clrSelectedCard + "┌──┐" + clrReset)
		} else if v.playable {
			fmt.Print(clrPlayableCard + "┌──┐" + clrReset)
		} else {
			fmt.Print(clrNonPlayableCard + "┌──┐" + clrReset)
		}
	}
	fmt.Println("      │")

	fmt.Print("│      ")
	for i, v := range hand.cards {
		if hand.selectedCard == i {
			fmt.Print(clrSelectedCard + "│" + clrReset)
			fmt.Printf(clrSelectedCard+clrPlayableCardPower+"%2d"+clrReset, v.value)
			fmt.Print(clrSelectedCard + "│" + clrReset)
		} else if v.playable {
			fmt.Print(clrPlayableCard + "│" + clrReset)
			fmt.Printf(clrPlayableCard+clrPlayableCardPower+"%2d"+clrReset, v.value)
			fmt.Print(clrPlayableCard + "│" + clrReset)
		} else {
			fmt.Print(clrNonPlayableCard + "│" + clrReset)
			fmt.Printf(clrNonPlayableCard+clrNonPlayableCardPower+"%2d"+clrReset, v.value)
			fmt.Print(clrNonPlayableCard + "│" + clrReset)
		}
	}
	fmt.Println("      │")

	fmt.Printf("│  "+clrHealth+"%2d  "+clrReset, hand.health)
	for i, v := range hand.cards {
		if hand.selectedCard == i {
			fmt.Print(clrSelectedCard + "├──┤" + clrReset)
		} else if v.playable {
			fmt.Print(clrPlayableCard + "├──┤" + clrReset)
		} else {
			fmt.Print(clrNonPlayableCard + "├──┤" + clrReset)
		}
	}
	fmt.Printf("  "+clrPower+"%2d"+clrReset+"  │\n", hand.power)

	fmt.Print("│      ")
	for i, v := range hand.cards {
		if hand.selectedCard == i {
			fmt.Print(clrSelectedCard + "│" + clrReset)
			fmt.Printf(clrSelectedCard+clrPlayableCardDamage+"%2d"+clrReset, v.damage)
			fmt.Print(clrSelectedCard + "│" + clrReset)
		} else if v.playable {
			fmt.Print(clrPlayableCard + "│" + clrReset)
			fmt.Printf(clrPlayableCard+clrPlayableCardDamage+"%2d"+clrReset, v.damage)
			fmt.Print(clrPlayableCard + "│" + clrReset)
		} else {
			fmt.Print(clrNonPlayableCard + "│" + clrReset)
			fmt.Printf(clrNonPlayableCard+clrNonPlayableCardDamage+"%2d"+clrReset, v.damage)
			fmt.Print(clrNonPlayableCard + "│" + clrReset)
		}
	}
	fmt.Println("      │")

	fmt.Print("│      ")
	for i, v := range hand.cards {
		if hand.selectedCard == i {
			fmt.Print(clrSelectedCard + "└──┘" + clrReset)
		} else if v.playable {
			fmt.Print(clrPlayableCard + "└──┘" + clrReset)
		} else {
			fmt.Print(clrNonPlayableCard + "└──┘" + clrReset)
		}
	}
	fmt.Println("      │")
	fmt.Println("└────────────────────────────┘")
	/*
		fmt.Println("┌────┬───────────────┬────┐")
		for _, v := range hand.cards {
			if v.playable {
				fmt.Printf("\033[32m%2d\033[0m │", v.value)
			} else {
				fmt.Printf("\033[37m%2d\033[0m │", v.value)
			}
		}
		fmt.Printf("│"+clrHealth+"%3d"+clrReset+" │", hand.health)
		//fmt.Print("│", clrHealth, hand.health, clrReset, " │")
		fmt.Printf(clrPower+"%3d"+clrReset+" │\n", hand.power)
		fmt.Println("└────┴" + strings.Repeat("─", 15) + "┴────┘") */
}

func drawTable(firstHand Hand, secondHand Hand) {
	//clearScreen()
	fmt.Print("\033[H\033[2J")

	drawHand(firstHand)

	//fmt.Println("╔" + strings.Repeat("═", 25) + "╗")
	//fmt.Println("║" + strings.Repeat(" ", 25) + "║")
	fmt.Println("╔════════════════════════════╗")
	fmt.Println("║                            ║")

	if firstHand.selectedCard != -1 {
		//tmpString := fmt.Sprint(firstHand.cards[firstHand.selectedCard].value, "+", firstHand.cards[firstHand.selectedCard].value, "*", firstHand.selectedPower)
		tmpString := fmt.Sprint(firstHand.cards[firstHand.selectedCard].value, "+", firstHand.cards[firstHand.selectedCard].value, "*??")
		fmt.Printf("║ "+clrPlayableCardPower+"%v"+clrReset, tmpString)

		fmt.Println(strings.Repeat(" ", 27-len(tmpString)) + "║")
		//fmt.Printf("║\033[32m%4d\033[0m+\033[31m%-20v\033[0m║\n", firstHand.cards[firstHand.selectedCard].value, "?")
	} else {
		fmt.Println("║                            ║")
	}

	fmt.Println("║                            ║")

	if secondHand.selectedCard != -1 {
		tmpString := fmt.Sprint(secondHand.cards[secondHand.selectedCard].value, "+", secondHand.cards[secondHand.selectedCard].value, "*", secondHand.selectedPower)
		fmt.Printf("║"+clrPlayableCardPower+"%27v"+clrReset, tmpString)
		fmt.Println(" ║")

		//fmt.Printf("║"+clrCardPower+"%20d"+clrReset+"+"+clrCardPower+"%-4d"+clrReset+"║\n", secondHand.cards[secondHand.selectedCard].value, secondHand.selectedPower)
		//fmt.Printf("║\033[32m%20d\033[0m+\033[31m%-4d\033[0m║\n", secondHand.cards[secondHand.selectedCard].value, secondHand.selectedPower)

		//selectedCards += "\033[32m" + strconv.Itoa(secondHand.cards[secondHand.selectedCard].value) + "\033[0m+"
		//selectedCards += "\033[31m" + strconv.Itoa(secondHand.selectedPower) + "\033[0m    ++"
	} else {
		fmt.Println("║                            ║")
	}

	fmt.Println("║                            ║")
	fmt.Println("╚════════════════════════════╝")

	drawHand(secondHand)

}

func drawBattle(firstHand Hand, secondHand Hand) {

	compTotalPower := firstHand.cards[firstHand.selectedCard].value + firstHand.selectedPower*firstHand.cards[firstHand.selectedCard].value
	compTotalPowerString := fmt.Sprint(firstHand.cards[firstHand.selectedCard].value, "+", firstHand.cards[firstHand.selectedCard].value, "*", firstHand.selectedPower, "=", compTotalPower)
	userTotalPower := secondHand.cards[secondHand.selectedCard].value + secondHand.selectedPower*secondHand.cards[secondHand.selectedCard].value
	userTotalPowerString := fmt.Sprint(secondHand.cards[secondHand.selectedCard].value, "+", secondHand.cards[secondHand.selectedCard].value, "*", secondHand.selectedPower, "=", userTotalPower)
	totalLen := len(compTotalPowerString) + len(userTotalPowerString) + 4
	fmt.Print("\033[H\033[2J")

	drawHand(firstHand)

	fmt.Println("╔════════════════════════════╗")
	fmt.Println("║                            ║")
	fmt.Print("║" + strings.Repeat(" ", (28-totalLen)/2))
	fmt.Print(clrPlayableCardPower + compTotalPowerString + clrReset + " vs " + clrPlayableCardPower + userTotalPowerString + clrReset)
	fmt.Println(strings.Repeat(" ", (29-totalLen)/2) + "║")
	if userTotalPower > compTotalPower {
		fmt.Println("║" + strings.Repeat(" ", 9) + clrGoodMessage + "USER  WINS" + clrReset + strings.Repeat(" ", 9) + "║")
		fmt.Printf("║"+strings.Repeat(" ", 9)+clrGoodMessage+"DAMAGE:%3d"+clrReset+strings.Repeat(" ", 9)+"║\n", secondHand.cards[secondHand.selectedCard].damage)
	} else if userTotalPower < compTotalPower {
		fmt.Println("║" + strings.Repeat(" ", 9) + clrBadMessage + "COMP  WINS" + clrReset + strings.Repeat(" ", 9) + "║")
		fmt.Printf("║"+strings.Repeat(" ", 9)+clrBadMessage+"DAMAGE:%3d"+clrReset+strings.Repeat(" ", 9)+"║\n", firstHand.cards[firstHand.selectedCard].damage)
	} else {
		fmt.Println("║                            ║")
		fmt.Println("║" + strings.Repeat(" ", 12) + clrGoodMessage + "DRAW" + clrReset + strings.Repeat(" ", 12) + "║")
	}

	fmt.Println("║                            ║")
	fmt.Println("╚════════════════════════════╝")

	drawHand(secondHand)

}

func processUserTurn(userHand Hand) Hand {
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
			if cardNumber > HandSize || cardNumber < 1 {
				fmt.Println("Incorrect card number. Card number range is 1 ..", HandSize)
				continue
			}
			if !userHand.cards[cardNumber-1].playable {
				fmt.Println("Card already played. Please choose another")
				continue
			}
		}
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
	return userHand
}

func processCompTurn(compHand Hand, userHand Hand) Hand {
	var cardNumber, cardPower int
	var playableCards int
	for i, v := range compHand.cards {
		if v.playable {
			playableCards++
			cardNumber = i
			cardPower = compHand.power
		}
	}
	if playableCards > 1 {
		cardNumber, cardPower = GetNextMove(compHand, userHand)
	}
	//	compHand.selectedCard = -1
	/* 	for {
	   		cardNumber = rand.Intn(HandSize)
	   		if compHand.cards[cardNumber].playable {
	   			break
	   		}
	   	}
	   	cardPower = rand.Intn(compHand.power + 1)
	*/
	compHand.selectedCard = cardNumber
	compHand.selectedPower = cardPower
	return compHand
}

func main() {

	rand.Seed(time.Now().UnixNano())

	userHand := initHand()
	compHand := initHand()
	Logger.Debug(compHand)
	Logger.Debug(userHand)

	Logger.Debug(compHand)
	Logger.Debug(userHand)
	//Logger.Fatal("!TEMP END!")
	//Bool value to select current player
	isUserTurn := rand.Intn(2) == 0
	//isTableDirty := true
	var totalUserPower, totalCompPower int
	//var cardNumber, cardPower int
	//var err error

	for i := 0; i < HandSize; i++ {
		userHand.selectedCard = -1
		compHand.selectedCard = -1
		compHand.active = !isUserTurn
		userHand.active = isUserTurn
		drawTable(compHand, userHand)
		if isUserTurn {
			fmt.Println("USER TURN")
			userHand = processUserTurn(userHand)
			drawTable(compHand, userHand)
			compHand = processCompTurn(compHand, userHand)
		} else {
			fmt.Println("COMP TURN")
			compHand = processCompTurn(compHand, userHand)
			drawTable(compHand, userHand)
			userHand = processUserTurn(userHand)
		}
		drawTable(compHand, userHand)
		fmt.Print("Press 'Enter' for the turn results...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		drawBattle(compHand, userHand)
		userHand.power -= userHand.selectedPower
		compHand.power -= compHand.selectedPower
		userHand.cards[userHand.selectedCard].playable = false
		compHand.cards[compHand.selectedCard].playable = false
		totalUserPower = userHand.cards[userHand.selectedCard].value * (userHand.selectedPower + 1)
		totalCompPower = compHand.cards[compHand.selectedCard].value * (compHand.selectedPower + 1)
		if totalUserPower > totalCompPower {
			compHand.health -= userHand.cards[userHand.selectedCard].damage
		} else {
			userHand.health -= compHand.cards[compHand.selectedCard].damage
		}

		if userHand.health < 1 || compHand.health < 1 {
			break
		}
		if i < HandSize-1 {
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

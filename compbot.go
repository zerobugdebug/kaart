package main

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"sort"
	"time"
)

//Genetic algorithm parameters
var (
	maxPopulationSize      int     = 100   //size of the population
	generationsLimit       int     = 2     //how many generations to generate
	crossoverRate          float32 = 0.9   //how often to do crossover 0%-100% in decimal
	mutationRate           float32 = 0.9   //how often to do mutation 0%-100% in decimal
	elitismRate            float32 = 0.2   //how many of the best indviduals to keep intact
	deadend                float32 = 10000 //round number to split between unscheduled tasks and real hours to complete
	tourneySampleSize      int     = 3     //sample size for the tournament selection, should be less than population size-number of elites
	crossoverParentsNumber int     = 2     //number of parents for the crossover
	maxCrossoverLength     int     = 3     //max number of sequential tasks to cross between Chromosomes
	maxMutatedGenes        int     = 3     //maximum number of mutated genes, min=2
	mutationTypePreference float32 = 0.5   //prefered mutation type rate. 0 = 100% swap mutation, 1 = 100% displacement mutation
)

const (
	threadsNum int = 256 //number of go routines to run simultaneously
	//handSize   int = 4
	//maxPower   int = 12
)

//Gene represents one gene in the chromosome
type gene struct {
	order int
	power int
}

//Chromosome is the type for the chromosome
type chromosome struct {
	genes   []gene
	fitness float32
}

//Population is a struct for the chromosomes and their hashes
type population struct {
	hashes      map[uint64]int
	chromosomes []chromosome
}

/*
//Packing slice of bytes to uint64 using the minimum amount of bits
func packBytes(unpacked []byte) uint64 {
	var packed uint64
	bitSize := bits.Len(uint(MaxPower))
	for _, v := range unpacked {
		packed = packed<<bitSize + uint64(v)
	}
	return packed
}
*/

func calcGenesHash(genes []gene) uint64 {

	//Convert slice into string representation
	allGenes := fmt.Sprint(genes)
	Logger.Debug("allGenes=", allGenes)

	//Calculate hash
	hashAlg := fnv.New64a()
	hashAlg.Write([]byte(allGenes))
	return hashAlg.Sum64()
}

//Calculate hash for the chromosome
func calcChromosomeHash(chromosome chromosome) uint64 {
	return calcGenesHash(chromosome.genes)
}

//Calculate hash for the chromosomes
func calcChromosomesHash(chromosomes []chromosome) map[uint64]int {
	hashMap := make(map[uint64]int)
	for i, v := range chromosomes {
		hashMap[calcChromosomeHash(v)] = i
	}
	return hashMap
}

//GenerateChromosome will generate a new chromosome for the GA
func generateChromosome(hand Hand) chromosome {
	var newChromosome chromosome
	var lenHand int
	for _, v := range hand.cards {
		if v.playable {
			lenHand++
		}
	}

	cardsOrder := rand.Perm(lenHand)
	var cardsPower []int
	cardsPower = make([]int, lenHand)
	totalPower := 0
	for i := range cardsPower {
		cardsPower[i] = rand.Intn(hand.power - totalPower + 1)
		totalPower += cardsPower[i]
	}
	Logger.Debug(cardsOrder)
	Logger.Debug(cardsPower)
	//Convert relative cardOrder to absolute
	cardsOrder = convertCardOrder(cardsOrder, hand)
	Logger.Debug(cardsOrder)

	//Filling up genes
	newChromosome.genes = make([]gene, lenHand)
	for i := range newChromosome.genes {
		newChromosome.genes[i].order = cardsOrder[i]
		newChromosome.genes[i].power = cardsPower[i]
	}
	Logger.Debug(newChromosome)
	return newChromosome
}

//GeneratePopulation will generate full population
func generatePopulation(hand Hand) population {
	var population population
	var chromosome chromosome
	var remainingChromosomesNumber, numPlayableCards, maxAvailableVariants int

	remainingChromosomesNumber = maxPopulationSize
	//Calculate maximum number of permutations
	for _, v := range hand.cards {
		if v.playable {
			numPlayableCards++
		}
	}
	//Generate temp slice to store card order numbers
	tempSlice := make([]int, numPlayableCards)
	for i := range tempSlice {
		tempSlice[i] = i
	}
	//Get the number of all possible orders of the cards for the specific amount of cards
	maxAvailableVariants = len(GetAllPermutations(tempSlice))
	Logger.Debug("maxAvailableVariants=", maxAvailableVariants)
	//Check to see if only card order variants emough to cover maxPopulationSize
	if maxAvailableVariants < maxPopulationSize {
		//Calculate number of possible combinations to get to the available hand.power and multiple it to the card order variants
		maxAvailableVariants *= len(GetAllPermutationsForSum(numPlayableCards, hand.power))
	}
	Logger.Debug("maxAvailableVariants=", maxAvailableVariants)
	//Select the min(maxAvailableVariants, maxPopulationSize) as remainingChromosomesNumber
	if maxAvailableVariants < maxPopulationSize {
		remainingChromosomesNumber = maxAvailableVariants
	}
	Logger.Debug("remainingChromosomesNumber=", remainingChromosomesNumber)

	population.hashes = make(map[uint64]int)
	for condition := true; condition; condition = remainingChromosomesNumber > 0 {
		chromosome = generateChromosome(hand)
		Logger.Debug(chromosome)
		hash := calcChromosomeHash(chromosome)
		Logger.Debug(hash)
		if _, ok := population.hashes[hash]; !ok {
			population.hashes[hash] = len(population.chromosomes)
			population.chromosomes = append(population.chromosomes, chromosome)
			remainingChromosomesNumber--
		}
	}
	Logger.Debug(population)
	return population
}

func copyChromosome(oldChromosome chromosome) chromosome {
	var newChromosome chromosome
	newChromosome.genes = make([]gene, len(oldChromosome.genes))
	copy(newChromosome.genes, oldChromosome.genes)
	newChromosome.fitness = oldChromosome.fitness
	return newChromosome
}

func copyChromosomes(oldChromosomes []chromosome) []chromosome {
	var newChromosomes []chromosome
	for _, v := range oldChromosomes {
		newChromosomes = append(newChromosomes, copyChromosome(v))
	}
	return newChromosomes
}

func copyPopulation(oldPopulation population) population {
	var newPopulation population
	//Copy chromosomes
	newPopulation.chromosomes = copyChromosomes(oldPopulation.chromosomes)

	//Copy hashes
	newPopulation.hashes = make(map[uint64]int)
	for k, v := range oldPopulation.hashes {
		newPopulation.hashes[k] = v
	}
	return newPopulation
}

//TransmogrifyPopulation will apply crossovers and mutations on non-elite Chromosomes
func transmogrifyPopulation(pop population) population {
	elitesNum := int(elitismRate * float32(len(pop.chromosomes)))
	//Logger.Info("elitesNum=", elitesNum)
	var newPopulation population
	var tempChromosomes []chromosome
	//Keep elites in the new population
	//	newPopulation = population[:elitesNum]
	//Logger.Info("OldElite=", population[0])
	newPopulation.chromosomes = copyChromosomes(pop.chromosomes[:elitesNum])
	//Recalculate hash for the elites
	//newPopulation.hashes = calcChromosomesHash(newPopulation.Chromosomes)
	//Logger.Info("NewElite=", newPopulation[0])
	Logger.Debug("newPopulation size with elites =", len(newPopulation.chromosomes))
	Logger.Debug("Best elite fitness =", newPopulation.chromosomes[0].fitness)
	//LoggerFile.Info("ELITES:", newPopulation[0].tasks)
	remainingChromosomesNumber := len(pop.chromosomes) - elitesNum
	Logger.Debug("remainingChromosomesNumber =", remainingChromosomesNumber)
	//Generate len(population)-elitesNum additonal Chromosomes
	for condition := true; condition; condition = remainingChromosomesNumber > 0 {
		tempChromosomes = make([]chromosome, crossoverParentsNumber)
		//Select crossoverParentsNumber from the population with Torunament Selection
		tempChromosomes = tourneySelect(pop.chromosomes, crossoverParentsNumber)
		Logger.Debug("tempPopulation size after tourney =", len(tempChromosomes))
		//Apply crossover to the tempPopulation
		//             tempChromosomes = crossoverChromosomesOX1(tempChromosomes)
		Logger.Debug("tempPopulation size after crossover =", len(tempChromosomes))
		//Apply mutation to the tempPopulation
		tempChromosomes = mutateChromosomes(tempChromosomes)
		Logger.Debug("tempPopulation size after mutation =", len(tempChromosomes))
		//Append tempPopulation to the new population, if indviduals are new
		for _, v := range tempChromosomes {
			tempHash := calcChromosomeHash(v)
			//If hash doesn't exist in the hashes map
			if _, ok := newPopulation.hashes[tempHash]; !ok {
				//Add hash with value of index of current Chromosome
				newPopulation.hashes[tempHash] = len(newPopulation.chromosomes)
				//Add Chromosome to the Chromosomes slice
				newPopulation.chromosomes = append(newPopulation.chromosomes, copyChromosome(v))
				remainingChromosomesNumber--
			}
		}

		Logger.Debug("newPopulation size =", len(newPopulation.chromosomes))
		//Update remaining number of Chromosomes to generate
		Logger.Debug("remainingChromosomesNumber =", remainingChromosomesNumber)
		Logger.Debug("condition =", condition)
	}

	Logger.Debug("newPopulation.hashes=", newPopulation.hashes)
	//Cut extra Chromosomes generated by mutation/crossover
	newPopulation.chromosomes = newPopulation.chromosomes[:len(pop.chromosomes)]
	return newPopulation

}

//Tournament selection for the crossover
func tourneySelect(chromosomes []chromosome, number int) []chromosome {
	//Create slice of randmoly permutated Chromosomes numbers
	sampleOrder := rand.Perm(len(chromosomes))
	Logger.Debug("sampleOrder =", sampleOrder)

	var bestChromosomes []chromosome
	var bestChromosomeNumber int
	var sampleOrderNumber int
	var bestChromosomeFitness float32
	for i := 0; i < number; i++ {
		Logger.Debug("Processing Chromosome =", i)

		bestChromosomeNumber = 0
		sampleOrderNumber = 0
		bestChromosomeFitness = float32(math.MaxFloat32)
		//Select best Chromosome number from first tourneySampleSize elements in sampleOrder
		for j, v := range sampleOrder[:tourneySampleSize] {
			Logger.Debugf("Processing sample %v, sample value %v", j, v)
			if chromosomes[v].fitness < bestChromosomeFitness {
				bestChromosomeNumber = v
				bestChromosomeFitness = chromosomes[v].fitness
				sampleOrderNumber = j
				Logger.Debug("bestChromosomeNumber =", bestChromosomeNumber)
				Logger.Debug("bestChromosomeFitness =", bestChromosomeFitness)
				Logger.Debug("sampleOrderNumber =", sampleOrderNumber)

			}
		}
		//Add best Chromosome to return slice
		bestChromosomes = append(bestChromosomes, chromosomes[bestChromosomeNumber])
		Logger.Debug("bestChromosomes size =", len(bestChromosomes))

		//Remove best Chromosome number from the selection
		//Using copy-last&truncate algorithm, due to O(1) complexity
		sampleOrder[sampleOrderNumber] = sampleOrder[len(sampleOrder)-1]
		sampleOrder = sampleOrder[:len(sampleOrder)-1]
		//Shuffle remaining Chromosome numbers
		rand.Shuffle(len(sampleOrder), func(i, j int) { sampleOrder[i], sampleOrder[j] = sampleOrder[j], sampleOrder[i] })
		Logger.Debug("new sampleOrder =", sampleOrder)

	}
	return bestChromosomes
}

func displacementMutation(chromosome chromosome) chromosome {
	//Randomly select number of genes to mutate, but at least 1
	//numOfGenesToMutate := rand.Intn(maxMutatedGenes) + 1
	/* 	for i := 0; i < numOfGenesToMutate; i++ {
	   		//Generate random old position for the gene between 0 and one element before last
	   		oldPosition := rand.Intn(len(Chromosome.genes)*8 - 1)
	   		//Generate random new position for the gene between oldPosition+1 and last element
	   		newPosition := rand.Intn(len(Chromosome.genes)*8-oldPosition-1) + oldPosition + 1
	   		//Store the original taskID at the oldPosition
	   		oldBit := Chromosome.genes[oldPosition/8] &
	   		//Shift all taskIDs one task back
	   		for j := range Chromosome.tasks[oldPosition:newPosition] {
	   			Chromosome.tasks[oldPosition+j].taskID = Chromosome.tasks[oldPosition+j+1].taskID
	   		}
	   		//Restore the original taskID to the newPosition
	   		Chromosome.tasks[newPosition].taskID = oldTaskID
	   	}
	*/
	return chromosome
}

func swapMutation(chromosome chromosome) chromosome {
	//Randomly select number of genes to mutate, but at least 1
	/* 	numOfGenesToMutate := rand.Intn(maxMutatedGenes-1) + 1
	   	sampleOrder := rand.Perm(len(Chromosome.tasks))
	   	for i := 0; i < numOfGenesToMutate; i++ {
	   		//Swap taskIDs for the task with number sampleOrder[i] and sampleOrder[len(Chromosome.tasks)-1] to make it easier to account for the border values
	   		Chromosome.tasks[sampleOrder[i]].taskID, Chromosome.tasks[sampleOrder[len(Chromosome.tasks)-i-1]].taskID = Chromosome.tasks[sampleOrder[len(Chromosome.tasks)-i-1]].taskID, Chromosome.tasks[sampleOrder[i]].taskID
	   	}
	*/return chromosome

}

func mutateChromosomes(chromosomes []chromosome) []chromosome {
	var mutatedChromosomes []chromosome
	/* 	//var crossoverStart, crossoverEnd, crossoverLen int
	   	//Copy parent to child Chromosomes slice
	   	//mutatedChromosomes = make([]Chromosome, len(Chromosomes))
	   	mutatedChromosomes = copyChromosomes(Chromosomes)
	   	for i := range mutatedChromosomes {
	   		//Check if we need to mutate
	   		if rand.Float32() < mutationRate {
	   			if rand.Float32() < mutationTypePreference {
	   				//Do the displacement mutation
	   				mutatedChromosomes[i] = displacementMutation(mutatedChromosomes[i])
	   			} else {
	   				//Do the swap mutation
	   				mutatedChromosomes[i] = swapMutation(mutatedChromosomes[i])
	   			}
	   		}
	   	} */
	return mutatedChromosomes
}

func calcChromosomesFitness(chromosomes []chromosome, compHand Hand, userHand Hand) {
	for i, v := range chromosomes {
		chromosomes[i].fitness = calcChromosomeFitness(v, compHand, userHand)
	}
}

func calcChromosomeFitness(chromosome chromosome, compHand Hand, userHand Hand) float32 {
	var totalMatches, totalWins, numCards int
	var winPercentage float32
	for _, v := range userHand.cards {
		if v.playable {
			numCards++
		}
	}
	//Generate temp slice to store card order numbers
	tempSlice := make([]int, numCards)
	for i := range tempSlice {
		tempSlice[i] = i
	}
	//Get all possible orders of the cards for the specific amount of cards
	cardOrders := GetAllPermutations(tempSlice)

	//Get all possible power combinations for numCards
	maxAvailablePower := userHand.power
	cardPowers := GetAllPermutationsForSum(numCards, maxAvailablePower)
	Logger.Debug(compHand.cards)
	Logger.Debug(userHand.cards)
	//TODO: Skip incorrect variants of cardOrder, if we know user selected card
	for _, cardOrder := range cardOrders {
		Logger.Debug(cardOrder)
		cardOrder = convertCardOrder(cardOrder, userHand)
		Logger.Debug(cardOrder)
		for _, cardPower := range cardPowers {
			//Logger.Debug(cardPower)
			totalMatches++
			compHealth := compHand.health
			userHealth := userHand.health
			for i, v := range cardOrder {
				//Logger.Debug(compHand.cards[chromosome.genes[i].order].value, chromosome.genes[i].power)
				//Logger.Debug(userHand.cards[v].value, cardPower[i])
				compAttack := compHand.cards[chromosome.genes[i].order].value * (chromosome.genes[i].power + 1)
				userAttack := userHand.cards[v].value * (cardPower[i] + 1)
				if compAttack > userAttack {
					userHealth -= compHand.cards[chromosome.genes[i].order].damage
				} else {
					compHealth -= userHand.cards[v].damage
				}
				//Logger.Debug(userHealth, ":", compHealth)

				if userHealth < 1 || compHealth < 1 {
					break
				}
			}
			if userHealth < 1 || userHealth <= compHealth {
				totalWins++
			}
		}
	}
	Logger.Debug(totalMatches)
	Logger.Debug(totalWins)
	winPercentage = float32(totalWins) / float32(totalMatches)
	return winPercentage
}

//convert relative card order (excluding played) to absolute order in hand
func convertCardOrder(cardOrder []int, hand Hand) []int {
	for i, v := range hand.cards {
		if !v.playable {
			for j, w := range cardOrder {
				Logger.Debug(i, v)
				Logger.Debug(j, w)
				Logger.Debug(cardOrder)
				if w >= i {
					cardOrder[j]++
					Logger.Debug(cardOrder)
				}
			}
		}
	}
	return cardOrder
}

func sortChromosomes(chromosomes []chromosome) {
	//Sort chromosomes in the order of fitness (descending) - from largest to smallest
	sort.Slice(chromosomes, func(i, j int) bool {
		return chromosomes[i].fitness > chromosomes[j].fitness
	})
}

func filterHand(hand Hand) Hand {
	var newHand Hand
	//Copying original hand
	newHand.health = hand.health
	newHand.power = hand.power
	newHand.selectedCard = hand.selectedCard
	newHand.selectedPower = hand.selectedPower
	for i, v := range hand.cards {
		if v.playable || i == hand.selectedCard {
			newHand.cards = append(newHand.cards, v)
		}
	}
	return newHand
}

//GetNextMove will return card number and power for the next comp move
func GetNextMove(compHand Hand, userHand Hand) (int, int) {
	//	var nextCardNumber, nextPower int
	rand.Seed(time.Now().UnixNano())
	newCompHand := filterHand(compHand)
	newUserHand := filterHand(userHand)
	Logger.Debug(newCompHand)
	Logger.Debug(newUserHand)
	population := generatePopulation(compHand)
	calcChromosomesFitness(population.chromosomes, compHand, userHand)
	sortChromosomes(population.chromosomes)

	//calcChromosomeFitness(population.chromosomes[0], newCompHand, newUserHand)
	Logger.Debug(population)
	Logger.Debug(population.chromosomes[0])
	return population.chromosomes[0].genes[0].order, population.chromosomes[0].genes[0].power
}

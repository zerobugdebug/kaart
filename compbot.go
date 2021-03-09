package main

import (
	"math"
	"math/bits"
	"math/rand"
)

//Genetic algorithm parameters
var (
	populationSize         int     = 10    //size of the population
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

//Chromosome is the type for the chromosome
type Chromosome struct {
	genes   uint64
	fitness float32
}

//Population is a struct for the chromosomes and their hashes
type Population struct {
	hashes      map[uint64]int
	Chromosomes []Chromosome
}

//Packing slice of bytes to uint64 using the minimum amount of bits
func packBytes(unpacked []byte) uint64 {
	var packed uint64
	bitSize := bits.Len(uint(MaxPower))
	for _, v := range unpacked {
		packed = packed<<bitSize + uint64(v)
	}
	return packed
}

/*
func calcGenesHash(genes uint16) uint64 {
	//Calculate hash
	hashAlg := fnv.New64a()
	bytes:= make([]byte, 4)
	binary.BigEndian.PutUint16(bytes, genes)
	hashAlg.Write(bytes)
	return hashAlg.Sum64()
}

//Calculate hash for the chromosome
func calcChromosomeHash(chromosome Chromosome) uint64 {
	return calcGenesHash(chromosome.genes)
}

//Calculate hash for the chromosomes
func calcChromosomesHash(chromosomes []Chromosome) map[uint64]int {
	hashMap := make(map[uint64]int)
	for i, v := range chromosomes {
		hashMap[calcChromosomeHash(v)] = i
	}
	return hashMap
} */

//GenerateChromosome will generate a new chromosome for the GA
func GenerateChromosome() Chromosome {
	var newChromosome Chromosome
	cardsOrder := rand.Perm(HandSize)[:HandSize-1]
	var cardsPower []int
	cardsPower = make([]int, HandSize-1)
	totalPower := 0
	for i := range cardsPower {
		cardsPower[i] = rand.Intn(MaxPower - totalPower + 1)
		totalPower += cardsPower[i]
	}
	Logger.Debug(cardsOrder)
	Logger.Debug(cardsPower)
	unpackedBytes := make([]byte, HandSize*(HandSize-1))
	bitSize := bits.Len(uint(MaxPower))
	for i, v := range cardsOrder {
		unpackedBytes[i*bitSize+v] = byte(cardsPower[i])
	}
	Logger.Debug(unpackedBytes)
	newChromosome.genes = packBytes(unpackedBytes)
	Logger.Debug(newChromosome)
	return newChromosome
}

//GeneratePopulation will generate full population
func GeneratePopulation() Population {
	var population Population
	var chromosome Chromosome
	remainingChromosomesNumber := populationSize
	population.hashes = make(map[uint64]int)
	for condition := true; condition; condition = remainingChromosomesNumber > 0 {
		chromosome = GenerateChromosome()
		if _, ok := population.hashes[chromosome.genes]; !ok {
			population.hashes[chromosome.genes] = len(population.Chromosomes)
			population.Chromosomes = append(population.Chromosomes, chromosome)
			remainingChromosomesNumber--
		}
	}
	Logger.Debug(population)
	return population
}

func copyChromosome(oldChromosome Chromosome) Chromosome {
	var newChromosome Chromosome
	newChromosome.genes = oldChromosome.genes
	newChromosome.fitness = oldChromosome.fitness
	return newChromosome
}

func copyChromosomes(oldChromosomes []Chromosome) []Chromosome {
	var newChromosomes []Chromosome
	for _, v := range oldChromosomes {
		newChromosomes = append(newChromosomes, copyChromosome(v))
	}
	return newChromosomes
}

func copyPopulation(oldPopulation Population) Population {
	var newPopulation Population
	//Copy chromosomes
	newPopulation.Chromosomes = copyChromosomes(oldPopulation.Chromosomes)

	//Copy hashes
	newPopulation.hashes = make(map[uint64]int)
	for k, v := range oldPopulation.hashes {
		newPopulation.hashes[k] = v
	}
	return newPopulation
}

//TransmogrifyPopulation will apply crossovers and mutations on non-elite Chromosomes
func transmogrifyPopulation(pop Population) Population {
	elitesNum := int(elitismRate * float32(len(pop.Chromosomes)))
	//Logger.Info("elitesNum=", elitesNum)
	var newPopulation Population
	var tempChromosomes []Chromosome
	//Keep elites in the new population
	//	newPopulation = population[:elitesNum]
	//Logger.Info("OldElite=", population[0])
	newPopulation.Chromosomes = copyChromosomes(pop.Chromosomes[:elitesNum])
	//Recalculate hash for the elites
	//newPopulation.hashes = calcChromosomesHash(newPopulation.Chromosomes)
	//Logger.Info("NewElite=", newPopulation[0])
	Logger.Debug("newPopulation size with elites =", len(newPopulation.Chromosomes))
	Logger.Debug("Best elite fitness =", newPopulation.Chromosomes[0].fitness)
	//LoggerFile.Info("ELITES:", newPopulation[0].tasks)
	remainingChromosomesNumber := len(pop.Chromosomes) - elitesNum
	Logger.Debug("remainingChromosomesNumber =", remainingChromosomesNumber)
	//Generate len(population)-elitesNum additonal Chromosomes
	for condition := true; condition; condition = remainingChromosomesNumber > 0 {
		tempChromosomes = make([]Chromosome, crossoverParentsNumber)
		//Select crossoverParentsNumber from the population with Torunament Selection
		tempChromosomes = tourneySelect(pop.Chromosomes, crossoverParentsNumber)
		Logger.Debug("tempPopulation size after tourney =", len(tempChromosomes))
		//Apply crossover to the tempPopulation
		//             tempChromosomes = crossoverChromosomesOX1(tempChromosomes)
		Logger.Debug("tempPopulation size after crossover =", len(tempChromosomes))
		//Apply mutation to the tempPopulation
		tempChromosomes = mutateChromosomes(tempChromosomes)
		Logger.Debug("tempPopulation size after mutation =", len(tempChromosomes))
		//Append tempPopulation to the new population, if indviduals are new
		for _, v := range tempChromosomes {
			//tempHash := calcChromosomeHash(v)
			//If hash doesn't exist in the hashes map
			if _, ok := newPopulation.hashes[v.genes]; !ok {
				//Add hash with value of index of current Chromosome
				newPopulation.hashes[v.genes] = len(newPopulation.Chromosomes)
				//Add Chromosome to the Chromosomes slice
				newPopulation.Chromosomes = append(newPopulation.Chromosomes, copyChromosome(v))
				remainingChromosomesNumber--
			}
		}

		Logger.Debug("newPopulation size =", len(newPopulation.Chromosomes))
		//Update remaining number of Chromosomes to generate
		Logger.Debug("remainingChromosomesNumber =", remainingChromosomesNumber)
		Logger.Debug("condition =", condition)
	}

	Logger.Debug("newPopulation.hashes=", newPopulation.hashes)
	//Cut extra Chromosomes generated by mutation/crossover
	newPopulation.Chromosomes = newPopulation.Chromosomes[:len(pop.Chromosomes)]
	return newPopulation

}

//Tournament selection for the crossover
func tourneySelect(chromosomes []Chromosome, number int) []Chromosome {
	//Create slice of randmoly permutated Chromosomes numbers
	sampleOrder := rand.Perm(len(chromosomes))
	Logger.Debug("sampleOrder =", sampleOrder)

	var bestChromosomes []Chromosome
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

func displacementMutation(Chromosome Chromosome) Chromosome {
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
	return Chromosome
}

func swapMutation(Chromosome Chromosome) Chromosome {
	//Randomly select number of genes to mutate, but at least 1
	/* 	numOfGenesToMutate := rand.Intn(maxMutatedGenes-1) + 1
	   	sampleOrder := rand.Perm(len(Chromosome.tasks))
	   	for i := 0; i < numOfGenesToMutate; i++ {
	   		//Swap taskIDs for the task with number sampleOrder[i] and sampleOrder[len(Chromosome.tasks)-1] to make it easier to account for the border values
	   		Chromosome.tasks[sampleOrder[i]].taskID, Chromosome.tasks[sampleOrder[len(Chromosome.tasks)-i-1]].taskID = Chromosome.tasks[sampleOrder[len(Chromosome.tasks)-i-1]].taskID, Chromosome.tasks[sampleOrder[i]].taskID
	   	}
	*/return Chromosome

}

func mutateChromosomes(Chromosomes []Chromosome) []Chromosome {
	var mutatedChromosomes []Chromosome
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

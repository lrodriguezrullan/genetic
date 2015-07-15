package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Chromosome struct {
	s       string
	fitness int
	mutant  bool
}

var poolSize = 250
var totalGenerations = 10000
var mutationRate = 0.002
var crossoverRate = 0.75

var pool = make([]Chromosome, poolSize)
var newPool = make([]Chromosome, poolSize)
var mutationRateInt = int(mutationRate * 1000)
var crossoverRateInt = int(crossoverRate * 1000)
var maxFitness = 0

var totalFitness int

func main() {

	if math.Mod(float64(poolSize), 2.0) > 0 {
		fmt.Println("Invalid pool size. Pool size must be even.")
		return
	}

	rand.Seed(time.Now().Unix())

	initializePool()
	computeFitness()

	var generation int

	for generation = 2; generation <= totalGenerations && maxFitness != 13; generation++ {
		populateNewPool()
		computeFitness()
	}

	fmt.Printf("Generation %d\n", generation-1)
	printPool()
	fmt.Printf("maxFitness: %d\n", maxFitness)

}

func randomChromosome() string {

	temp := ""

	for c := 0; c < 24; c++ {
		v := rand.Intn(2)

		if v > 0 {
			temp += "1"
		} else {
			temp += "0"
		}
	}

	return temp

}

func initializePool() {
	for i := 0; i < poolSize; i++ {
		pool[i].s = randomChromosome()
	}
}

func printPool() {
	for i := 0; i < poolSize; i++ {

		mutant := ""
		if pool[i].mutant {
			mutant = "M"
		}

		fmt.Printf("%d: %s [%d, %d, %d, %d, %d, %d, %d, %d] (%d) %f%% %s\n",
			i, pool[i].s,
			binStringToInt(pool[i].s[0:3]),
			binStringToInt(pool[i].s[3:6]),
			binStringToInt(pool[i].s[6:9]),
			binStringToInt(pool[i].s[9:12]),
			binStringToInt(pool[i].s[12:15]),
			binStringToInt(pool[i].s[15:18]),
			binStringToInt(pool[i].s[18:21]),
			binStringToInt(pool[i].s[21:24]),
			pool[i].fitness,
			float64(pool[i].fitness)/float64(totalFitness),
			mutant)
	}
	fmt.Printf("Total Fitness: %d\n", totalFitness)
}

func computeFitness() {
	totalFitness = 0

	for i := 0; i < poolSize; i++ {
		fitness := getFitness(&pool[i])
		totalFitness += fitness
		if fitness > maxFitness {
			maxFitness = fitness
		}
	}
}

func getFitness(chromosome *Chromosome) int {

	fitness := 0

	var g = make([]int, 8)

	g[0] = binStringToInt(chromosome.s[0:3])
	g[1] = binStringToInt(chromosome.s[3:6])
	g[2] = binStringToInt(chromosome.s[6:9])
	g[3] = binStringToInt(chromosome.s[9:12])
	g[4] = binStringToInt(chromosome.s[12:15])
	g[5] = binStringToInt(chromosome.s[15:18])
	g[6] = binStringToInt(chromosome.s[18:21])
	g[7] = binStringToInt(chromosome.s[21:24])

	for c := 0; c < 7; c++ {
		for d := c + 1; d < 8; d++ {
			if g[c] == g[d] {
				fitness++
				break
			}
		}
	}

	// diagonals

	for c := 0; c < 7; c++ {
		for d := c + 1; d < 8; d++ {

			e := int(math.Abs(float64(c - d)))

			if g[c]+e == g[d] || g[c]-e == g[d] {
				fitness++
				break
			}

		}
	}

	chromosome.fitness = 13 - fitness

	return chromosome.fitness
}

func binStringToInt(s string) int {
	i := 0
	for c := len(s) - 1; c >= 0; c-- {
		if s[c] == '1' {
			i += int(math.Pow(2, float64(len(s)-c-1)))
		}

	}
	return i
}

func selectParent() Chromosome {
	randomNumber := rand.Intn(totalFitness)
	accumulatedFitness := 0

	var index = 0

	for c := 0; accumulatedFitness < randomNumber; c++ {
		index = int(math.Mod(float64(c), float64(poolSize)))
		accumulatedFitness += pool[index].fitness
	}

	return pool[index]
}

func populateNewPool() {

	for c := 0; c < poolSize; c += 2 {

		p1 := selectParent()
		p2 := selectParent()

		r := rand.Intn(1000)
		if r < crossoverRateInt {
			newPool[c], newPool[c+1] = crossover(p1, p2)
		} else {
			newPool[c], newPool[c+1] = p1, p2
		}

		mutate(&newPool[c])
		mutate(&newPool[c+1])
	}
	copy(pool, newPool)
}

func mutate(a *Chromosome) {

	newString := ""

	for i := 0; i < len(a.s); i++ {
		c := a.s[i]
		r := rand.Intn(1000)
		if r < mutationRateInt {
			if c == '0' {
				c = '1'
			} else {
				c = '0'
			}
			a.mutant = true
		}
		newString = newString + string(c)
	}

	a.s = newString

}

func crossover(a Chromosome, b Chromosome) (Chromosome, Chromosome) {
	xop := rand.Intn(len(a.s))

	var d, e Chromosome

	d.s = a.s[0:xop] + b.s[xop:]
	e.s = b.s[0:xop] + a.s[xop:]

	return d, e
}

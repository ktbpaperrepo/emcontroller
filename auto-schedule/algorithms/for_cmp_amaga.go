package algorithms

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/KeepTheBeats/routing-algorithms/random"
	"github.com/astaxie/beego"
	"github.com/wcharczuk/go-chart"

	asmodel "emcontroller/auto-schedule/model"
)

/**
An algorithm for comparison, with name "Accept More Applications Genetic Algorithm (AMAGA)". This algorithm only aim at accepting as more applications as possible, and it is a genetic algorithm.
In the experiment, we will compare MCSSGA with this algorithm.
*/

// Accept More Applications Genetic Algorithm (AMAGA)
type Amaga struct {
	ChromosomesCount     int // One chromosome is a solution
	IterationCount       int // In each iteration, a population will be generated. One population consists of some solutions.
	CrossoverProbability float64
	MutationProbability  float64

	// If in the past StopNoUpdateIteration iterations, the best solution so far has not been updated, we should end this algorithm and return the best solution so far.
	StopNoUpdateIteration int
	CurNoUpdateIteration  int // record how many iterations the best solution has not updated currently.

	// these 2 member variables record the best solution in all iteration as well as its fitness value
	BestFitnessRecords []float64
	BestSolnRecords    []asmodel.Solution

	// record the best fitness value in every iteration, to show the evolution trend of the populations
	BestFitnessEachIter []float64
}

func NewAmaga(chromosomesCount int, iterationCount int, crossoverProbability float64, mutationProbability float64, stopNoUpdateIteration int) *Amaga {
	return &Amaga{
		ChromosomesCount:      chromosomesCount,
		IterationCount:        iterationCount,
		CrossoverProbability:  crossoverProbability,
		MutationProbability:   mutationProbability,
		StopNoUpdateIteration: stopNoUpdateIteration,
		CurNoUpdateIteration:  0,
		BestFitnessRecords:    nil,
		BestSolnRecords:       nil,
		BestFitnessEachIter:   nil,
	}
}

func (a *Amaga) Schedule(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, error) {
	beego.Info("Using scheduling algorithm:", AmagaName)

	// randomly generate the init population
	var initPopulation []asmodel.Solution = a.initialize(clouds, apps, appsOrder)

	// there are IterationCount+1 iterations in total, this is the No. 0 iteration
	currentPopulation := a.selectionOperator(clouds, apps, initPopulation) // Iteration No. 0

	// No. 1 iteration to No. m.IterationCount iteration
	for iteration := 1; iteration <= a.IterationCount; iteration++ {

		currentPopulation = a.crossoverOperator(clouds, apps, appsOrder, currentPopulation)

		currentPopulation = a.mutationOperator(clouds, apps, appsOrder, currentPopulation)

		currentPopulation = a.selectionOperator(clouds, apps, currentPopulation)

		// If we did not find better solutions in the past some iterations, we stop the algorithm and return the result.
		if a.CurNoUpdateIteration > a.StopNoUpdateIteration {
			break
		}
	}

	beego.Info("Best fitness in each iteration:", a.BestFitnessEachIter)
	beego.Info("Final BestFitnessRecords:", a.BestFitnessRecords)
	beego.Info("Total iteration number (the following 2 should be equal): ", len(a.BestFitnessRecords), len(a.BestSolnRecords))
	return a.BestSolnRecords[len(a.BestSolnRecords)-1], nil
}

// use "best effort random" method to generate some solutions as the init population
func (a *Amaga) initialize(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) []asmodel.Solution {
	var initPopulation []asmodel.Solution
	for i := 0; i < a.ChromosomesCount; i++ {
		var oneSolution asmodel.Solution = CmpRandomAcceptMostSolution(clouds, apps, appsOrder)
		initPopulation = append(initPopulation, oneSolution)
	}

	return initPopulation
}

// the selection operator of Genetic Algorithm
func (a *Amaga) selectionOperator(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, population []asmodel.Solution) []asmodel.Solution {

	// calculate the fitness of every chromosome in the current (old) population
	fitnesses := make([]float64, len(population))
	for i := 0; i < len(population); i++ {
		fitnesses[i] = a.Fitness(clouds, apps, population[i])
	}

	beego.Info("fitness values this iteration:", fitnesses)

	// do selection to generate a new population
	var newPopulation []asmodel.Solution
	pickHelper := make([]int, len(fitnesses)) // for binary tournament selection

	// to record the solution with the highest fitness value in the new population generated in this iteration
	var bestFitThisIter float64 = -math.MaxFloat64 // initialized with a very small value
	var bestFitThisIterIdx int = 0                 // the index of the solution with the highest fitness value in the old population

	// in every population, there should be m.ChromosomesCount chromosomes
	for i := 0; i < a.ChromosomesCount; i++ {
		var selChrIdx int // selected chromosome index in the input population

		// binary tournament selection
		picked := random.RandomPickN(pickHelper, 2)
		if fitnesses[picked[0]] > fitnesses[picked[1]] { // the larger, the better
			selChrIdx = picked[0]
		} else {
			selChrIdx = picked[1]
		}

		// put the selected solution into the new population
		newChromosome := asmodel.SolutionCopy(population[selChrIdx])
		newPopulation = append(newPopulation, newChromosome)

		// If the selected solution has the highest fitness value so far in this iteration, we save it.
		selFitness := fitnesses[selChrIdx]
		if selFitness > bestFitThisIter {
			bestFitThisIter = selFitness
			bestFitThisIterIdx = selChrIdx
		}
	}

	// At the end of selection operator, we record the best solution in every iteration.
	var bestFitAllIter float64
	var bestSolnAllIter asmodel.Solution

	if len(a.BestFitnessRecords) != len(a.BestSolnRecords) { // the 2 lengths should be equal, this check is for safety.
		panic(fmt.Sprintf("len(m.BestFitnessRecords) [%d] is not equal to len(m.BestSolnRecords) [%d]", len(a.BestFitnessRecords), len(a.BestSolnRecords)))
	}

	// We only record the best solutions until the current iteration.
	if len(a.BestFitnessRecords) == 0 { // In the 1st iteration, m.BestFitnessRecords and m.BestSolnRecords are nil.
		bestFitAllIter = bestFitThisIter
		bestSolnAllIter = population[bestFitThisIterIdx]
		a.CurNoUpdateIteration = 0
	} else {
		bestFitAllIter = a.BestFitnessRecords[len(a.BestFitnessRecords)-1]
		bestSolnAllIter = a.BestSolnRecords[len(a.BestSolnRecords)-1]
		if bestFitThisIter > bestFitAllIter {
			bestFitAllIter = bestFitThisIter
			bestSolnAllIter = population[bestFitThisIterIdx]
			a.CurNoUpdateIteration = 0
		} else {
			a.CurNoUpdateIteration++
		}
	}

	// record them
	a.BestFitnessEachIter = append(a.BestFitnessEachIter, bestFitThisIter)
	a.BestFitnessRecords = append(a.BestFitnessRecords, bestFitAllIter)
	a.BestSolnRecords = append(a.BestSolnRecords, asmodel.SolutionCopy(bestSolnAllIter))

	return newPopulation
}

// the fitness function of this algorithm
func (a *Amaga) Fitness(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, chromosome asmodel.Solution) float64 {
	var fitnessValue float64

	for appName, _ := range apps {
		if chromosome.AppsSolution[appName].Accepted { // This algorithm only aim at accepting as more applications as possible.
			fitnessValue += 1
		}
	}

	return fitnessValue
}

// the crossover operator of Genetic Algorithm.
func (a *Amaga) crossoverOperator(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, population []asmodel.Solution) []asmodel.Solution {
	// If a chromosome has less than 2 genes, we cannot do crossover.
	if len(apps) <= 1 {
		return population
	}

	// randomly choose the chromosomes that need crossover. We save their indexes.
	var idxNeedCrossover []int
	for i := 0; i < len(population); i++ {
		if random.RandomFloat64(0, 1) < a.CrossoverProbability {
			idxNeedCrossover = append(idxNeedCrossover, i)
		}
	}

	//beego.Info("idxNeedCrossover:", idxNeedCrossover) // for debug

	var crossoveredPopulation []asmodel.Solution

	// randomly choose chromosome pairs to do crossover
	var whetherCrossover []bool = make([]bool, len(population))
	for len(idxNeedCrossover) > 1 { // we can only do crossover when we have at list 2 chromosomes

		/**
		Before doing crossover, we do 3 things:
		1. choose two indexes of chromosomes for crossover;
		2. mark them in whetherCrossover
		3. delete them from idxNeedCrossover;
		*/

		// choose first index
		first := random.RandomInt(0, len(idxNeedCrossover)-1)
		firstIndex := idxNeedCrossover[first]
		// mark
		whetherCrossover[firstIndex] = true
		// delete
		idxNeedCrossover = append(idxNeedCrossover[:first], idxNeedCrossover[first+1:]...)

		// choose second index
		second := random.RandomInt(0, len(idxNeedCrossover)-1)
		secondIndex := idxNeedCrossover[second]
		// mark
		whetherCrossover[secondIndex] = true
		// delete
		idxNeedCrossover = append(idxNeedCrossover[:second], idxNeedCrossover[second+1:]...)

		/**
		Then, we do crossover.
		*/

		// get the 2 chromosomes for crossover. We use copy to avoid changing the original population
		firstChromosome := asmodel.SolutionCopy(population[firstIndex])
		secondChromosome := asmodel.SolutionCopy(population[secondIndex])

		newFirstChromosome, newSecondChromosome := AllPossTwoPointCrossover(firstChromosome, secondChromosome, clouds, apps, appsOrder)

		// append the two new chromosomes in crossoveredPopulation
		crossoveredPopulation = append(crossoveredPopulation, newFirstChromosome, newSecondChromosome)
	}

	//beego.Info("whetherCrossover:", whetherCrossover) // for debug

	// directly put the chromosomes without doing crossover to the new population
	for i := 0; i < len(population); i++ {
		if !whetherCrossover[i] {
			crossoveredPopulation = append(crossoveredPopulation, asmodel.SolutionCopy(population[i]))
		}
	}

	return crossoveredPopulation
}

// the mutation operator of Genetic Algorithm
func (a *Amaga) mutationOperator(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, population []asmodel.Solution) []asmodel.Solution {
	var mutatedPopulation []asmodel.Solution = make([]asmodel.Solution, len(population))

	for i := 0; i < len(population); i++ { // a chromosome

		for { // We repeat mutating a chromosome until the mutated new one is acceptable.
			mutatedChromosome := asmodel.GenEmptySoln()

			// gene-based mutation
			for appName, oriGene := range population[i].AppsSolution {
				// every gene has the probability "m.MutationProbability" to mutate
				if random.RandomFloat64(0, 1) < a.MutationProbability {
					mutatedChromosome.AppsSolution[appName] = a.geneMutate(clouds, oriGene) // mutate
				} else {
					mutatedChromosome.AppsSolution[appName] = asmodel.SasCopy(oriGene) // do not mutate
				}
			}

			// refine the mutated chromosome and check whether it is acceptable
			mutatedChromosome, acceptable := CmpRefineSoln(clouds, apps, appsOrder, mutatedChromosome)
			if acceptable {
				mutatedPopulation[i] = mutatedChromosome
				break
			}
		}

		/**
		In the above loop, we do not need to save the unacceptable solutions/chromosomes, because the mutated solutions/chromosomes are generated randomly and it is almost impossible to exclude the tried solutions from the following random attempts, so even if we save the unacceptable solutions/chromosomes, it can only save some time to run the function RefineSoln, but cannot reduce the number of random attempts, which I think is not worthy enough.
		*/

	}

	return mutatedPopulation
}

// The function to mutate a gene. After the mutation, a gene should become a different one, unless it is not accepted originally.
func (a *Amaga) geneMutate(clouds map[string]asmodel.Cloud, ori asmodel.SingleAppSolution) asmodel.SingleAppSolution {
	var mutated asmodel.SingleAppSolution = asmodel.SasCopy(asmodel.RejSoln)

	cloudsToPick := asmodel.CloudMapCopy(clouds)
	if ori.Accepted {
		delete(cloudsToPick, ori.TargetCloudName) // after mutation, the target cloud should be different
	}

	mutated.Accepted = random.RandomInt(0, 1) == 0 // 50% accept 50% not
	if mutated.Accepted {                          // Only when accepted, this gene needs a target cloud.
		mutated.TargetCloudName, _ = randomCloudMapPick(cloudsToPick)
	}

	return mutated
}

// draw a.BestFitnessEachIter and a.BestFitnessRecords on a line chart, to show the evolution trend
func (a *Amaga) DrawEvoChart() {
	var drawChartFunc func(http.ResponseWriter, *http.Request) = func(res http.ResponseWriter, r *http.Request) {
		var xValuesAllBest []float64
		for i, _ := range a.BestFitnessRecords {
			xValuesAllBest = append(xValuesAllBest, float64(i))
		}

		graph := chart.Chart{
			Title: "Evolution",
			//TitleStyle: chart.Style{
			//	Show: true,
			//},
			//Width: 600,
			//Height: 1800,
			//DPI:    300,
			XAxis: chart.XAxis{
				Name:      "Iteration Number",
				NameStyle: chart.StyleShow(),
				Style:     chart.StyleShow(),
				ValueFormatter: func(v interface{}) string {
					return strconv.FormatInt(int64(v.(float64)), 10)
				},
			},
			YAxis: chart.YAxis{
				AxisType:  chart.YAxisSecondary,
				Name:      "Fitness",
				NameStyle: chart.StyleShow(),
				Style:     chart.StyleShow(),
			},
			Background: chart.Style{
				Padding: chart.Box{
					Top:  50,
					Left: 20,
				},
			},
			Series: []chart.Series{
				chart.ContinuousSeries{
					Name:    "Best Fitness in all iteration",
					XValues: xValuesAllBest,
					YValues: a.BestFitnessRecords,
				},
				chart.ContinuousSeries{
					Name:    "Best Fitness in each iterations",
					XValues: xValuesAllBest,
					YValues: a.BestFitnessEachIter,
					Style: chart.Style{
						Show:            true,
						StrokeDashArray: []float64{5.0, 3.0, 2.0, 3.0},
						StrokeWidth:     1,
					},
				},
			},
		}

		graph.Elements = []chart.Renderable{
			chart.LegendThin(&graph),
		}

		res.Header().Set("Content-Type", "image/png")
		err := graph.Render(chart.PNG, res)
		if err != nil {
			log.Println("Error: graph.Render(chart.PNG, res)", err)
		}
	}

	http.HandleFunc("/", drawChartFunc)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error: http.ListenAndServe(\":8080\", nil)", err)
	}
}

package algorithms

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/KeepTheBeats/routing-algorithms/random"
	"github.com/astaxie/beego"
	"github.com/wcharczuk/go-chart"

	asmodel "emcontroller/auto-schedule/model"
)

/**
An algorithm for comparison, with name "Diktyo-GA". The paper "Diktyo: Network-Aware Scheduling in Container-based Clouds" proposes a Mixed-Integer Linear Programming (MILP) algorithm with the objective to reduce the network latencies among applications. I do not have enough time to implement an MILP algorithm in multi-cloud manager, and maybe an MILP algorithm will use too long scheduling time, so I customize the algorithm to suit my model and change it to a Genetic Algorithm (GA). I name this algorithm as "Diktyo-GA".
In the experiment, we will compare MCSSGA with this algorithm.
*/

// Diktyo-GA
type Diktyoga struct {
	ChromosomesCount     int // One chromosome is a solution
	IterationCount       int // In each iteration, a population will be generated. One population consists of some solutions.
	CrossoverProbability float64
	MutationProbability  float64

	// If in the past StopNoUpdateIteration iterations, the best solution so far has not been updated, we should end this algorithm and return the best solution so far.
	StopNoUpdateIteration int
	CurNoUpdateIteration  int // record how many iterations the best solution has not updated currently.

	MaxReachableRtt float64 // The biggest RTT between any 2 (or 1) reachable clouds, used to calculate fitness values. unit millisecond (ms)
	AvgDepNum       float64 // Average dependent application number of all applications

	// these 2 member variables record the best solution in all iteration as well as its fitness value
	BestFitnessRecords []float64
	BestSolnRecords    []asmodel.Solution

	// record the best fitness value in every iteration, to show the evolution trend of the populations
	BestFitnessEachIter []float64
}

func NewDiktyoga(chromosomesCount int, iterationCount int, crossoverProbability float64, mutationProbability float64, stopNoUpdateIteration int) *Diktyoga {
	return &Diktyoga{
		ChromosomesCount:      chromosomesCount,
		IterationCount:        iterationCount,
		CrossoverProbability:  crossoverProbability,
		MutationProbability:   mutationProbability,
		StopNoUpdateIteration: stopNoUpdateIteration,
		CurNoUpdateIteration:  0,
		MaxReachableRtt:       0,
		BestFitnessRecords:    nil,
		BestSolnRecords:       nil,
		BestFitnessEachIter:   nil,
	}
}

// Traverse all clouds to find the max RTT between any 2 (or 1) reachable clouds, and set it as the MaxReachableRtt of Diktyoga.
func (d *Diktyoga) SetMaxReaRtt(clouds map[string]asmodel.Cloud) {
	var maxReaRtt float64 = 0
	for _, srcCloud := range clouds {
		for _, ns := range srcCloud.NetState {
			if ns.Rtt > maxReaRtt && ns.Rtt < maxAccRttMs {
				maxReaRtt = ns.Rtt
			}
		}
	}

	d.MaxReachableRtt = maxReaRtt
}

// Traverse all applications to find the average dependency number of all apps.
func (d *Diktyoga) SetAvgDepNum(apps map[string]asmodel.Application) {
	var sumDepNum int = 0
	for _, app := range apps {
		sumDepNum += len(app.Dependencies)
	}
	d.AvgDepNum = float64(sumDepNum) / float64(len(apps))
}

func (d *Diktyoga) Schedule(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, error) {
	beego.Info("Using scheduling algorithm:", DiktyogaName)
	d.SetMaxReaRtt(clouds)
	beego.Info("MaxReachableRtt:", d.MaxReachableRtt)
	d.SetAvgDepNum(apps)
	beego.Info("AvgDepNum:", d.AvgDepNum)

	// randomly generate the init population
	var initPopulation []asmodel.Solution = d.initialize(clouds, apps, appsOrder)

	// there are IterationCount+1 iterations in total, this is the No. 0 iteration
	currentPopulation := d.selectionOperator(clouds, apps, initPopulation) // Iteration No. 0

	// No. 1 iteration to No. m.IterationCount iteration
	for iteration := 1; iteration <= d.IterationCount; iteration++ {

		currentPopulation = d.crossoverOperator(clouds, apps, appsOrder, currentPopulation)

		currentPopulation = d.mutationOperator(clouds, apps, appsOrder, currentPopulation)

		currentPopulation = d.selectionOperator(clouds, apps, currentPopulation)

		// If we did not find better solutions in the past some iterations, we stop the algorithm and return the result.
		if d.CurNoUpdateIteration > d.StopNoUpdateIteration {
			break
		}
	}

	beego.Info("Best fitness in each iteration:", d.BestFitnessEachIter)
	beego.Info("Final BestFitnessRecords:", d.BestFitnessRecords)
	beego.Info("Total iteration number (the following 2 should be equal): ", len(d.BestFitnessRecords), len(d.BestSolnRecords))
	return d.BestSolnRecords[len(d.BestSolnRecords)-1], nil
}

// use "best effort random" method to generate some solutions as the init population
func (d *Diktyoga) initialize(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) []asmodel.Solution {
	var initPopulation []asmodel.Solution
	for i := 0; i < d.ChromosomesCount; i++ {
		var oneSolution asmodel.Solution = CmpRandomAcceptMostSolution(clouds, apps, appsOrder)
		initPopulation = append(initPopulation, oneSolution)
	}

	return initPopulation
}

// the selection operator of Genetic Algorithm
func (d *Diktyoga) selectionOperator(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, population []asmodel.Solution) []asmodel.Solution {

	// calculate the fitness of every chromosome in the current (old) population
	fitnesses := make([]float64, len(population))

	var fitMu sync.Mutex  // the slice in golang is not safe for concurrent read/write
	var wg sync.WaitGroup // calculate the fitness of every chromosome in parallel
	for i := 0; i < len(population); i++ {
		wg.Add(1)
		go func(chromIdx int) {
			defer wg.Done()

			thisFitness := d.Fitness(clouds, apps, population[chromIdx])

			fitMu.Lock()
			fitnesses[chromIdx] = thisFitness
			fitMu.Unlock()
		}(i)
	}
	wg.Wait()

	beego.Info("fitness values this iteration:", fitnesses)

	// do selection to generate a new population
	var newPopulation []asmodel.Solution
	pickHelper := make([]int, len(fitnesses)) // for binary tournament selection

	// to record the solution with the highest fitness value in the new population generated in this iteration
	var bestFitThisIter float64 = -math.MaxFloat64 // initialized with a very small value
	var bestFitThisIterIdx int = 0                 // the index of the solution with the highest fitness value in the old population

	// in every population, there should be m.ChromosomesCount chromosomes
	for i := 0; i < d.ChromosomesCount; i++ {
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

	if len(d.BestFitnessRecords) != len(d.BestSolnRecords) { // the 2 lengths should be equal, this check is for safety.
		panic(fmt.Sprintf("len(m.BestFitnessRecords) [%d] is not equal to len(m.BestSolnRecords) [%d]", len(d.BestFitnessRecords), len(d.BestSolnRecords)))
	}

	// We only record the best solutions until the current iteration.
	if len(d.BestFitnessRecords) == 0 { // In the 1st iteration, m.BestFitnessRecords and m.BestSolnRecords are nil.
		bestFitAllIter = bestFitThisIter
		bestSolnAllIter = population[bestFitThisIterIdx]
		d.CurNoUpdateIteration = 0
	} else {
		bestFitAllIter = d.BestFitnessRecords[len(d.BestFitnessRecords)-1]
		bestSolnAllIter = d.BestSolnRecords[len(d.BestSolnRecords)-1]
		if bestFitThisIter > bestFitAllIter {
			bestFitAllIter = bestFitThisIter
			bestSolnAllIter = population[bestFitThisIterIdx]
			d.CurNoUpdateIteration = 0
		} else {
			d.CurNoUpdateIteration++
		}
	}

	// record them
	d.BestFitnessEachIter = append(d.BestFitnessEachIter, bestFitThisIter)
	d.BestFitnessRecords = append(d.BestFitnessRecords, bestFitAllIter)
	d.BestSolnRecords = append(d.BestSolnRecords, asmodel.SolutionCopy(bestSolnAllIter))

	return newPopulation
}

// the fitness function of this algorithm
func (d *Diktyoga) Fitness(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, chromosome asmodel.Solution) float64 {
	var fitnessValue float64

	for appName, _ := range apps {
		fitnessValue += d.fitnessOneApp(clouds, apps, chromosome, appName)
	}

	return fitnessValue
}

// calculate the fitness value contributed by an application. According to the paper Diktyo, this fitness function only considers acceptance rate and network latency.
func (d *Diktyoga) fitnessOneApp(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, chromosome asmodel.Solution, thisAppName string) float64 {
	var thisAppFitness float64 // result

	// if an application is rejected, it contributes a big negative fitness value, this is to encourage higher priority-weighted acceptance rate
	if !chromosome.AppsSolution[thisAppName].Accepted {
		thisAppFitness = -(d.MaxReachableRtt * d.AvgDepNum) / 2
	} else {

		// if this app is accepted, all its dependent apps are also accepted, which is guaranteed by our dependency acceptable check

		netPart := d.MaxReachableRtt * d.AvgDepNum // the base network part of fitness.
		for _, dep := range apps[thisAppName].Dependencies {
			depAppName := dep.AppName

			// calculate the network part of the fitness value of this dependency
			thisCloudName := chromosome.AppsSolution[thisAppName].TargetCloudName
			depCloudName := chromosome.AppsSolution[depAppName].TargetCloudName
			thisNodeName := chromosome.AppsSolution[thisAppName].K8sNodeName
			depNodeName := chromosome.AppsSolution[depAppName].K8sNodeName

			var thisRtt float64 // RTT from this application to the dependent application

			if thisNodeName == depNodeName { // We consider the RTT inside a same VM as 0.
				thisRtt = 0
			} else {
				thisRtt = clouds[thisCloudName].NetState[depCloudName].Rtt
			}

			netPart -= thisRtt
			/**
			The network delay between this app and its dependent ones will reduce this fitness.
			This fitness function tend to accept apps with fewer dependencies, this is an unfair point, but I cannot find a better way, and this way seems like the most fair one that I can find yet now.
			*/
		}

		// For an application with many dependencies, after the above calculation netPart may become a negative value, so we should set it to 0 in this condition, because otherwise accepting an application may be worse than rejecting it.
		if netPart < 0 {
			netPart = 0
		}

		thisAppFitness = netPart
	}

	return thisAppFitness
}

// the crossover operator of Genetic Algorithm.
func (d *Diktyoga) crossoverOperator(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, population []asmodel.Solution) []asmodel.Solution {
	// If a chromosome has less than 2 genes, we cannot do crossover.
	if len(apps) <= 1 {
		return population
	}

	// randomly choose the chromosomes that need crossover. We save their indexes.
	var idxNeedCrossover []int
	for i := 0; i < len(population); i++ {
		if random.RandomFloat64(0, 1) < d.CrossoverProbability {
			idxNeedCrossover = append(idxNeedCrossover, i)
		}
	}

	//beego.Info("idxNeedCrossover:", idxNeedCrossover) // for debug

	var crossoveredPopulation []asmodel.Solution

	// randomly choose chromosome pairs to do crossover
	var whetherCrossover []bool = make([]bool, len(population))

	var crpoMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	var wg sync.WaitGroup // calculate the fitness of every chromosome in parallel

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

		// this AllPossTwoPointCrossover has big workload, and it can be concurrent so we do it concurrently.
		wg.Add(1)
		go func() {
			defer wg.Done()
			newFirstChromosome, newSecondChromosome := CmpAllPossTwoPointCrossover(firstChromosome, secondChromosome, clouds, apps, appsOrder)
			// append the two new chromosomes in crossoveredPopulation
			crpoMu.Lock()
			crossoveredPopulation = append(crossoveredPopulation, newFirstChromosome, newSecondChromosome)
			crpoMu.Unlock()
		}()
	}
	wg.Wait()

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
func (d *Diktyoga) mutationOperator(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, population []asmodel.Solution) []asmodel.Solution {
	var mutatedPopulation []asmodel.Solution = make([]asmodel.Solution, len(population))

	var popuMu sync.Mutex // the slice in golang is not safe for concurrent read/write
	var wg sync.WaitGroup // mutate every chromosome in parallel

	for i := 0; i < len(population); i++ { // a chromosome
		wg.Add(1)

		go func(chromIdx int) {
			defer wg.Done()

			for { // We repeat mutating a chromosome until the mutated new one is acceptable.
				mutatedChromosome := asmodel.GenEmptySoln()

				// gene-based mutation
				for appName, oriGene := range population[chromIdx].AppsSolution {
					// every gene has the probability "m.MutationProbability" to mutate
					if random.RandomFloat64(0, 1) < d.MutationProbability {
						mutatedChromosome.AppsSolution[appName] = d.geneMutate(clouds, oriGene) // mutate
					} else {
						mutatedChromosome.AppsSolution[appName] = asmodel.SasCopy(oriGene) // do not mutate
					}
				}

				// refine the mutated chromosome and check whether it is acceptable
				mutatedChromosome, acceptable := CmpRefineSoln(clouds, apps, appsOrder, mutatedChromosome)
				if acceptable {
					popuMu.Lock()
					mutatedPopulation[chromIdx] = mutatedChromosome
					popuMu.Unlock()
					break
				}
			}

			/**
			In the above loop, we do not need to save the unacceptable solutions/chromosomes, because the mutated solutions/chromosomes are generated randomly and it is almost impossible to exclude the tried solutions from the following random attempts, so even if we save the unacceptable solutions/chromosomes, it can only save some time to run the function RefineSoln, but cannot reduce the number of random attempts, which I think is not worthy enough.
			*/
		}(i)
	}
	wg.Wait()

	return mutatedPopulation
}

// The function to mutate a gene. After the mutation, a gene should become a different one, unless it is not accepted originally.
func (d *Diktyoga) geneMutate(clouds map[string]asmodel.Cloud, ori asmodel.SingleAppSolution) asmodel.SingleAppSolution {
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

// draw d.BestFitnessEachIter and d.BestFitnessRecords on a line chart, to show the evolution trend
func (d *Diktyoga) DrawEvoChart() {
	var drawChartFunc func(http.ResponseWriter, *http.Request) = func(res http.ResponseWriter, r *http.Request) {
		var xValuesAllBest []float64
		for i, _ := range d.BestFitnessRecords {
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
					YValues: d.BestFitnessRecords,
				},
				chart.ContinuousSeries{
					Name:    "Best Fitness in each iterations",
					XValues: xValuesAllBest,
					YValues: d.BestFitnessEachIter,
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

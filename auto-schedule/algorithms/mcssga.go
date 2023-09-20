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
	chart "github.com/wcharczuk/go-chart"

	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
)

/**
NOTE:
In our fitness function, we do not need to consider to make unacceptable chromosomes harder to be selected (such as "fitness /= 3 for un accepted ones"), because:
1. this is complicated;
2. we need to use the allocatedCpuCore to calculate fitness values, but unacceptable solution cannot have this.
Instead, in our init function, mutation selector, crossover selector, we should guarantee that all generated solutions/chromosomes are acceptable. For mutation and crossover if we get an unacceptable solution, we can repeat, excluding the generated unacceptable ones, until we get an acceptable solution.
*/

const (
	// TODO: currently, this is measured. In the future, this may be able to:
	// possibility 1. be set by a dedicated user request;
	// possibility 2. be set in the user request for auto-scheduling;
	// possibility 3. be set for every application separately;
	// DefaultExpAppCompuTimeOneCpu float64 = 50
	DefaultExpAppCompuTimeOneCpu float64 = 35 // expected computation time by one CPU core, unit: ms

	enlargerScaleMaxCompuTime float64 = 1    // minimum allocated CPU is 1, so the max computation time should enlarger 1 times.
	enlargerScaleMaxRTT       float64 = 1.25 // because the minimum computation part is not 0, minimum RTT should also not be 0.
)

// Multi-Cloud Service Scheduling Genetic Algorithm (MCSSGA)
type Mcssga struct {
	ChromosomesCount     int // One chromosome is a solution
	IterationCount       int // In each iteration, a population will be generated. One population consists of some solutions.
	CrossoverProbability float64
	MutationProbability  float64

	// If in the past StopNoUpdateIteration iterations, the best solution so far has not been updated, we should end this algorithm and return the best solution so far.
	StopNoUpdateIteration int
	CurNoUpdateIteration  int // record how many iterations the best solution has not updated currently.

	MaxReachableRtt       float64            // The biggest RTT between any 2 (or 1) reachable clouds, used to calculate fitness values. unit millisecond (ms)
	AvgDepNum             float64            // Average dependent application number of all applications
	ExpAppCompuTimeOneCpu float64            // the expected computation time of every application by one CPU core.  unit millisecond (ms)
	FitnessNonPriDp       map[string]float64 // the record for dynamic programming in Fitness calculation, to reduce the scheduling time.

	// these 2 member variables record the best solution in all iteration as well as its fitness value
	BestFitnessRecords []float64
	BestSolnRecords    []asmodel.Solution

	// record the best fitness value in every iteration, to show the evolution trend of the populations
	BestFitnessEachIter []float64
}

func NewMcssga(chromosomesCount int, iterationCount int, crossoverProbability float64, mutationProbability float64, stopNoUpdateIteration int, exTimeOneCpu float64) *Mcssga {
	return &Mcssga{
		ChromosomesCount:      chromosomesCount,
		IterationCount:        iterationCount,
		CrossoverProbability:  crossoverProbability,
		MutationProbability:   mutationProbability,
		StopNoUpdateIteration: stopNoUpdateIteration,
		CurNoUpdateIteration:  0,
		MaxReachableRtt:       0,
		ExpAppCompuTimeOneCpu: exTimeOneCpu,
		BestFitnessRecords:    nil,
		BestSolnRecords:       nil,
		BestFitnessEachIter:   nil,
	}
}

// Traverse all clouds to find the max RTT between any 2 (or 1) reachable clouds, and set it as the MaxReachableRtt of Mcssga.
func (m *Mcssga) SetMaxReaRtt(clouds map[string]asmodel.Cloud) {
	var maxReaRtt float64 = 0
	for _, srcCloud := range clouds {
		for _, ns := range srcCloud.NetState {
			if ns.Rtt > maxReaRtt && ns.Rtt < maxAccRttMs {
				maxReaRtt = ns.Rtt
			}
		}
	}

	m.MaxReachableRtt = maxReaRtt
}

// Traverse all applications to find the average dependency number of all apps.
func (m *Mcssga) SetAvgDepNum(apps map[string]asmodel.Application) {
	var sumDepNum int = 0
	for _, app := range apps {
		sumDepNum += len(app.Dependencies)
	}
	m.AvgDepNum = float64(sumDepNum) / float64(len(apps))
}

func (m *Mcssga) Schedule(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, error) {
	beego.Info("Using scheduling algorithm:", McssgaName)
	m.SetMaxReaRtt(clouds)
	beego.Info("MaxReachableRtt:", m.MaxReachableRtt)
	m.SetAvgDepNum(apps)
	beego.Info("AvgDepNum:", m.AvgDepNum)

	beego.Info("Clouds:")
	for _, cloud := range clouds {
		beego.Info(models.JsonString(cloud))
	}
	beego.Info("Applications:")
	for _, app := range apps {
		beego.Info(models.JsonString(app))
	}
	beego.Info("Applications:", models.JsonString(apps))
	beego.Info("Clouds:", models.JsonString(clouds))
	beego.Info("appsOrder:", models.JsonString(appsOrder))

	// randomly generate the init population
	var initPopulation []asmodel.Solution = m.initialize(clouds, apps, appsOrder)
	//beego.Info("initPopulation:")
	//for _, soln := range initPopulation {
	//	beego.Info(models.JsonString(soln))
	//}

	// there are IterationCount+1 iterations in total, this is the No. 0 iteration
	currentPopulation := m.selectionOperator(clouds, apps, initPopulation) // Iteration No. 0

	// No. 1 iteration to No. m.IterationCount iteration
	for iteration := 1; iteration <= m.IterationCount; iteration++ {

		currentPopulation = m.crossoverOperator(clouds, apps, appsOrder, currentPopulation)

		currentPopulation = m.mutationOperator(clouds, apps, appsOrder, currentPopulation)

		currentPopulation = m.selectionOperator(clouds, apps, currentPopulation)

		// If we did not find better solutions in the past some iterations, we stop the algorithm and return the result.
		if m.CurNoUpdateIteration > m.StopNoUpdateIteration {
			break
		}
	}

	beego.Info("Best fitness in each iteration:", m.BestFitnessEachIter)
	beego.Info("Final BestFitnessRecords:", m.BestFitnessRecords)
	beego.Info("Total iteration number (the following 2 should be equal): ", len(m.BestFitnessRecords), len(m.BestSolnRecords))
	return m.BestSolnRecords[len(m.BestSolnRecords)-1], nil
}

// randomly generate some solutions as the init population
func (m *Mcssga) initialize(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) []asmodel.Solution {
	var initPopulation []asmodel.Solution
	for i := 0; i < m.ChromosomesCount; i++ {
		var oneSolution asmodel.Solution = RandomAcceptMostSolution(clouds, apps, appsOrder)
		initPopulation = append(initPopulation, oneSolution)
	}
	return initPopulation
}

// the crossover operator of Genetic Algorithm.
func (m *Mcssga) crossoverOperator(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, population []asmodel.Solution) []asmodel.Solution {
	// If a chromosome has less than 2 genes, we cannot do crossover.
	if len(apps) <= 1 {
		return population
	}

	// randomly choose the chromosomes that need crossover. We save their indexes.
	var idxNeedCrossover []int
	for i := 0; i < len(population); i++ {
		if random.RandomFloat64(0, 1) < m.CrossoverProbability {
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

// Randomly explore all possibilities of 2-point crossover, to try to get an acceptable solution. If this function cannot find an acceptable solution after trying all possibilities, it will return the original 2 chromosomes without doing crossover.
func AllPossTwoPointCrossover(firstChromosome asmodel.Solution, secondChromosome asmodel.Solution, clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, asmodel.Solution) {
	// in our unit tests, we will set both the input cloud and apps as nil
	var testMode bool = clouds == nil && apps == nil

	if len(firstChromosome.AppsSolution) != len(secondChromosome.AppsSolution) || len(firstChromosome.AppsSolution) != len(appsOrder) {
		panic(fmt.Sprintf("len(firstChromosome.AppsSolution) is %d; len(secondChromosome.AppsSolution) is %d; len(appsOrder) is %d. They should be equal.", len(firstChromosome.AppsSolution), len(secondChromosome.AppsSolution), len(appsOrder)))
	}

	// the number of genes in a chromosome, also the number of applications to schedule.
	geneNumber := len(firstChromosome.AppsSolution)

	/**
	We exchange the genes of the 2 chromosomes in the closed interval [point1, point2]. The width of this closed interval ranges from 1 to geneNumber-1.
	The following loop randomly traverse all possibility of different width and point1 which determine point2.
	*/

	// build an array to help select point widths randomly
	var possiblePointWidths []int
	for pointWidth := 1; pointWidth <= geneNumber-1; pointWidth++ {
		possiblePointWidths = append(possiblePointWidths, pointWidth)
	}
	for len(possiblePointWidths) > 0 {
		// randomly select a possible point width, and then remove it from the array, in order not to select it again.
		widthIdx := random.RandomInt(0, len(possiblePointWidths)-1)
		pointWidth := possiblePointWidths[widthIdx]
		possiblePointWidths = append(possiblePointWidths[:widthIdx], possiblePointWidths[widthIdx+1:]...)
		if testMode {
			beego.Info("pointWidth is:", pointWidth) // for debug
		}

		// build an array to help select point1 randomly
		var possiblePoint1 []int
		for point1 := 0; calcPoint2(point1, pointWidth) < geneNumber; point1++ {
			possiblePoint1 = append(possiblePoint1, point1)
		}
		for len(possiblePoint1) > 0 {
			// randomly select a possible point1, and then remove it from the array, in order not to select it again.
			pointIdx := random.RandomInt(0, len(possiblePoint1)-1)
			point1 := possiblePoint1[pointIdx]
			possiblePoint1 = append(possiblePoint1[:pointIdx], possiblePoint1[pointIdx+1:]...)

			// calculate point 2 by point 1
			point2 := calcPoint2(point1, pointWidth)

			if testMode {
				beego.Info("point1, point2:", point1, point2) // for debug
			} else {

				/**
				Then we do crossover with the randomly selected point1 and point2.
				We set the tryFunc here, because with this the unit tests will be easier to make.
				*/

				// if in this possibility the 2 crossovered chromosomes are acceptable, return them.
				crossoveredChromosome1, crossoveredChromosome2 := twoPointCrossover(firstChromosome, secondChromosome, appsOrder, point1, point2)

				// refine the 2 crossovered chromosomes and check whether they are acceptable. If both of them are acceptable, we return them as the result.
				if crossoveredChromosome1, acceptable1 := RefineSoln(clouds, apps, appsOrder, crossoveredChromosome1); acceptable1 {
					if crossoveredChromosome2, acceptable2 := RefineSoln(clouds, apps, appsOrder, crossoveredChromosome2); acceptable2 {
						return crossoveredChromosome1, crossoveredChromosome2
					}
				}
			}

		}
		if testMode {
			fmt.Println() // for debug
		}
	}

	return firstChromosome, secondChromosome
}

// try to do crossover on 2 chromosomes with the determined point1 and point2
func twoPointCrossover(firstChromosome asmodel.Solution, secondChromosome asmodel.Solution, appsOrder []string, point1, point2 int) (asmodel.Solution, asmodel.Solution) {
	var crossoveredFirstCh, crossoveredSecondCh asmodel.Solution = asmodel.GenEmptySoln(), asmodel.GenEmptySoln()

	for orderIdx, appName := range appsOrder {
		if orderIdx >= point1 && orderIdx <= point2 { // in the closed interval [point1, point2], exchange
			crossoveredFirstCh.AppsSolution[appName] = asmodel.SasCopy(secondChromosome.AppsSolution[appName])
			crossoveredSecondCh.AppsSolution[appName] = asmodel.SasCopy(firstChromosome.AppsSolution[appName])
		} else { // not in the closed interval [point1, point2], do not exchange
			crossoveredFirstCh.AppsSolution[appName] = asmodel.SasCopy(firstChromosome.AppsSolution[appName])
			crossoveredSecondCh.AppsSolution[appName] = asmodel.SasCopy(secondChromosome.AppsSolution[appName])
		}
	}

	return crossoveredFirstCh, crossoveredSecondCh
}

// This function is to calculate point 2. We will exchange the genes of the 2 chromosomes in the closed interval [point1, point2].
func calcPoint2(point1 int, pointWidth int) int {
	return point1 + pointWidth - 1
}

// the mutation operator of Genetic Algorithm
func (m *Mcssga) mutationOperator(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string, population []asmodel.Solution) []asmodel.Solution {
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
					if random.RandomFloat64(0, 1) < m.MutationProbability {
						mutatedChromosome.AppsSolution[appName] = m.geneMutate(clouds, oriGene) // mutate
					} else {
						mutatedChromosome.AppsSolution[appName] = asmodel.SasCopy(oriGene) // do not mutate
					}
				}

				// refine the mutated chromosome and check whether it is acceptable
				mutatedChromosome, acceptable := RefineSoln(clouds, apps, appsOrder, mutatedChromosome)
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
func (m *Mcssga) geneMutate(clouds map[string]asmodel.Cloud, ori asmodel.SingleAppSolution) asmodel.SingleAppSolution {
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

// the selection operator of Genetic Algorithm
func (m *Mcssga) selectionOperator(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, population []asmodel.Solution) []asmodel.Solution {

	// calculate the fitness of every chromosome in the current (old) population
	fitnesses := make([]float64, len(population))

	var fitMu sync.Mutex  // the slice in golang is not safe for concurrent read/write
	var wg sync.WaitGroup // calculate the fitness of every chromosome in parallel
	for i := 0; i < len(population); i++ {
		wg.Add(1)
		go func(chromIdx int) {
			defer wg.Done()

			thisFitness := m.Fitness(clouds, apps, population[chromIdx])

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
	for i := 0; i < m.ChromosomesCount; i++ {
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

	if len(m.BestFitnessRecords) != len(m.BestSolnRecords) { // the 2 lengths should be equal, this check is for safety.
		panic(fmt.Sprintf("len(m.BestFitnessRecords) [%d] is not equal to len(m.BestSolnRecords) [%d]", len(m.BestFitnessRecords), len(m.BestSolnRecords)))
	}

	// We only record the best solutions until the current iteration.
	if len(m.BestFitnessRecords) == 0 { // In the 1st iteration, m.BestFitnessRecords and m.BestSolnRecords are nil.
		bestFitAllIter = bestFitThisIter
		bestSolnAllIter = population[bestFitThisIterIdx]
		m.CurNoUpdateIteration = 0
	} else {
		bestFitAllIter = m.BestFitnessRecords[len(m.BestFitnessRecords)-1]
		bestSolnAllIter = m.BestSolnRecords[len(m.BestSolnRecords)-1]
		if bestFitThisIter > bestFitAllIter {
			bestFitAllIter = bestFitThisIter
			bestSolnAllIter = population[bestFitThisIterIdx]
			m.CurNoUpdateIteration = 0
		} else {
			m.CurNoUpdateIteration++
		}
	}

	// record them
	m.BestFitnessEachIter = append(m.BestFitnessEachIter, bestFitThisIter)
	m.BestFitnessRecords = append(m.BestFitnessRecords, bestFitAllIter)
	m.BestSolnRecords = append(m.BestSolnRecords, asmodel.SolutionCopy(bestSolnAllIter))

	return newPopulation
}

// the fitness function of this algorithm
func (m *Mcssga) Fitness(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, chromosome asmodel.Solution) float64 {
	var fitnessValue float64

	m.FitnessNonPriDp = make(map[string]float64) // clear the dp record
	for appName, _ := range apps {
		fitnessValue += m.fitnessOneApp(clouds, apps, chromosome, appName)
	}

	return fitnessValue
}

// calculate the fitness value contributed by an application
func (m *Mcssga) fitnessOneApp(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, chromosome asmodel.Solution, thisAppName string) float64 {
	thisPri := apps[thisAppName].Priority // the fitness values should be weighted by applications' priorities.
	return m.fitnessOneAppNonPri(clouds, apps, chromosome, thisAppName) * float64(thisPri)
}

// calculate the fitness value contributed by an application without the consideration of its priority
func (m *Mcssga) fitnessOneAppNonPri(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, chromosome asmodel.Solution, thisAppName string) float64 {

	var thisAppFitnessNonPri float64 // result

	// if an application is rejected, it contributes a big negative fitness value, this is to encourage higher priority-weighted acceptance rate
	if !chromosome.AppsSolution[thisAppName].Accepted {
		thisAppFitnessNonPri = -(m.ExpAppCompuTimeOneCpu + m.MaxReachableRtt*m.AvgDepNum) / 2
	} else {

		// if this app is accepted, all its dependent apps are also accepted, which is guaranteed by our dependency acceptable check

		// calculate the computation part of this application
		thisAlloCpu := chromosome.AppsSolution[thisAppName].AllocatedCpuCore

		// the maximum possible computation part of fitness of an application without priority should be "m.ExpAppCompuTimeOneCpu"
		thisAppPart := m.ExpAppCompuTimeOneCpu - m.ExpAppCompuTimeOneCpu/thisAlloCpu

		/**
		An application's fitness is only affected by its computation part and the network part to its dependent applications. We do not need to consider its dependent applications' computation parts, because that parts are already considered in the dependent applications' fitness values.
		*/

		netPart := m.MaxReachableRtt * m.AvgDepNum // the base network part of fitness.
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

		thisAppFitnessNonPri = thisAppPart + netPart
	}

	return thisAppFitnessNonPri
}

// Deprecated: this fitness function tend to accept the applications with more dependencies, which is wrong.
// calculate the fitness value contributed by an application without the consideration of its priority
func (m *Mcssga) oldFitnessOneAppNonPri(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, chromosome asmodel.Solution, thisAppName string) float64 {

	// dynamic programming to reduce the scheduling time
	if recordedFitnessNonPri, exist := m.FitnessNonPriDp[thisAppName]; exist {
		return recordedFitnessNonPri
	}

	var thisAppFitnessNonPri float64 // result

	// if an application is rejected, it contributes a very big negative fitness value, this is to encourage higher priority-weighted acceptance rate
	if !chromosome.AppsSolution[thisAppName].Accepted {
		thisAppFitnessNonPri = -(m.ExpAppCompuTimeOneCpu*enlargerScaleMaxCompuTime + m.MaxReachableRtt*enlargerScaleMaxRTT) * (float64(len(apps)) / 5.0)
	} else {

		// if this app is accepted, all its dependent apps are also accepted, which is guaranteed by our dependency acceptable check

		// calculate the computation part of this application
		thisAlloCpu := chromosome.AppsSolution[thisAppName].AllocatedCpuCore
		//thisReqCpu := apps[thisAppName].Resources.CpuCore
		// the maximum possible computation part of fitness of an application without priority should be "m.ExpAppCompuTimeOneCpu"
		thisAppPart := m.ExpAppCompuTimeOneCpu*enlargerScaleMaxCompuTime - m.ExpAppCompuTimeOneCpu/thisAlloCpu

		/**
		If this application does not have dependencies, its fitness will only be affected by itself.
		if this application has dependencies, its fitness will be affected by all its dependent apps and their dependent apps, so this is a recursive function.
		If this application has dependencies, we calculate the sum of the fitness contributed by all dependencies, but do not calculate the average value, because in practice all dependencies should be accessed.
		*/

		var sumAllDeps float64 = 0
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
			netPart := m.MaxReachableRtt*enlargerScaleMaxRTT - thisRtt // because the minimum computation part is not 0, this minimum should also not be 0.

			// calculate the dependent application's fitness value, including its computation part, network part, and dependent apps part.
			depAppPart := m.oldFitnessOneAppNonPri(clouds, apps, chromosome, depAppName)

			// add the 2 parts of the fitness value of this dependency to the sum
			thisDepFitnessNonPri := netPart + depAppPart
			sumAllDeps += thisDepFitnessNonPri
		}

		// For an app without dependencies, sumAllDeps will be 0.
		thisAppFitnessNonPri = thisAppPart + sumAllDeps
	}

	m.FitnessNonPriDp[thisAppName] = thisAppFitnessNonPri
	return thisAppFitnessNonPri
}

// draw m.BestFitnessEachIter and m.BestFitnessRecords on a line chart, to show the evolution trend
func (m *Mcssga) DrawEvoChart() {
	var drawChartFunc func(http.ResponseWriter, *http.Request) = func(res http.ResponseWriter, r *http.Request) {
		var xValuesAllBest []float64
		for i, _ := range m.BestFitnessRecords {
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
					YValues: m.BestFitnessRecords,
				},
				chart.ContinuousSeries{
					Name:    "Best Fitness in each iterations",
					XValues: xValuesAllBest,
					YValues: m.BestFitnessEachIter,
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

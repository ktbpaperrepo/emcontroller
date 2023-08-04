package algorithms

import (
	asmodel "emcontroller/auto-schedule/model"
	"emcontroller/models"
	"fmt"
	"github.com/KeepTheBeats/routing-algorithms/random"
	"github.com/astaxie/beego"
)

/**
NOTE:
In our fitness function, we do not need to consider to make unacceptable chromosomes harder to be selected (such as "fitness /= 3 for un accepted ones"), because:
1. this is complicated;
2. we need to use the allocatedCpuCore to calculate fitness values, but unacceptable solution cannot have this.
Instead, in our init function, mutation selector, crossover selector, we should guarantee that all generated solutions/chromosomes are acceptable. For mutation and crossover if we get an unacceptable solution, we can repeat, excluding the generated unacceptable ones, until we get an acceptable solution.
*/

// Multi-Cloud Service Scheduling Genetic Algorithm (MCSSGA)
type Mcssga struct {
	ChromosomesCount     int // One chromosome is a solution
	IterationCount       int // In each iteration, a population will be generated. One population consists of some solutions.
	CrossoverProbability float64
	MutationProbability  float64

	// If in the past StopNoUpdateIteration iterations, the best solution so far has not been updated, we should end this algorithm and return the best solution so far.
	StopNoUpdateIteration int
	CurNoUpdateIteration  int // record how many iterations the best solution has not updated currently.

	MaxReachableRtt float64 // The biggest RTT between any 2 (or 1) reachable clouds, used to calculate fitness values.

	// these 2 member variables record the best solution in each iteration as well as its fitness value
	BestFitnessRecords []float64
	BestSolnRecords    []asmodel.Solution
}

func NewMcssga(chromosomesCount int, iterationCount int, crossoverProbability float64, mutationProbability float64, stopNoUpdateIteration int) *Mcssga {
	return &Mcssga{
		ChromosomesCount:      chromosomesCount,
		IterationCount:        iterationCount,
		CrossoverProbability:  crossoverProbability,
		MutationProbability:   mutationProbability,
		StopNoUpdateIteration: stopNoUpdateIteration,
		CurNoUpdateIteration:  0,
		MaxReachableRtt:       0,
		BestFitnessRecords:    nil,
		BestSolnRecords:       nil,
	}
}

// Traverse all clouds to find the max RTT between any 2 (or 1) reachable clouds, and set it as the MaxReachableRtt of Mcssga.
func (m *Mcssga) setMaxReaRtt(clouds map[string]asmodel.Cloud) {
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

func (m *Mcssga) Schedule(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, error) {
	m.setMaxReaRtt(clouds)
	beego.Info("MaxReachableRtt:", m.MaxReachableRtt)

	beego.Info("Clouds:")
	for _, cloud := range clouds {
		beego.Info(models.JsonString(cloud))
	}
	beego.Info("Applications:")
	for _, app := range apps {
		beego.Info(models.JsonString(app))
	}

	// randomly generate the init population
	var initPopulation []asmodel.Solution = m.initialize(clouds, apps, appsOrder)
	beego.Info("initPopulation:")
	for _, soln := range initPopulation {
		beego.Info(models.JsonString(soln))
	}

	// there are IterationCount+1 iterations in total, this is the No. 0 iteration
	currentPopulation := m.selectionOperator(clouds, apps, initPopulation) // Iteration No. 0

	// No. 1 iteration to No. m.IterationCount iteration
	for iteration := 1; iteration <= m.IterationCount; iteration++ {
		
		// TODO: crossover and mutation

		currentPopulation = m.selectionOperator(clouds, apps, currentPopulation)

		// If we did not find better solutions in the past some iterations, we stop the algorithm and return the result.
		if m.CurNoUpdateIteration > m.StopNoUpdateIteration {
			break
		}
	}

	beego.Info("Final BestFitnessRecords:", m.BestFitnessRecords)
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

// the selection operator of Genetic Algorithm
func (m *Mcssga) selectionOperator(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, population []asmodel.Solution) []asmodel.Solution {

	// calculate the fitness of every chromosome in the current (old) population
	fitnesses := make([]float64, len(population))
	for i := 0; i < len(population); i++ {
		fitnesses[i] = m.Fitness(clouds, apps, population[i])
	}

	beego.Info("fitness values this iteration:", fitnesses)

	// do selection to generate a new population
	var newPopulation []asmodel.Solution
	pickHelper := make([]int, len(fitnesses)) // for binary tournament selection

	// to record the solution with the highest fitness value in the new population generated in this iteration
	var bestFitThisIter float64 = -maxAccRttMs // initialized with a very small value
	var bestFitThisIterIdx int = 0             // the index of the solution with the highest fitness value in the old population

	// in every population, there should be m.ChromosomesCount chromosomes
	for i := 0; i < m.ChromosomesCount; i++ {
		var selChrIdx int // selected chromosome index in the input population

		// binary tournament selection
		picked := random.RandomPickN(pickHelper, 2)
		if fitnesses[picked[0]] > fitnesses[picked[1]] {
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

	// We only record the best solutions until every iteration.
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
	m.BestFitnessRecords = append(m.BestFitnessRecords, bestFitAllIter)
	m.BestSolnRecords = append(m.BestSolnRecords, asmodel.SolutionCopy(bestSolnAllIter))

	return newPopulation
}

// the fitness function of this algorithm
func (m *Mcssga) Fitness(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, chromosome asmodel.Solution) float64 {
	var fitnessValue float64

	for appName, _ := range apps {
		fitnessValue += m.fitnessOneApp(clouds, apps, chromosome, appName)
	}

	return fitnessValue
}

// calculate the fitness value contributed by an application
func (m *Mcssga) fitnessOneApp(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, chromosome asmodel.Solution, thisAppName string) float64 {
	var thisAppFitness float64 // result

	thisPri := apps[thisAppName].Priority // the fitness values should be weighted by applications' priorities.
	depNum := float64(len(apps[thisAppName].Dependencies))
	if depNum > 0 {
		// If this application has dependencies, we calculate the sum of the fitness contributed by all dependencies, and then calculate the average value.

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
			netPart := m.MaxReachableRtt - thisRtt

			// calculate this application's computation part of the fitness value of this dependency
			thisAlloCpu := chromosome.AppsSolution[thisAppName].AllocatedCpuCore
			thisReqCpu := apps[thisAppName].Resources.CpuCore

			thisAppPart := m.MaxReachableRtt * (thisAlloCpu / thisReqCpu)

			// calculate the dependent application's computation part of the fitness value of this dependency
			depAlloCpu := chromosome.AppsSolution[depAppName].AllocatedCpuCore
			depReqCpu := apps[depAppName].Resources.CpuCore

			depAppPart := m.MaxReachableRtt * (depAlloCpu / depReqCpu)

			// add the 3 parts of the fitness value of this dependency to the sum
			thisDepFitness := netPart + thisAppPart + depAppPart
			sumAllDeps += thisDepFitness

		}

		// weighted by applications' priorities.
		thisAppFitness = sumAllDeps / depNum * float64(thisPri)
	} else {
		// If this application does not have dependencies, its fitness will only be contributed by itself.

		thisAlloCpu := chromosome.AppsSolution[thisAppName].AllocatedCpuCore
		thisReqCpu := apps[thisAppName].Resources.CpuCore

		thisAppPart := m.MaxReachableRtt * (thisAlloCpu / thisReqCpu)

		// weighted by applications' priorities.
		thisAppFitness = thisAppPart * float64(thisPri)
	}

	return thisAppFitness
}

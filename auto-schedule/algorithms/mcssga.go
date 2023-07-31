package algorithms

import (
	"github.com/astaxie/beego"

	asmodel "emcontroller/auto-schedule/model"
)

// Multi-Cloud Service Scheduling Genetic Algorithm (MCSSGA)
type Mcssga struct {
	ChromosomesCount      int // One chromosome is a solution
	IterationCount        int // In each iteration, a population will be generated. One population consists of some solutions.
	CrossoverProbability  float64
	MutationProbability   float64
	StopNoUpdateIteration int
}

func NewMcssga(chromosomesCount int, iterationCount int, crossoverProbability float64, mutationProbability float64, stopNoUpdateIteration int) *Mcssga {
	return &Mcssga{
		ChromosomesCount:      chromosomesCount,
		IterationCount:        iterationCount,
		CrossoverProbability:  crossoverProbability,
		MutationProbability:   mutationProbability,
		StopNoUpdateIteration: stopNoUpdateIteration,
	}
}

func (m *Mcssga) Schedule(clouds map[string]asmodel.Cloud, apps map[string]asmodel.Application, appsOrder []string) (asmodel.Solution, error) {
	beego.Info("Clouds:")
	for _, cloud := range clouds {
		beego.Info(asmodel.JsonString(cloud))
	}
	beego.Info("Applications:")
	for _, app := range apps {
		beego.Info(asmodel.JsonString(app))
	}

	// randomly generate the init population
	var initPopulation []asmodel.Solution = m.initialize(clouds, apps, appsOrder)
	beego.Info("initPopulation:")
	for _, soln := range initPopulation {
		beego.Info(asmodel.JsonString(soln))
	}

	return asmodel.Solution{}, nil
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

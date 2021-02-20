package evolve

import (
	"fmt"

	cxcore "github.com/skycoin/cx/cx"
)

func (pop *Population) Evolve() {
	errors := make([]float64, pop.PopulationSize)
	numIter := pop.Iterations
	solProt := pop.FunctionToEvolve
	fnToEvolveName := solProt.Name
	sPrgrm := cxcore.Serialize(pop.Individuals[0], 0)
	targetError := pop.TargetError
	inputs := pop.Inputs
	outputs := pop.Outputs

	// Evolution process.
	for c := 0; c < int(numIter); c++ {
		// Selection process.
		pop1Idx, pop2Idx := tournamentSelection(errors, 0.5, true)
		dead1Idx, dead2Idx := tournamentSelection(errors, 0.5, false)

		pop1MainPkg, err := pop.Individuals[pop1Idx].GetPackage(cxcore.MAIN_PKG)
		if err != nil {
			panic(err)
		}
		parent1, err := pop1MainPkg.GetFunction(fnToEvolveName)
		if err != nil {
			panic(err)
		}

		pop2MainPkg, err := pop.Individuals[pop2Idx].GetPackage(cxcore.MAIN_PKG)
		if err != nil {
			panic(err)
		}
		parent2, err := pop2MainPkg.GetFunction(fnToEvolveName)
		if err != nil {
			panic(err)
		}

		// Crossover process.
		crossoverFn := pop.getCrossoverFn()
		child1, child2 := crossoverFn(parent1, parent2)
		// child1 := parent1
		// child2 := parent2

		// Mutation process.
		_ = sPrgrm
		_ = dead1Idx
		_ = dead2Idx
		_ = child1
		_ = child2
		randomMutation(pop, sPrgrm)

		// Replacing individuals in population.
		replaceSolution(pop.Individuals[dead1Idx], fnToEvolveName, child1)
		replaceSolution(pop.Individuals[dead2Idx], fnToEvolveName, child2)

		// Evaluation process.
		for i, _ := range pop.Individuals {
			errors[i] = perByteEvaluation(pop.Individuals[i], solProt, inputs, outputs)
			if errors[i] <= targetError {
				fmt.Printf("\nFound solution:\n\n")
				pop.Individuals[i].PrintProgram()
				return
			}
		}

		avg := 0.0
		for _, err := range errors {
			avg += err
		}
		fmt.Printf("%v\n", float64(avg) / float64(len(errors)))
	}
}

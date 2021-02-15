package main

import (
	"os"
	evolve "github.com/skycoin/cx-evolves/evolve"
	actions "github.com/skycoin/cx/cxgo/actions"
	cxgo "github.com/skycoin/cx/cxgo/cxgo"
	cxcore "github.com/skycoin/cx/cx"
)

func InitEvolve() {
	// Registering this library as a CX library.
	cxcore.RegisterPackage("evolve")
	cxcore.Op(cxcore.GetOpCodeCount(), "evolve.evolve", evolve.OpEvolve, cxcore.In(cxcore.Slice(cxcore.TYPE_AFF), cxcore.Slice(cxcore.TYPE_AFF), cxcore.Slice(cxcore.TYPE_AFF), cxcore.Slice(cxcore.TYPE_AFF), cxcore.Slice(cxcore.TYPE_AFF), cxcore.AI32, cxcore.AI32, cxcore.AI32, cxcore.AF64), nil)
}

func RunEvolve() {
	// Creating a new CX program.
	actions.PRGRM = cxcore.MakeProgram()

	// Reading flags.
	options := evolve.DefaultCmdFlags()
	evolve.ParseFlags(&options, os.Args[1:])
	args := evolve.CommandLine.Args()

	// Reading source code and parsing it to a valid CX program.
	_, sourceCode, fileNames := cxcore.ParseArgsForCX(args, true)
	cxgo.ParseSourceCode(sourceCode, fileNames)
	cxgo.AddInitFunction(actions.PRGRM)

	actions.PRGRM.RunCompiled(0, nil)
}

func main() {
	InitEvolve()
	RunEvolve()
}

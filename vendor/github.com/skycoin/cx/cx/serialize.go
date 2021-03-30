package cxcore

import (
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

type sIndex struct {
	ProgramOffset     int32
	CallsOffset       int32
	PackagesOffset    int32
	StructsOffset     int32
	FunctionsOffset   int32
	ExpressionsOffset int32
	ArgumentsOffset   int32
	IntegersOffset    int32
	NamesOffset       int32
	MemoryOffset      int32
}

type sProgram struct {
	PackagesOffset       int32
	PackagesSize         int32
	CurrentPackageOffset int32

	InputsOffset int32
	InputsSize   int32

	OutputsOffset int32
	OutputsSize   int32

	CallStackOffset int32
	CallStackSize   int32

	CallCounter int32

	MemoryOffset int32
	MemorySize   int32

	HeapPointer  int32
	StackPointer int32
	StackSize    int32
	HeapSize     int32
	HeapStartsAt int32

	Terminated int32

	BCPackageCount int32

	VersionOffset int32
	VersionSize   int32
}

type sCall struct {
	OperatorOffset int32
	Line           int32
	FramePointer   int32
}

type sPackage struct {
	NameOffset            int32
	NameSize              int32
	ImportsOffset         int32
	ImportsSize           int32
	StructsOffset         int32
	StructsSize           int32
	GlobalsOffset         int32
	GlobalsSize           int32
	FunctionsOffset       int32
	FunctionsSize         int32
	CurrentFunctionOffset int32
	CurrentStructOffset   int32
}

type sStruct struct {
	NameOffset   int32
	NameSize     int32
	FieldsOffset int32
	FieldsSize   int32

	Size int32

	PackageOffset int32
}

type sFunction struct {
	NameOffset        int32
	NameSize          int32
	InputsOffset      int32
	InputsSize        int32
	OutputsOffset     int32
	OutputsSize       int32
	ExpressionsOffset int32
	ExpressionsSize   int32
	Size              int32
	Length            int32

	ListOfPointersOffset int32
	ListOfPointersSize   int32

	// We're going to determine this when procesing the expressions. Check sExpression type
	// IsNative                        int32
	// OpCode                          int32

	CurrentExpressionOffset int32
	PackageOffset           int32
}

type sExpression struct {
	OperatorOffset int32
	// we add these two fields here so we don't add every native sFunction to the serialization
	// the CX runtime already knows about the natives properties. We just need the code if IsNative = true
	IsNative int32
	OpCode   int32

	InputsOffset  int32
	InputsSize    int32
	OutputsOffset int32
	OutputsSize   int32

	LabelOffset int32
	LabelSize   int32
	ThenLines   int32
	ElseLines   int32

	ScopeOperation int32

	IsMethodCall    int32
	IsStructLiteral int32
	IsArrayLiteral  int32
	IsUndType       int32
	IsBreak         int32
	IsContinue      int32

	FunctionOffset int32
	PackageOffset  int32
}

type sArgument struct {
	NameOffset       int32
	NameSize         int32
	Type             int32
	CustomTypeOffset int32
	Size             int32
	TotalSize        int32

	Offset int32

	IndirectionLevels           int32
	DereferenceLevels           int32
	DereferenceOperationsOffset int32
	DereferenceOperationsSize   int32
	DeclarationSpecifiersOffset int32
	DeclarationSpecifiersSize   int32

	IsSlice      int32
	IsArray      int32
	IsArrayFirst int32
	IsPointer    int32
	IsReference  int32

	IsDereferenceFirst int32
	IsStruct           int32
	IsRest             int32
	IsLocalDeclaration int32
	IsShortDeclaration int32
	PreviouslyDeclared int32

	PassBy     int32
	DoesEscape int32

	LengthsOffset int32
	LengthsSize   int32
	IndexesOffset int32
	IndexesSize   int32
	FieldsOffset  int32
	FieldsSize    int32
	InputsOffset  int32
	InputsSize    int32
	OutputsOffset int32
	OutputsSize   int32

	PackageOffset int32
}

type sAll struct {
	Index   sIndex
	Program sProgram

	Packages     []sPackage
	PackagesMap  map[string]int
	Structs      []sStruct
	StructsMap   map[string]int
	Functions    []sFunction
	FunctionsMap map[string]int

	Expressions []sExpression
	Arguments   []sArgument
	Calls       []sCall

	Names    []byte
	NamesMap map[string]int
	Integers []int32

	Memory []byte
}

func serializeName(name string, s *sAll) (int32, int32) {
	if name == "" {
		return int32(-1), int32(-1)
	}

	size := encoder.Size(name)

	off, found := s.NamesMap[name]
	if found {
		return int32(off), int32(size)
	}
	off = len(s.Names)
	s.Names = append(s.Names, encoder.Serialize(name)...)
	s.NamesMap[name] = off

	return int32(off), int32(size)
}

func indexPackage(pkg *CXPackage, s *sAll) {
	if _, found := s.PackagesMap[pkg.Name]; !found {
		s.PackagesMap[pkg.Name] = len(s.PackagesMap)
	} else {
		panic("duplicated package in serialization process")
	}
}

func indexStruct(strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if _, found := s.StructsMap[strctName]; !found {
		s.StructsMap[strctName] = len(s.StructsMap)
	} else {
		panic("duplicated struct in serialization process")
	}
}

func indexFunction(fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if _, found := s.FunctionsMap[fnName]; !found {
		s.FunctionsMap[fnName] = len(s.FunctionsMap)
	} else {
		panic("duplicated function in serialization process")
	}
}

func serializeBoolean(val bool) int32 {
	if val {
		return 1
	}
	return 0
}

func serializeIntegers(ints []int, s *sAll) (int32, int32) {
	if len(ints) == 0 {
		return int32(-1), int32(-1)
	}
	off := len(s.Integers)
	l := len(ints)

	ints32 := make([]int32, l)
	for i, int := range ints {
		ints32[i] = int32(int)
	}

	s.Integers = append(s.Integers, ints32...)

	return int32(off), int32(l)
}

func serializeArgument(arg *CXArgument, s *sAll) int {
	s.Arguments = append(s.Arguments, sArgument{})
	argOff := len(s.Arguments) - 1

	sNil := int32(-1)

	s.Arguments[argOff].NameOffset, s.Arguments[argOff].NameSize = serializeName(arg.Name, s)

	s.Arguments[argOff].Type = int32(arg.Type)

	if arg.CustomType == nil {
		s.Arguments[argOff].CustomTypeOffset = sNil
	} else {
		strctName := arg.CustomType.Package.Name + "." + arg.CustomType.Name
		if strctOff, found := s.StructsMap[strctName]; found {
			s.Arguments[argOff].CustomTypeOffset = int32(strctOff)
		} else {
			panic("struct reference not found")
		}
	}

	s.Arguments[argOff].Size = int32(arg.Size)
	s.Arguments[argOff].TotalSize = int32(arg.TotalSize)
	s.Arguments[argOff].Offset = int32(arg.Offset)
	s.Arguments[argOff].IndirectionLevels = int32(arg.IndirectionLevels)
	s.Arguments[argOff].DereferenceLevels = int32(arg.DereferenceLevels)

	s.Arguments[argOff].DereferenceOperationsOffset,
		s.Arguments[argOff].DereferenceOperationsSize = serializeIntegers(arg.DereferenceOperations, s)

	s.Arguments[argOff].DeclarationSpecifiersOffset,
		s.Arguments[argOff].DeclarationSpecifiersSize = serializeIntegers(arg.DeclarationSpecifiers, s)

	s.Arguments[argOff].IsSlice = serializeBoolean(arg.IsSlice)
	s.Arguments[argOff].IsArray = serializeBoolean(arg.IsArray)
	s.Arguments[argOff].IsArrayFirst = serializeBoolean(arg.IsArrayFirst)
	s.Arguments[argOff].IsPointer = serializeBoolean(arg.IsPointer)
	s.Arguments[argOff].IsReference = serializeBoolean(arg.IsReference)

	s.Arguments[argOff].IsDereferenceFirst = serializeBoolean(arg.IsDereferenceFirst)
	s.Arguments[argOff].IsStruct = serializeBoolean(arg.IsStruct)
	s.Arguments[argOff].IsRest = serializeBoolean(arg.IsRest)
	s.Arguments[argOff].IsLocalDeclaration = serializeBoolean(arg.IsLocalDeclaration)
	s.Arguments[argOff].IsShortDeclaration = serializeBoolean(arg.IsShortDeclaration)
	s.Arguments[argOff].PreviouslyDeclared = serializeBoolean(arg.PreviouslyDeclared)

	s.Arguments[argOff].PassBy = int32(arg.PassBy)
	s.Arguments[argOff].DoesEscape = serializeBoolean(arg.DoesEscape)

	s.Arguments[argOff].LengthsOffset, s.Arguments[argOff].LengthsSize = serializeIntegers(arg.Lengths, s)
	s.Arguments[argOff].IndexesOffset, s.Arguments[argOff].IndexesSize = serializeSliceOfArguments(arg.Indexes, s)
	s.Arguments[argOff].FieldsOffset, s.Arguments[argOff].FieldsSize = serializeSliceOfArguments(arg.Fields, s)
	s.Arguments[argOff].InputsOffset, s.Arguments[argOff].InputsSize = serializeSliceOfArguments(arg.Inputs, s)
	s.Arguments[argOff].OutputsOffset, s.Arguments[argOff].OutputsSize = serializeSliceOfArguments(arg.Outputs, s)

	if pkgOff, found := s.PackagesMap[arg.Package.Name]; found {
		s.Arguments[argOff].PackageOffset = int32(pkgOff)
	} else {
		panic("package reference not found")
	}

	return argOff
}

func serializeSliceOfArguments(args []*CXArgument, s *sAll) (int32, int32) {
	if len(args) == 0 {
		return int32(-1), int32(-1)
	}
	idxs := make([]int, len(args))
	for i, arg := range args {
		idxs[i] = serializeArgument(arg, s)
	}
	return serializeIntegers(idxs, s)
}

func serializeCalls(calls []CXCall, s *sAll) (int32, int32) {
	if len(calls) == 0 {
		return int32(-1), int32(-1)
	}
	idxs := make([]int, len(calls))
	for i, call := range calls {
		idxs[i] = serializeCall(&call, s)
	}
	return serializeIntegers(idxs, s)

}

func serializeExpression(expr *CXExpression, s *sAll) int {
	s.Expressions = append(s.Expressions, sExpression{})
	exprOff := len(s.Expressions) - 1
	sExpr := &s.Expressions[exprOff]

	sNil := int32(-1)

	if expr.Operator == nil {
		// then it's a declaration
		sExpr.OperatorOffset = sNil
		sExpr.IsNative = serializeBoolean(false)
		sExpr.OpCode = int32(-1)
	} else if expr.Operator.IsNative {
		sExpr.OperatorOffset = sNil
		sExpr.IsNative = serializeBoolean(true)
		sExpr.OpCode = int32(expr.Operator.OpCode)
	} else {
		sExpr.IsNative = serializeBoolean(false)
		sExpr.OpCode = sNil

		opName := expr.Operator.Package.Name + "." + expr.Operator.Name
		if opOff, found := s.FunctionsMap[opName]; found {
			sExpr.OperatorOffset = int32(opOff)
		}
	}

	sExpr.InputsOffset, sExpr.InputsSize = serializeSliceOfArguments(expr.Inputs, s)
	sExpr.OutputsOffset, sExpr.OutputsSize = serializeSliceOfArguments(expr.Outputs, s)

	sExpr.LabelOffset, sExpr.LabelSize = serializeName(expr.Label, s)
	sExpr.ThenLines = int32(expr.ThenLines)
	sExpr.ElseLines = int32(expr.ElseLines)
	sExpr.ScopeOperation = int32(expr.ScopeOperation)

	sExpr.IsMethodCall = serializeBoolean(expr.IsMethodCall)
	sExpr.IsStructLiteral = serializeBoolean(expr.IsStructLiteral)
	sExpr.IsArrayLiteral = serializeBoolean(expr.IsArrayLiteral)
	sExpr.IsUndType = serializeBoolean(expr.IsUndType)
	sExpr.IsBreak = serializeBoolean(expr.IsBreak)
	sExpr.IsContinue = serializeBoolean(expr.IsContinue)

	fnName := expr.Function.Package.Name + "." + expr.Function.Name
	if fnOff, found := s.FunctionsMap[fnName]; found {
		sExpr.FunctionOffset = int32(fnOff)
	} else {
		panic("function reference not found")
	}

	if pkgOff, found := s.PackagesMap[expr.Package.Name]; found {
		sExpr.PackageOffset = int32(pkgOff)
	} else {
		panic("package reference not found")
	}

	return exprOff
}

func serializeCall(call *CXCall, s *sAll) int {
	s.Calls = append(s.Calls, sCall{})
	callOff := len(s.Calls) - 1
	sCall := &s.Calls[callOff]

	opName := call.Operator.Package.Name + "." + call.Operator.Name
	if opOff, found := s.FunctionsMap[opName]; found {
		sCall.OperatorOffset = int32(opOff)
		sCall.Line = int32(call.Line)
		sCall.FramePointer = int32(call.FramePointer)
	} else {
		panic("function reference not found")
	}

	return callOff
}

func serializeProgram(prgrm *CXProgram, s *sAll) {
	s.Program = sProgram{}
	sPrgrm := &s.Program
	sPrgrm.PackagesOffset = int32(0)
	sPrgrm.PackagesSize = int32(len(prgrm.Packages))

	if pkgOff, found := s.PackagesMap[prgrm.CurrentPackage.Name]; found {
		sPrgrm.CurrentPackageOffset = int32(pkgOff)
	} else {
		panic("package reference not found")
	}

	sPrgrm.InputsOffset, sPrgrm.InputsSize = serializeSliceOfArguments(prgrm.Inputs, s)
	sPrgrm.OutputsOffset, sPrgrm.OutputsSize = serializeSliceOfArguments(prgrm.Outputs, s)

	sPrgrm.CallStackOffset, sPrgrm.CallStackSize = serializeCalls(prgrm.CallStack[:prgrm.CallCounter], s)

	sPrgrm.CallCounter = int32(prgrm.CallCounter)

	sPrgrm.MemoryOffset = int32(0)
	sPrgrm.MemorySize = int32(len(PROGRAM.Memory))

	sPrgrm.HeapPointer = int32(prgrm.HeapPointer)
	sPrgrm.StackPointer = int32(prgrm.StackPointer)
	sPrgrm.StackSize = int32(prgrm.StackSize)
	sPrgrm.HeapSize = int32(prgrm.HeapSize)
	sPrgrm.HeapStartsAt = int32(prgrm.HeapStartsAt)

	sPrgrm.Terminated = serializeBoolean(prgrm.Terminated)
	sPrgrm.VersionOffset, sPrgrm.VersionSize = serializeName(prgrm.Version, s)
}

func sStructArguments(strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if strctOff, found := s.StructsMap[strctName]; found {
		sStrct := &s.Structs[strctOff]
		sStrct.FieldsOffset, sStrct.FieldsSize = serializeSliceOfArguments(strct.Fields, s)
	} else {
		panic("struct reference not found")
	}
}

func sFunctionArguments(fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if fnOff, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[fnOff]

		sFn.InputsOffset, sFn.InputsSize = serializeSliceOfArguments(fn.Inputs, s)
		sFn.OutputsOffset, sFn.OutputsSize = serializeSliceOfArguments(fn.Outputs, s)
		sFn.ListOfPointersOffset, sFn.ListOfPointersSize = serializeSliceOfArguments(fn.ListOfPointers, s)
	} else {
		panic("function reference not found")
	}
}

func sPackageName(pkg *CXPackage, s *sAll) {
	sPkg := &s.Packages[s.PackagesMap[pkg.Name]]
	sPkg.NameOffset, sPkg.NameSize = serializeName(pkg.Name, s)
}

func sStructName(strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	sStrct := &s.Structs[s.StructsMap[strctName]]
	sStrct.NameOffset, sStrct.NameSize = serializeName(strct.Name, s)
}

func sFunctionName(fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if off, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[off]
		sFn.NameOffset, sFn.NameSize = serializeName(fn.Name, s)
	} else {
		panic("function reference not found")
	}
}

func sPackageGlobals(pkg *CXPackage, s *sAll) {
	if pkgOff, found := s.PackagesMap[pkg.Name]; found {
		sPkg := &s.Packages[pkgOff]
		sPkg.GlobalsOffset, sPkg.GlobalsSize = serializeSliceOfArguments(pkg.Globals, s)
	} else {
		panic("package reference not found")
	}
}

func sPackageImports(pkg *CXPackage, s *sAll) {
	l := len(pkg.Imports)
	if l == 0 {
		s.Packages[s.PackagesMap[pkg.Name]].ImportsOffset = int32(-1)
		s.Packages[s.PackagesMap[pkg.Name]].ImportsSize = int32(-1)
		return
	}
	imps := make([]int32, l)
	for i, imp := range pkg.Imports {
		if idx, found := s.PackagesMap[imp.Name]; found {
			imps[i] = int32(idx)
		} else {
			panic("import package reference not found")
		}
	}

	s.Packages[s.PackagesMap[pkg.Name]].ImportsOffset = int32(len(s.Integers))
	s.Packages[s.PackagesMap[pkg.Name]].ImportsSize = int32(l)
	s.Integers = append(s.Integers, imps...)
}

func sStructPackage(strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if pkgOff, found := s.PackagesMap[strct.Package.Name]; found {
		if off, found := s.StructsMap[strctName]; found {
			sStrct := &s.Structs[off]
			sStrct.PackageOffset = int32(pkgOff)
		} else {
			panic("struct reference not found")
		}
	} else {
		panic("struct's package reference not found")
	}
}

func sFunctionPackage(fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if pkgOff, found := s.PackagesMap[fn.Package.Name]; found {
		if off, found := s.FunctionsMap[fnName]; found {
			sFn := &s.Functions[off]
			sFn.PackageOffset = int32(pkgOff)
		} else {
			panic("function reference not found")
		}
	} else {
		panic("function's package reference not found")
	}
}

func sPackageIntegers(pkg *CXPackage, s *sAll) {
	if pkgOff, found := s.PackagesMap[pkg.Name]; found {
		sPkg := &s.Packages[pkgOff]

		if pkg.CurrentFunction == nil {
			// package has no functions
			sPkg.CurrentFunctionOffset = int32(-1)
		} else {
			currFnName := pkg.CurrentFunction.Package.Name + "." + pkg.CurrentFunction.Name

			if fnOff, found := s.FunctionsMap[currFnName]; found {
				sPkg.CurrentFunctionOffset = int32(fnOff)
			} else {
				panic("function reference not found")
			}
		}

		if pkg.CurrentStruct == nil {
			// package has no structs
			sPkg.CurrentStructOffset = int32(-1)
		} else {
			currStrctName := pkg.CurrentStruct.Package.Name + "." + pkg.CurrentStruct.Name

			if strctOff, found := s.StructsMap[currStrctName]; found {
				sPkg.CurrentStructOffset = int32(strctOff)
			} else {
				panic("struct reference not found")
			}
		}
	} else {
		panic("package reference not found")
	}
}

func sStructIntegers(strct *CXStruct, s *sAll) {
	strctName := strct.Package.Name + "." + strct.Name
	if off, found := s.StructsMap[strctName]; found {
		sStrct := &s.Structs[off]
		sStrct.Size = int32(strct.Size)
	} else {
		panic("struct reference not found")
	}
}

func sFunctionIntegers(fn *CXFunction, s *sAll) {
	fnName := fn.Package.Name + "." + fn.Name
	if off, found := s.FunctionsMap[fnName]; found {
		sFn := &s.Functions[off]
		sFn.Size = int32(fn.Size)
		sFn.Length = int32(fn.Length)
	} else {
		panic("function reference not found")
	}
}

func initSerialization(prgrm *CXProgram, s *sAll) {
	s.PackagesMap = make(map[string]int)
	s.StructsMap = make(map[string]int)
	s.FunctionsMap = make(map[string]int)
	s.NamesMap = make(map[string]int)

	s.Calls = make([]sCall, prgrm.CallCounter)
	s.Packages = make([]sPackage, len(prgrm.Packages))

	// s.Memory = prgrm.Memory[:PROGRAM.HeapStartsAt+PROGRAM.HeapPointer]
	s.Memory = prgrm.Memory

	var numStrcts int
	var numFns int

	for _, pkg := range prgrm.Packages {
		numStrcts += len(pkg.Structs)
		numFns += len(pkg.Functions)
	}

	s.Structs = make([]sStruct, numStrcts)
	s.Functions = make([]sFunction, numFns)
	// args and exprs need to be appended as they are found
}

// SplitSerialize ...
// WHAT DOES THIS DO? WHY ARE THERE NO COMMENTS?
func splitSerialize(prgrm *CXProgram, s *sAll, fnCounter, strctCounter *int32, from, to int) {
	// indexing packages and serializing their names
	for _, pkg := range prgrm.Packages[from:to] {
		indexPackage(pkg, s)
		sPackageName(pkg, s)
	}
	// we first needed to populate references to all packages
	// now we add the imports' references
	for _, pkg := range prgrm.Packages[from:to] {
		sPackageImports(pkg, s)
	}

	// structs
	for _, pkg := range prgrm.Packages[from:to] {
		for _, strct := range pkg.Structs {
			indexStruct(strct, s)
			sStructName(strct, s)
			sStructPackage(strct, s)
			sStructIntegers(strct, s)
		}
	}
	// we first needed to populate references to all structs
	// now we add fields
	for _, pkg := range prgrm.Packages[from:to] {
		for _, strct := range pkg.Structs {
			sStructArguments(strct, s)
		}
	}

	// globals
	for _, pkg := range prgrm.Packages[from:to] {
		sPackageGlobals(pkg, s)
	}

	// functions
	for _, pkg := range prgrm.Packages[from:to] {
		for _, fn := range pkg.Functions {
			indexFunction(fn, s)
			sFunctionName(fn, s)
			sFunctionPackage(fn, s)
			sFunctionIntegers(fn, s)
			sFunctionArguments(fn, s)
		}
	}

	// package elements' offsets and sizes
	for _, pkg := range prgrm.Packages[from:to] {
		if pkgOff, found := s.PackagesMap[pkg.Name]; found {
			sPkg := &s.Packages[pkgOff]

			if len(pkg.Structs) == 0 {
				sPkg.StructsOffset = int32(-1)
				sPkg.StructsSize = int32(-1)
			} else {
				sPkg.StructsOffset = *strctCounter
				lenStrcts := int32(len(pkg.Structs))
				sPkg.StructsSize = lenStrcts
				*strctCounter += lenStrcts
			}

			if len(pkg.Functions) == 0 {
				sPkg.FunctionsOffset = int32(-1)
				sPkg.FunctionsSize = int32(-1)
			} else {
				sPkg.FunctionsOffset = *fnCounter
				lenFns := int32(len(pkg.Functions))
				sPkg.FunctionsSize = lenFns
				*fnCounter += lenFns
			}
		} else {
			panic("package reference not found")
		}
	}

	// package integers
	// we needed the references to all functions and structs first
	for _, pkg := range prgrm.Packages[from:to] {
		sPackageIntegers(pkg, s)
	}

	// expressions
	for _, pkg := range prgrm.Packages[from:to] {
		for _, fn := range pkg.Functions {
			fnName := fn.Package.Name + "." + fn.Name
			if fnOff, found := s.FunctionsMap[fnName]; found {
				sFn := &s.Functions[fnOff]

				if len(fn.Expressions) == 0 {
					sFn.ExpressionsOffset = int32(-1)
					sFn.ExpressionsSize = int32(-1)
					sFn.CurrentExpressionOffset = int32(-1)
				} else {
					exprs := make([]int, len(fn.Expressions))
					for i, expr := range fn.Expressions {
						exprIdx := serializeExpression(expr, s)
						if fn.CurrentExpression == expr {
							// sFn.CurrentExpressionOffset = int32(exprIdx)
							sFn.CurrentExpressionOffset = int32(i)
						}
						exprs[i] = exprIdx
					}

					sFn.ExpressionsOffset, sFn.ExpressionsSize = serializeIntegers(exprs, s)
				}
			} else {
				panic("function reference not found")
			}
		}
	}
}

// Serialize ...
func Serialize(prgrm *CXProgram, split int) (byts []byte) {
	// prgrm.PrintProgram()

	s := sAll{}
	initSerialization(prgrm, &s)

	var fnCounter int32
	var strctCounter int32
	splitSerialize(prgrm, &s, &fnCounter, &strctCounter, 0, split)
	splitSerialize(prgrm, &s, &fnCounter, &strctCounter, split, len(prgrm.Packages))

	// program
	serializeProgram(prgrm, &s)

	s.Index = sIndex{}
	sIdx := &s.Index

	// assigning relative offset

	idxSize := encoder.Size(s.Index)
	prgrmSize := encoder.Size(s.Program)
	callSize := encoder.Size(s.Calls)
	pkgSize := encoder.Size(s.Packages)
	strctSize := encoder.Size(s.Structs)
	fnSize := encoder.Size(s.Functions)
	exprSize := encoder.Size(s.Expressions)
	argSize := encoder.Size(s.Arguments)
	intSize := encoder.Size(s.Integers)

	// assigning absolute offset
	sIdx.ProgramOffset += int32(idxSize)
	sIdx.CallsOffset += sIdx.ProgramOffset + int32(prgrmSize)
	sIdx.PackagesOffset += sIdx.CallsOffset + int32(callSize)
	sIdx.StructsOffset += sIdx.PackagesOffset + int32(pkgSize)
	sIdx.FunctionsOffset += sIdx.StructsOffset + int32(strctSize)
	sIdx.ExpressionsOffset += sIdx.FunctionsOffset + int32(fnSize)
	sIdx.ArgumentsOffset += sIdx.ExpressionsOffset + int32(exprSize)
	sIdx.IntegersOffset += sIdx.ArgumentsOffset + int32(argSize)
	sIdx.NamesOffset += sIdx.IntegersOffset + int32(intSize)
	sIdx.MemoryOffset += sIdx.NamesOffset + int32(len(s.Names))

	// serializing everything
	byts = append(byts, encoder.Serialize(s.Index)...)
	byts = append(byts, encoder.Serialize(s.Program)...)
	byts = append(byts, encoder.Serialize(s.Calls)...)
	byts = append(byts, encoder.Serialize(s.Packages)...)
	byts = append(byts, encoder.Serialize(s.Structs)...)
	byts = append(byts, encoder.Serialize(s.Functions)...)
	byts = append(byts, encoder.Serialize(s.Expressions)...)
	byts = append(byts, encoder.Serialize(s.Arguments)...)
	byts = append(byts, encoder.Serialize(s.Integers)...)
	byts = append(byts, s.Names...)
	byts = append(byts, s.Memory...)

	return byts
}

func opSerialize(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	_ = inp1

	var slcOff int
	byts := Serialize(PROGRAM, 0)
	for _, b := range byts {
		slcOff = WriteToSlice(slcOff, []byte{b})
	}

	WriteI32(out1Offset, int32(slcOff))
}

func opDeserialize(expr *CXExpression, fp int) {
	inp := expr.Inputs[0]

	inpOffset := GetFinalOffset(fp, inp)

	off := Deserialize_i32(PROGRAM.Memory[inpOffset : inpOffset+TYPE_POINTER_SIZE])

	_l := PROGRAM.Memory[off+OBJECT_HEADER_SIZE : off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE]
	l := Deserialize_i32(_l[4:8])

	Deserialize(PROGRAM.Memory[off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE : off+OBJECT_HEADER_SIZE+SLICE_HEADER_SIZE+l]) // BUG : should be l * elt.TotalSize ?
}

func dsName(off int32, size int32, s *sAll) string {
	if size < 1 {
		return ""
	}

	var name string
	DeserializeRaw(s.Names[off:off+size], &name)

	return name
}

func dsPackages(s *sAll, prgrm *CXProgram) {
	var fnCounter int32
	var strctCounter int32

	for i, sPkg := range s.Packages {
		// initializing packages with their names,
		// empty functions, structs, imports and globals
		// and current function and struct
		pkg := CXPackage{}
		prgrm.Packages[i] = &pkg

		pkg.Name = dsName(sPkg.NameOffset, sPkg.NameSize, s)

		if sPkg.ImportsSize > 0 {
			prgrm.Packages[i].Imports = make([]*CXPackage, sPkg.ImportsSize)
		}

		if sPkg.FunctionsSize > 0 {
			prgrm.Packages[i].Functions = make([]*CXFunction, sPkg.FunctionsSize)

			for j, sFn := range s.Functions[sPkg.FunctionsOffset : sPkg.FunctionsOffset+sPkg.FunctionsSize] {
				var fn CXFunction
				fn.Name = dsName(sFn.NameOffset, sFn.NameSize, s)
				prgrm.Packages[i].Functions[j] = &fn
			}
		}

		if sPkg.StructsSize > 0 {
			prgrm.Packages[i].Structs = make([]*CXStruct, sPkg.StructsSize)

			for j, sStrct := range s.Structs[sPkg.StructsOffset : sPkg.StructsOffset+sPkg.StructsSize] {
				var strct CXStruct
				strct.Name = dsName(sStrct.NameOffset, sStrct.NameSize, s)
				prgrm.Packages[i].Structs[j] = &strct
			}
		}

		if sPkg.GlobalsSize > 0 {
			prgrm.Packages[i].Globals = make([]*CXArgument, sPkg.GlobalsSize)
		}

		// // CurrentFunction
		// if sPkg.FunctionsSize > 0 {
		// 	prgrm.Packages[i].CurrentFunction = prgrm.Packages[i].Functions[sPkg.CurrentFunctionOffset-fnCounter]
		// }

		// CurrentStruct
		if sPkg.StructsSize > 0 {
			prgrm.Packages[i].CurrentStruct = prgrm.Packages[i].Structs[sPkg.CurrentStructOffset-strctCounter]
		}

		fnCounter += sPkg.FunctionsSize
		strctCounter += sPkg.StructsSize
	}

	// imports
	for i, sPkg := range s.Packages {
		if sPkg.ImportsSize > 0 {
			// getting indexes of imports
			idxs := dsIntegers(sPkg.ImportsOffset, sPkg.ImportsSize, s)

			for j, idx := range idxs {
				prgrm.Packages[i].Imports[j] = getImport(&s.Packages[idx], s, prgrm)
			}
		}
	}

	// globals
	for i, sPkg := range s.Packages {
		if sPkg.GlobalsSize > 0 {
			prgrm.Packages[i].Globals = dsArguments(sPkg.GlobalsOffset, sPkg.GlobalsSize, s, prgrm)
		}
	}

	// structs
	for i, sPkg := range s.Packages {
		if sPkg.StructsSize > 0 {
			for j, sStrct := range s.Structs[sPkg.StructsOffset : sPkg.StructsOffset+sPkg.StructsSize] {
				dsStruct(&sStrct, prgrm.Packages[i].Structs[j], s, prgrm)
			}
		}
	}

	// functions
	for i, sPkg := range s.Packages {
		if sPkg.FunctionsSize > 0 {
			for j, sFn := range s.Functions[sPkg.FunctionsOffset : sPkg.FunctionsOffset+sPkg.FunctionsSize] {
				dsFunction(&sFn, prgrm.Packages[i].Functions[j], s, prgrm)
			}
		}
	}

	// current package
	prgrm.CurrentPackage = prgrm.Packages[s.Program.CurrentPackageOffset]
}

func dsStruct(sStrct *sStruct, strct *CXStruct, s *sAll, prgrm *CXProgram) {
	strct.Name = dsName(sStrct.NameOffset, sStrct.NameSize, s)
	strct.Fields = dsArguments(sStrct.FieldsOffset, sStrct.FieldsSize, s, prgrm)
	strct.Size = int(sStrct.Size)
	strct.Package = prgrm.Packages[sStrct.PackageOffset]
}

func dsArguments(off int32, size int32, s *sAll, prgrm *CXProgram) []*CXArgument {
	if size < 1 {
		return nil
	}

	// getting indexes of arguments
	idxs := dsIntegers(off, size, s)

	// sArgs := s.Arguments[off : off + size]
	args := make([]*CXArgument, size)
	for i, idx := range idxs {
		args[i] = dsArgument(&s.Arguments[idx], s, prgrm)
	}
	return args
}

func getCustomType(sArg *sArgument, s *sAll, prgrm *CXProgram) *CXStruct {
	if sArg.CustomTypeOffset < 0 {
		return nil
	}

	customTypePkg := prgrm.Packages[s.Structs[sArg.CustomTypeOffset].PackageOffset]
	sStrct := s.Structs[sArg.CustomTypeOffset]
	customTypeName := dsName(sStrct.NameOffset, sStrct.NameSize, s)

	for _, strct := range customTypePkg.Structs {
		if strct.Name == customTypeName {
			return strct
		}
	}

	return nil
}

func dsArgument(sArg *sArgument, s *sAll, prgrm *CXProgram) *CXArgument {
	var arg CXArgument
	arg.Name = dsName(sArg.NameOffset, sArg.NameSize, s)
	arg.Type = int(sArg.Type)

	arg.CustomType = getCustomType(sArg, s, prgrm)

	arg.Size = int(sArg.Size)
	arg.TotalSize = int(sArg.TotalSize)
	arg.Offset = int(sArg.Offset)
	arg.IndirectionLevels = int(sArg.IndirectionLevels)
	arg.DereferenceLevels = int(sArg.DereferenceLevels)
	arg.PassBy = int(sArg.PassBy)

	arg.DereferenceOperations = dsIntegers(sArg.DereferenceOperationsOffset, sArg.DereferenceOperationsSize, s)
	arg.DeclarationSpecifiers = dsIntegers(sArg.DeclarationSpecifiersOffset, sArg.DeclarationSpecifiersSize, s)

	arg.IsSlice = dsBool(sArg.IsSlice)
	arg.IsArray = dsBool(sArg.IsArray)
	arg.IsArrayFirst = dsBool(sArg.IsArrayFirst)
	arg.IsPointer = dsBool(sArg.IsPointer)
	arg.IsReference = dsBool(sArg.IsReference)
	arg.IsDereferenceFirst = dsBool(sArg.IsDereferenceFirst)
	arg.IsStruct = dsBool(sArg.IsStruct)
	arg.IsRest = dsBool(sArg.IsRest)
	arg.IsLocalDeclaration = dsBool(sArg.IsLocalDeclaration)
	arg.IsShortDeclaration = dsBool(sArg.IsShortDeclaration)
	arg.PreviouslyDeclared = dsBool(sArg.PreviouslyDeclared)
	arg.DoesEscape = dsBool(sArg.DoesEscape)

	arg.Lengths = dsIntegers(sArg.LengthsOffset, sArg.LengthsSize, s)
	arg.Indexes = dsArguments(sArg.IndexesOffset, sArg.IndexesSize, s, prgrm)
	arg.Fields = dsArguments(sArg.FieldsOffset, sArg.FieldsSize, s, prgrm)
	arg.Inputs = dsArguments(sArg.InputsOffset, sArg.InputsSize, s, prgrm)
	arg.Outputs = dsArguments(sArg.OutputsOffset, sArg.OutputsSize, s, prgrm)

	arg.Package = prgrm.Packages[sArg.PackageOffset]

	return &arg
}

func getOperator(sExpr *sExpression, s *sAll, prgrm *CXProgram) *CXFunction {
	if sExpr.OperatorOffset < 0 {
		return nil
	}

	opPkg := prgrm.Packages[s.Functions[sExpr.OperatorOffset].PackageOffset]
	sOp := s.Functions[sExpr.OperatorOffset]
	opName := dsName(sOp.NameOffset, sOp.NameSize, s)

	for _, fn := range opPkg.Functions {
		if fn.Name == opName {
			return fn
		}
	}

	return nil
}

func getImport(sImp *sPackage, s *sAll, prgrm *CXProgram) *CXPackage {
	impName := dsName(sImp.NameOffset, sImp.NameSize, s)

	for _, pkg := range prgrm.Packages {
		if pkg.Name == impName {
			return pkg
		}
	}

	return nil
}

func getFunction(sExpr *sExpression, s *sAll, prgrm *CXProgram) *CXFunction {
	if sExpr.FunctionOffset < 0 {
		return nil
	}

	fnPkg := prgrm.Packages[s.Functions[sExpr.FunctionOffset].PackageOffset]
	sFn := s.Functions[sExpr.FunctionOffset]
	fnName := dsName(sFn.NameOffset, sFn.NameSize, s)

	for _, fn := range fnPkg.Functions {
		if fn.Name == fnName {
			return fn
		}
	}

	return nil
}

func dsExpressions(off int32, size int32, s *sAll, prgrm *CXProgram) []*CXExpression {
	if size < 1 {
		return nil
	}

	// getting indexes of expressions
	idxs := dsIntegers(off, size, s)

	// sExprs := s.Expressions[off : off + size]
	exprs := make([]*CXExpression, size)
	for i, idx := range idxs {
		exprs[i] = dsExpression(&s.Expressions[idx], s, prgrm)
	}
	return exprs
}

func dsExpression(sExpr *sExpression, s *sAll, prgrm *CXProgram) *CXExpression {
	var expr CXExpression

	if dsBool(sExpr.IsNative) {
		expr.Operator = Natives[int(sExpr.OpCode)]
	} else {
		expr.Operator = getOperator(sExpr, s, prgrm)
	}

	expr.Inputs = dsArguments(sExpr.InputsOffset, sExpr.InputsSize, s, prgrm)
	expr.Outputs = dsArguments(sExpr.OutputsOffset, sExpr.OutputsSize, s, prgrm)

	expr.Label = dsName(sExpr.LabelOffset, sExpr.LabelSize, s)

	expr.ThenLines = int(sExpr.ThenLines)
	expr.ElseLines = int(sExpr.ElseLines)
	expr.ScopeOperation = int(sExpr.ScopeOperation)

	expr.IsMethodCall = dsBool(sExpr.IsMethodCall)
	expr.IsStructLiteral = dsBool(sExpr.IsStructLiteral)
	expr.IsArrayLiteral = dsBool(sExpr.IsArrayLiteral)
	expr.IsUndType = dsBool(sExpr.IsUndType)
	expr.IsBreak = dsBool(sExpr.IsBreak)
	expr.IsContinue = dsBool(sExpr.IsContinue)

	expr.Function = getFunction(sExpr, s, prgrm)
	expr.Package = prgrm.Packages[sExpr.PackageOffset]

	return &expr
}

func dsFunction(sFn *sFunction, fn *CXFunction, s *sAll, prgrm *CXProgram) {
	fn.Name = dsName(sFn.NameOffset, sFn.NameSize, s)
	fn.Inputs = dsArguments(sFn.InputsOffset, sFn.InputsSize, s, prgrm)
	fn.Outputs = dsArguments(sFn.OutputsOffset, sFn.OutputsSize, s, prgrm)
	fn.ListOfPointers = dsArguments(sFn.ListOfPointersOffset, sFn.ListOfPointersSize, s, prgrm)
	fn.Expressions = dsExpressions(sFn.ExpressionsOffset, sFn.ExpressionsSize, s, prgrm)
	fn.Size = int(sFn.Size)
	fn.Length = int(sFn.Length)

	if sFn.CurrentExpressionOffset > 0 {
		fn.CurrentExpression = fn.Expressions[sFn.CurrentExpressionOffset]
	}

	fn.Package = prgrm.Packages[sFn.PackageOffset]
}

func dsBool(val int32) bool {
	return val == 1
}

func dsIntegers(off int32, size int32, s *sAll) []int {
	if size < 1 {
		return nil
	}
	ints := s.Integers[off : off+size]
	res := make([]int, len(ints))
	for i, in := range ints {
		res[i] = int(in)
	}

	return res
}

// initDeserialization initializes the CXProgram fields that represent a CX program. This should be refactored, as the names Deserialize and initDeserialization create some naming conflict.
func initDeserialization(prgrm *CXProgram, s *sAll) {
	prgrm.Memory = s.Memory
	prgrm.Packages = make([]*CXPackage, len(s.Packages))
	prgrm.CallStack = make([]CXCall, CALLSTACK_SIZE)
	prgrm.HeapStartsAt = int(s.Program.HeapStartsAt)
	prgrm.HeapPointer = int(s.Program.HeapPointer)
	prgrm.StackSize = int(s.Program.StackSize)
	prgrm.HeapSize = int(s.Program.HeapSize)
	prgrm.Version = dsName(s.Program.VersionOffset, s.Program.VersionSize, s)

	dsPackages(s, prgrm)
}

// Deserialize deserializes a serialized CX program back to its golang struct representation.
func Deserialize(byts []byte) (prgrm *CXProgram) {
	prgrm = &CXProgram{}
	idxSize := encoder.Size(sIndex{})

	var s sAll

	DeserializeRaw(byts[:idxSize], &s.Index)
	DeserializeRaw(byts[s.Index.ProgramOffset:s.Index.CallsOffset], &s.Program)
	DeserializeRaw(byts[s.Index.CallsOffset:s.Index.PackagesOffset], &s.Calls)
	DeserializeRaw(byts[s.Index.PackagesOffset:s.Index.StructsOffset], &s.Packages)
	DeserializeRaw(byts[s.Index.StructsOffset:s.Index.FunctionsOffset], &s.Structs)
	DeserializeRaw(byts[s.Index.FunctionsOffset:s.Index.ExpressionsOffset], &s.Functions)
	DeserializeRaw(byts[s.Index.ExpressionsOffset:s.Index.ArgumentsOffset], &s.Expressions)
	DeserializeRaw(byts[s.Index.ArgumentsOffset:s.Index.IntegersOffset], &s.Arguments)
	DeserializeRaw(byts[s.Index.IntegersOffset:s.Index.NamesOffset], &s.Integers)
	s.Names = byts[s.Index.NamesOffset:s.Index.MemoryOffset]
	s.Memory = byts[s.Index.MemoryOffset:]

	initDeserialization(prgrm, &s)

	// prgrm.PrintProgram()

	return prgrm
}

// CopyProgramState copies the program state from `prgrm1` to `prgrm2`.
func CopyProgramState(sPrgrm1, sPrgrm2 *[]byte) {
	idxSize := encoder.Size(sIndex{})

	var index1 sIndex
	var index2 sIndex

	DeserializeRaw((*sPrgrm1)[:idxSize], &index1)
	DeserializeRaw((*sPrgrm2)[:idxSize], &index2)

	var prgrm1Info sProgram
	DeserializeRaw((*sPrgrm1)[index1.ProgramOffset:index1.CallsOffset], &prgrm1Info)

	var prgrm2Info sProgram
	DeserializeRaw((*sPrgrm2)[index2.ProgramOffset:index2.CallsOffset], &prgrm2Info)

	// the stack segment should be 0 for prgrm1, but just in case
	var prgrmState []byte
	prgrmState = append(prgrmState, make([]byte, prgrm2Info.StackSize)...)
	// We are only interested on extracting the data segment
	prgrmState = append(prgrmState, (*sPrgrm1)[index1.NamesOffset+prgrm1Info.StackSize:index1.NamesOffset+prgrm1Info.StackSize+(prgrm2Info.HeapStartsAt-prgrm2Info.StackSize)]...)

	for i, byt := range prgrmState {
		(*sPrgrm2)[i+int(index2.MemoryOffset)] = byt
	}
}

// GetSerializedMemoryOffset returns the offset at which the memory of a serialized CX program starts.
// Only used in one place? Delete if possible
// ./cxgo/cxlexer/state.go:30:	memOffset := cxcore.GetSerializedMemoryOffset(progS)
func GetSerializedMemoryOffset(sPrgrm []byte) int {
	idxSize := encoder.Size(sIndex{})
	var index sIndex
	DeserializeRaw(sPrgrm[:idxSize], &index)
	return int(index.MemoryOffset)
}

// GetSerializedStackSize returns the stack size of a serialized CX program starts.
func GetSerializedStackSize(sPrgrm []byte) int {
	idxSize := encoder.Size(sIndex{})
	var index sIndex
	DeserializeRaw(sPrgrm[:idxSize], &index)

	var prgrmInfo sProgram
	DeserializeRaw(sPrgrm[index.ProgramOffset:index.CallsOffset], &prgrmInfo)

	return int(prgrmInfo.StackSize)
}

// GetSerializedDataSize returns the size of the data segment of a serialized CX program.
func GetSerializedDataSize(sPrgrm []byte) int {
	idxSize := encoder.Size(sIndex{})
	var index sIndex
	DeserializeRaw(sPrgrm[:idxSize], &index)

	var prgrmInfo sProgram
	DeserializeRaw(sPrgrm[index.ProgramOffset:index.CallsOffset], &prgrmInfo)

	return int(prgrmInfo.HeapStartsAt - prgrmInfo.StackSize)
}

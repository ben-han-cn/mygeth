// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

type (
	executionFunc       func(pc *uint64, env *EVM, contract *Contract, memory *Memory, stack *Stack) ([]byte, error)
	gasFunc             func(params.GasTable, *EVM, *Contract, *Stack, *Memory, uint64) (uint64, error) // last parameter is the requested memory size as a uint64
	stackValidationFunc func(*Stack) error
	memorySizeFunc      func(*Stack) *big.Int
)

var errGasUintOverflow = errors.New("gas uint64 overflow")

type operation struct {
	// op is the operation function
	execute executionFunc
	// validateStack validates the stack (size) for the operation
	validateStack stackValidationFunc
	// memorySize returns the memory size required for the operation
	memorySize memorySizeFunc
	// halts indicates whether the operation shoult halt further execution
	// and return
	halts bool
	// jumps indicates whether operation made a jump. This prevents the program
	// counter from further incrementing.
	jumps bool
	// writes determines whether this a state modifying operation
	writes bool
	// valid is used to check whether the retrieved operation is valid and known
	valid bool
	// reverts determined whether the operation reverts state
	reverts bool
}

var (
	frontierInstructionSet  = NewFrontierInstructionSet()
	homesteadInstructionSet = NewHomesteadInstructionSet()
)

// NewHomesteadInstructionSet returns the frontier and homestead
// instructions that can be executed during the homestead phase.
func NewHomesteadInstructionSet() [256]operation {
	instructionSet := NewFrontierInstructionSet()
	instructionSet[DELEGATECALL] = operation{
		execute:       opDelegateCall,
		validateStack: makeStackFunc(6, 1),
		memorySize:    memoryDelegateCall,
		valid:         true,
	}
	return instructionSet
}

// NewFrontierInstructionSet returns the frontier instructions
// that can be executed during the frontier phase.
func NewFrontierInstructionSet() [256]operation {
	return [256]operation{
		STOP: {
			execute:       opStop,
			validateStack: makeStackFunc(0, 0),
			halts:         true,
			valid:         true,
		},
		ADD: {
			execute:       opAdd,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		MUL: {
			execute:       opMul,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		SUB: {
			execute:       opSub,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		DIV: {
			execute:       opDiv,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		SDIV: {
			execute:       opSdiv,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		MOD: {
			execute:       opMod,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		SMOD: {
			execute:       opSmod,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		ADDMOD: {
			execute:       opAddmod,
			validateStack: makeStackFunc(3, 1),
			valid:         true,
		},
		MULMOD: {
			execute:       opMulmod,
			validateStack: makeStackFunc(3, 1),
			valid:         true,
		},
		EXP: {
			execute:       opExp,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		SIGNEXTEND: {
			execute:       opSignExtend,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		LT: {
			execute:       opLt,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		GT: {
			execute:       opGt,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		SLT: {
			execute:       opSlt,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		SGT: {
			execute:       opSgt,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		EQ: {
			execute:       opEq,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		ISZERO: {
			execute:       opIszero,
			validateStack: makeStackFunc(1, 1),
			valid:         true,
		},
		AND: {
			execute:       opAnd,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		XOR: {
			execute:       opXor,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		OR: {
			execute:       opOr,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		NOT: {
			execute:       opNot,
			validateStack: makeStackFunc(1, 1),
			valid:         true,
		},
		BYTE: {
			execute:       opByte,
			validateStack: makeStackFunc(2, 1),
			valid:         true,
		},
		SHA3: {
			execute:       opSha3,
			validateStack: makeStackFunc(2, 1),
			memorySize:    memorySha3,
			valid:         true,
		},
		ADDRESS: {
			execute:       opAddress,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		BALANCE: {
			execute:       opBalance,
			validateStack: makeStackFunc(1, 1),
			valid:         true,
		},
		ORIGIN: {
			execute:       opOrigin,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		CALLER: {
			execute:       opCaller,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		CALLVALUE: {
			execute:       opCallValue,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		CALLDATALOAD: {
			execute:       opCalldataLoad,
			validateStack: makeStackFunc(1, 1),
			valid:         true,
		},
		CALLDATASIZE: {
			execute:       opCalldataSize,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		CALLDATACOPY: {
			execute:       opCalldataCopy,
			validateStack: makeStackFunc(3, 0),
			memorySize:    memoryCalldataCopy,
			valid:         true,
		},
		CODESIZE: {
			execute:       opCodeSize,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		CODECOPY: {
			execute:       opCodeCopy,
			validateStack: makeStackFunc(3, 0),
			memorySize:    memoryCodeCopy,
			valid:         true,
		},
		EXTCODESIZE: {
			execute:       opExtCodeSize,
			validateStack: makeStackFunc(1, 1),
			valid:         true,
		},
		EXTCODECOPY: {
			execute:       opExtCodeCopy,
			validateStack: makeStackFunc(4, 0),
			memorySize:    memoryExtCodeCopy,
			valid:         true,
		},
		BLOCKHASH: {
			execute:       opBlockhash,
			validateStack: makeStackFunc(1, 1),
			valid:         true,
		},
		COINBASE: {
			execute:       opCoinbase,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		TIMESTAMP: {
			execute:       opTimestamp,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		NUMBER: {
			execute:       opNumber,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		DIFFICULTY: {
			execute:       opDifficulty,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		POP: {
			execute:       opPop,
			validateStack: makeStackFunc(1, 0),
			valid:         true,
		},
		MLOAD: {
			execute:       opMload,
			validateStack: makeStackFunc(1, 1),
			memorySize:    memoryMLoad,
			valid:         true,
		},
		MSTORE: {
			execute:       opMstore,
			validateStack: makeStackFunc(2, 0),
			memorySize:    memoryMStore,
			valid:         true,
		},
		MSTORE8: {
			execute:       opMstore8,
			memorySize:    memoryMStore8,
			validateStack: makeStackFunc(2, 0),

			valid: true,
		},
		SLOAD: {
			execute:       opSload,
			validateStack: makeStackFunc(1, 1),
			valid:         true,
		},
		SSTORE: {
			execute:       opSstore,
			validateStack: makeStackFunc(2, 0),
			valid:         true,
			writes:        true,
		},
		JUMP: {
			execute:       opJump,
			validateStack: makeStackFunc(1, 0),
			jumps:         true,
			valid:         true,
		},
		JUMPI: {
			execute:       opJumpi,
			validateStack: makeStackFunc(2, 0),
			jumps:         true,
			valid:         true,
		},
		PC: {
			execute:       opPc,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		MSIZE: {
			execute:       opMsize,
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		JUMPDEST: {
			execute:       opJumpdest,
			validateStack: makeStackFunc(0, 0),
			valid:         true,
		},
		PUSH1: {
			execute:       makePush(1, 1),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH2: {
			execute:       makePush(2, 2),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH3: {
			execute:       makePush(3, 3),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH4: {
			execute:       makePush(4, 4),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH5: {
			execute:       makePush(5, 5),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH6: {
			execute:       makePush(6, 6),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH7: {
			execute:       makePush(7, 7),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH8: {
			execute:       makePush(8, 8),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH9: {
			execute:       makePush(9, 9),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH10: {
			execute:       makePush(10, 10),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH11: {
			execute:       makePush(11, 11),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH12: {
			execute:       makePush(12, 12),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH13: {
			execute:       makePush(13, 13),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH14: {
			execute:       makePush(14, 14),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH15: {
			execute:       makePush(15, 15),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH16: {
			execute:       makePush(16, 16),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH17: {
			execute:       makePush(17, 17),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH18: {
			execute:       makePush(18, 18),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH19: {
			execute:       makePush(19, 19),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH20: {
			execute:       makePush(20, 20),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH21: {
			execute:       makePush(21, 21),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH22: {
			execute:       makePush(22, 22),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH23: {
			execute:       makePush(23, 23),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH24: {
			execute:       makePush(24, 24),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH25: {
			execute:       makePush(25, 25),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH26: {
			execute:       makePush(26, 26),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH27: {
			execute:       makePush(27, 27),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH28: {
			execute:       makePush(28, 28),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH29: {
			execute:       makePush(29, 29),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH30: {
			execute:       makePush(30, 30),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH31: {
			execute:       makePush(31, 31),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		PUSH32: {
			execute:       makePush(32, 32),
			validateStack: makeStackFunc(0, 1),
			valid:         true,
		},
		DUP1: {
			execute:       makeDup(1),
			validateStack: makeDupStackFunc(1),
			valid:         true,
		},
		DUP2: {
			execute:       makeDup(2),
			validateStack: makeDupStackFunc(2),
			valid:         true,
		},
		DUP3: {
			execute:       makeDup(3),
			validateStack: makeDupStackFunc(3),
			valid:         true,
		},
		DUP4: {
			execute:       makeDup(4),
			validateStack: makeDupStackFunc(4),
			valid:         true,
		},
		DUP5: {
			execute:       makeDup(5),
			validateStack: makeDupStackFunc(5),
			valid:         true,
		},
		DUP6: {
			execute:       makeDup(6),
			validateStack: makeDupStackFunc(6),
			valid:         true,
		},
		DUP7: {
			execute:       makeDup(7),
			validateStack: makeDupStackFunc(7),
			valid:         true,
		},
		DUP8: {
			execute:       makeDup(8),
			validateStack: makeDupStackFunc(8),
			valid:         true,
		},
		DUP9: {
			execute:       makeDup(9),
			validateStack: makeDupStackFunc(9),
			valid:         true,
		},
		DUP10: {
			execute:       makeDup(10),
			validateStack: makeDupStackFunc(10),
			valid:         true,
		},
		DUP11: {
			execute:       makeDup(11),
			validateStack: makeDupStackFunc(11),
			valid:         true,
		},
		DUP12: {
			execute:       makeDup(12),
			validateStack: makeDupStackFunc(12),
			valid:         true,
		},
		DUP13: {
			execute:       makeDup(13),
			validateStack: makeDupStackFunc(13),
			valid:         true,
		},
		DUP14: {
			execute:       makeDup(14),
			validateStack: makeDupStackFunc(14),
			valid:         true,
		},
		DUP15: {
			execute:       makeDup(15),
			validateStack: makeDupStackFunc(15),
			valid:         true,
		},
		DUP16: {
			execute:       makeDup(16),
			validateStack: makeDupStackFunc(16),
			valid:         true,
		},
		SWAP1: {
			execute:       makeSwap(1),
			validateStack: makeSwapStackFunc(2),
			valid:         true,
		},
		SWAP2: {
			execute:       makeSwap(2),
			validateStack: makeSwapStackFunc(3),
			valid:         true,
		},
		SWAP3: {
			execute:       makeSwap(3),
			validateStack: makeSwapStackFunc(4),
			valid:         true,
		},
		SWAP4: {
			execute:       makeSwap(4),
			validateStack: makeSwapStackFunc(5),
			valid:         true,
		},
		SWAP5: {
			execute:       makeSwap(5),
			validateStack: makeSwapStackFunc(6),
			valid:         true,
		},
		SWAP6: {
			execute:       makeSwap(6),
			validateStack: makeSwapStackFunc(7),
			valid:         true,
		},
		SWAP7: {
			execute:       makeSwap(7),
			validateStack: makeSwapStackFunc(8),
			valid:         true,
		},
		SWAP8: {
			execute:       makeSwap(8),
			validateStack: makeSwapStackFunc(9),
			valid:         true,
		},
		SWAP9: {
			execute:       makeSwap(9),
			validateStack: makeSwapStackFunc(10),
			valid:         true,
		},
		SWAP10: {
			execute:       makeSwap(10),
			validateStack: makeSwapStackFunc(11),
			valid:         true,
		},
		SWAP11: {
			execute:       makeSwap(11),
			validateStack: makeSwapStackFunc(12),
			valid:         true,
		},
		SWAP12: {
			execute:       makeSwap(12),
			validateStack: makeSwapStackFunc(13),
			valid:         true,
		},
		SWAP13: {
			execute:       makeSwap(13),
			validateStack: makeSwapStackFunc(14),
			valid:         true,
		},
		SWAP14: {
			execute:       makeSwap(14),
			validateStack: makeSwapStackFunc(15),
			valid:         true,
		},
		SWAP15: {
			execute:       makeSwap(15),
			validateStack: makeSwapStackFunc(16),
			valid:         true,
		},
		SWAP16: {
			execute:       makeSwap(16),
			validateStack: makeSwapStackFunc(17),
			valid:         true,
		},
		LOG0: {
			execute:       makeLog(0),
			validateStack: makeStackFunc(2, 0),
			memorySize:    memoryLog,
			valid:         true,
		},
		LOG1: {
			execute:       makeLog(1),
			validateStack: makeStackFunc(3, 0),
			memorySize:    memoryLog,
			valid:         true,
		},
		LOG2: {
			execute:       makeLog(2),
			validateStack: makeStackFunc(4, 0),
			memorySize:    memoryLog,
			valid:         true,
		},
		LOG3: {
			execute:       makeLog(3),
			validateStack: makeStackFunc(5, 0),
			memorySize:    memoryLog,
			valid:         true,
		},
		LOG4: {
			execute:       makeLog(4),
			validateStack: makeStackFunc(6, 0),
			memorySize:    memoryLog,
			valid:         true,
		},
		CREATE: {
			execute:       opCreate,
			validateStack: makeStackFunc(3, 1),
			memorySize:    memoryCreate,
			valid:         true,
			writes:        true,
		},
		CALL: {
			execute:       opCall,
			validateStack: makeStackFunc(7, 1),
			memorySize:    memoryCall,
			valid:         true,
		},
		CALLCODE: {
			execute:       opCallCode,
			validateStack: makeStackFunc(7, 1),
			memorySize:    memoryCall,
			valid:         true,
		},
		RETURN: {
			execute:       opReturn,
			validateStack: makeStackFunc(2, 0),
			memorySize:    memoryReturn,
			halts:         true,
			valid:         true,
		},
		SELFDESTRUCT: {
			execute:       opSuicide,
			validateStack: makeStackFunc(1, 0),
			halts:         true,
			valid:         true,
			writes:        true,
		},
	}
}

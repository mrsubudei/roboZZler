package service

import (
	"strconv"
)

const (
	North = "north"
	East  = "east"
	West  = "west"
	South = "south"
)

const (
	StepForward    = "F"
	ForwardOnBlue  = "FB"
	ForwardOnRed   = "FR"
	ForwardOnGreen = "FG"
	TurnLeft       = "L"
	LeftOnBlue     = "LB"
	LeftOnRed      = "LR"
	LeftOnGreen    = "LG"
	TurnRight      = "R"
	RightOnBlue    = "RB"
	RightOnRed     = "RR"
	RightOnGreen   = "RG"
	PrintBlue      = "PB"
	PrintRed       = "PR"
	PrintGreen     = "PG"
)

var (
	defalutCmds = []string{StepForward, ForwardOnBlue, ForwardOnRed, ForwardOnGreen,
		TurnLeft, LeftOnBlue, LeftOnRed, LeftOnGreen, TurnRight, RightOnBlue,
		RightOnRed, RightOnGreen, PrintBlue, PrintRed, PrintGreen}
)

func Solve(board [][]string, toVisit map[[2]int]struct{}, cmdLines []int,
	startPos [2]int, startDir string) ([][]string, bool) {

	funcCmds := getFuncCmds(cmdLines)
	allCmds := append(funcCmds, defalutCmds...)
	cmdRanges := getFuncCmdsRanges(cmdLines)
	cmdCount := getCmdsCount(cmdLines)
	prefix := make([]string, cmdCount)
	foundCmds := []string{}
	isFinished := false

	bruteForce(allCmds, prefix, cmdCount-1, board, cmdRanges,
		&foundCmds, toVisit, &isFinished, startPos, startDir)

	if len(foundCmds) == 0 {
		return nil, false
	}

	ans := make([][]string, len(cmdLines))
	for i := 0; i < len(ans); i++ {
		ansLine := []string{}
		cmdRange := cmdRanges[i]
		for j := cmdRange[0]; j < cmdRange[1]; j++ {
			ansLine = append(ansLine, foundCmds[j])
		}
		ans[i] = ansLine
	}

	return ans, true
}

func bruteForce(cmds []string, prefix []string, depth int,
	board [][]string, cmdRanges [][2]int, ans *[]string,
	toVisit map[[2]int]struct{}, isFinished *bool, startPos [2]int,
	startDir string) {

	if depth == -1 {
		if !hasUnreachableCmds(prefix, cmdRanges) {
			copiedBoard := copyBoard(board)
			copiedToVisitCells := copyToVisitCells(toVisit)
			if tryToPass(copiedBoard, prefix, cmdRanges,
				copiedToVisitCells, startPos, startDir) {
				*isFinished = true
				*ans = prefix
			}
		}
		return
	}

	if *isFinished {
		return
	}

	for _, s := range cmds {
		if *isFinished {
			return
		}
		prefix[depth] = s
		bruteForce(cmds, prefix, depth-1, board,
			cmdRanges, ans, toVisit, isFinished, startPos, startDir)
	}
}

func tryToPass(board [][]string, cmds []string, cmdRanges [][2]int,
	toVisit map[[2]int]struct{}, startPos [2]int, startDir string) bool {

	toVisitCount := len(toVisit)
	cmdPos := 0
	curPos := startPos
	moves := 0

	for {
		if cmdPos >= len(cmds) {
			break
		}
		curCmd := cmds[cmdPos]
		newCmdLine, isCmdChanged := makeCmd(board, &curCmd, &curPos, &startDir)
		if isCmdChanged {
			cmdPos = cmdRanges[newCmdLine][0]
		} else {
			cmdPos++
		}

		if isOutOfZone(board, curPos) {
			break
		}
		if isNeededCellVisited(curPos, toVisit) {
			toVisitCount--
			if toVisitCount == 0 {
				break
			}
		}

		moves++
		// check if robot is wondering in circle
		if moves == 10000 {
			break
		}
	}

	return toVisitCount == 0
}

func hasUnreachableCmds(cmds []string, cmdRanges [][2]int) bool {
	for i := 0; i < len(cmdRanges); i++ {
		cmdRange := cmdRanges[i]
		self := strconv.Itoa(i)
		if cmdRange[1]-cmdRange[0] == 1 {
			if cmds[cmdRange[0]] == self {
				return true
			}
		}
		for j := cmdRange[0]; j < cmdRange[1]-1; j++ {
			if cmds[j] == self {
				return true
			}
		}
	}

	return false
}

func isOutOfZone(board [][]string, curPos [2]int) bool {
	if curPos[0] == 0 || curPos[0] == len(board)-1 ||
		curPos[1] == 0 || curPos[1] == len(board[0]) {
		return true
	}

	return false
}

func isNeededCellVisited(curPos [2]int, toVisit map[[2]int]struct{}) bool {
	_, ok := toVisit[curPos]
	if ok {
		delete(toVisit, curPos)
	}

	return ok
}

func makeCmd(board [][]string, cmd *string, curPos *[2]int, curDir *string) (int, bool) {
	isCmdChanged := false
	newCmdLine := 0

	switch *cmd {
	case StepForward:
		moveForward(curPos, curDir)
	case ForwardOnBlue:
		if board[curPos[0]][curPos[1]] == "B" {
			moveForward(curPos, curDir)
		}
	case ForwardOnRed:
		if board[curPos[0]][curPos[1]] == "R" {
			moveForward(curPos, curDir)
		}
	case ForwardOnGreen:
		if board[curPos[0]][curPos[1]] == "G" {
			moveForward(curPos, curDir)
		}
	case TurnLeft:
		turnLeft(curDir)
	case LeftOnBlue:
		if board[curPos[0]][curPos[1]] == "B" {
			turnLeft(curDir)
		}
	case LeftOnRed:
		if board[curPos[0]][curPos[1]] == "R" {
			turnLeft(curDir)
		}
	case LeftOnGreen:
		if board[curPos[0]][curPos[1]] == "G" {
			turnLeft(curDir)
		}
	case TurnRight:
		turnRight(curDir)
	case RightOnBlue:
		if board[curPos[0]][curPos[1]] == "B" {
			turnRight(curDir)
		}
	case RightOnRed:
		if board[curPos[0]][curPos[1]] == "R" {
			turnRight(curDir)
		}
	case RightOnGreen:
		if board[curPos[0]][curPos[1]] == "G" {
			turnRight(curDir)
		}
	case PrintBlue:
		board[curPos[0]][curPos[1]] = "B"
	case PrintRed:
		board[curPos[0]][curPos[1]] = "R"
	case PrintGreen:
		board[curPos[0]][curPos[1]] = "G"
	default:
		cmdLine, _ := strconv.Atoi(*cmd)
		newCmdLine = cmdLine
		isCmdChanged = true
	}

	return newCmdLine, isCmdChanged
}

func turnLeft(curDir *string) {
	switch *curDir {
	case North:
		*curDir = West
	case West:
		*curDir = South
	case South:
		*curDir = East
	case East:
		*curDir = North
	}
}

func turnRight(curDir *string) {
	switch *curDir {
	case North:
		*curDir = East
	case East:
		*curDir = South
	case South:
		*curDir = West
	case West:
		*curDir = North
	}
}

func moveForward(curPos *[2]int, curDir *string) {
	switch *curDir {
	case North:
		curPos[0] -= 1
	case East:
		curPos[1] += 1
	case West:
		curPos[1] -= 1
	case South:
		curPos[0] += 1
	}
}

func getCmdsCount(cmdLines []int) int {
	count := 0
	for _, v := range cmdLines {
		count += v
	}

	return count
}

func getFuncCmdsRanges(cmdLines []int) [][2]int {
	cmdsRanges := make([][2]int, len(cmdLines))

	from := 0
	for i, v := range cmdLines {
		cmdsRanges[i] = [2]int{from, from + v}
		from = v
	}

	return cmdsRanges
}

func getFuncCmds(cmdLines []int) []string {
	funcCommands := make([]string, len(cmdLines))

	for i := 0; i < len(cmdLines); i++ {
		funcCommands[i] = strconv.Itoa(i)
	}

	return funcCommands
}

func copyBoard(board [][]string) [][]string {
	copied := make([][]string, len(board))
	for i, v := range board {
		tmpSl := make([]string, len(v))
		copy(tmpSl, v)
		copied[i] = tmpSl
	}

	return copied
}

func copyToVisitCells(toVisit map[[2]int]struct{}) map[[2]int]struct{} {
	copied := make(map[[2]int]struct{}, len(toVisit))

	for key := range toVisit {
		copied[key] = struct{}{}
	}

	return copied
}

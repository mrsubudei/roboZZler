package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mrsubudei/roboZZler/internal/service"
)

var (
	ErrCommandLinesInvalid   = fmt.Errorf("invalid command lines")
	ErrToVisitCellsInvalid   = fmt.Errorf("invalid to visit cells")
	ErrBoardInvalid          = fmt.Errorf("invalid board")
	ErrStartPositionInvalid  = fmt.Errorf("invalid start postion")
	ErrStartDirectionInvalid = fmt.Errorf("invalid start direction")
)

type SolveReq struct {
	Board          [][]string `json:"board"`
	ToVisitCells   [][2]int   `json:"toVisitCells"`
	CommandLines   []int      `json:"commandLines"`
	StartPosition  [2]int     `json:"startPosition"`
	StartDirection string     `json:"startDirection"`
}

type SolveResponse struct {
	Commands [][]string `json:"commands"`
}

func SolveRobozzle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/solve" {
		http.Error(w, "Page is Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not Allowed", http.StatusMethodNotAllowed)
		return
	}

	req := SolveReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("req: %#v\n", req)
	if err := validatePayload(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	solved, ok := service.Solve(
		broadBoard(req.Board),
		getToVisitMap(req.ToVisitCells),
		req.CommandLines,
		[2]int{req.StartPosition[0] + 1, req.StartPosition[1] + 1},
		req.StartDirection,
	)
	if !ok {
		http.Error(w, "Has not been solved", http.StatusOK)
		return
	}

	resp := SolveResponse{
		Commands: solved,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(jsonResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getToVisitMap(toVisitCells [][2]int) map[[2]int]struct{} {
	m := make(map[[2]int]struct{}, len(toVisitCells))
	for _, v := range toVisitCells {
		m[[2]int{v[0] + 1, v[1] + 1}] = struct{}{}
	}

	return m
}

func validatePayload(req SolveReq) error {
	if !validateCommandLines(req.CommandLines) {
		return ErrCommandLinesInvalid
	}
	if validateToVisitCells(req.ToVisitCells) {
		return ErrToVisitCellsInvalid
	}
	if !validateBoard(req.Board) {
		return ErrBoardInvalid
	}
	if !validateStartPosition(req.Board, req.StartPosition) {
		return ErrStartPositionInvalid
	}
	if !validateStartDirection(req.StartDirection) {
		return ErrStartDirectionInvalid
	}

	return nil
}

func validateCommandLines(commandLines []int) bool {
	if len(commandLines) == 0 {
		return false
	}

	for _, v := range commandLines {
		if v == 0 {
			return false
		}
	}

	return true
}

func validateToVisitCells(ToVisitCells [][2]int) bool {
	return len(ToVisitCells) == 0
}

func validateStartDirection(startDirection string) bool {
	switch startDirection {
	case "north":
	case "east":
	case "west":
	case "south":
	default:
		return false
	}

	return true
}
func validateBoard(board [][]string) bool {
	if len(board) == 0 {
		return false
	}

	rowLength := len(board[0])
	for i := range board {
		if len(board[i]) != rowLength {
			return false
		}
		for j := range board[i] {
			if board[i][j] != "G" && board[i][j] != "B" &&
				board[i][j] != "R" && board[i][j] != "" {
				return false
			}
		}
	}

	return true
}

func validateStartPosition(board [][]string, startPosition [2]int) bool {
	if startPosition[0] < 0 || startPosition[0] >= len(board) {
		return false
	}

	if startPosition[1] < 0 || startPosition[1] >= len(board[0]) {
		return false
	}

	return true
}

func broadBoard(board [][]string) [][]string {
	broadedBoard := make([][]string, len(board)+2)
	for i, v := range board {
		row := make([]string, len(v)+2)
		row[0] = "#"
		row[len(row)-1] = "#"
		for j := 0; j < len(v); j++ {
			if v[j] == "G" || v[j] == "B" || v[j] == "R" {
				row[j+1] = v[j]
			} else {
				row[j+1] = "#"
			}
		}
		broadedBoard[i+1] = row
	}

	stubRow := make([]string, len(board[0])+2)
	for i := range stubRow {
		stubRow[i] = "#"
	}
	broadedBoard[0] = stubRow
	broadedBoard[len(broadedBoard)-1] = stubRow

	return broadedBoard
}

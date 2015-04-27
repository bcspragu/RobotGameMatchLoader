package match

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
)

type MatchRecord struct {
	Timestamp uint64
	MatchID   uint64 `json:"match_id"`
}

type MatchData struct {
	Score   []int
	History History
}

type History struct {
	Moves [][]Move
}

type Move struct {
	PlayerID int `json:"player_id"`
	RobotID  int `json:"robot_id"`
	HP       int
	Location []int
	Action   Action
}

type Action struct {
	Action   string
	Location []int
}

type Match struct {
	R1ID int `json:"r1_id"`
	R2ID int `json:"r2_id"`

	R1Name string `json:"r1_name"`
	R2Name string `json:"r2_name"`

	R1Score int `json:"r1_score"`
	R2Score int `json:"r2_score"`

	R1Ranking int `json:"r1_ranking"`
	R2Ranking int `json:"r2_ranking"`

	R1Rating float64 `json:"r1_rating"`
	R2Rating float64 `json:"r2_rating"`

	R1Time float64 `json:"r1_time"`
	R2Time float64 `json:"r2_time"`

	Timestamp uint64
	Ranked    bool
	Seed      uint64
	KFactor   int `json:"k_factor"`
	State     int
	Winner    int

	Data MatchData
}

func (h *History) UnmarshalJSON(b []byte) error {
	data, err := base64.StdEncoding.DecodeString(strings.Replace(string(b[1:len(b)-1]), `\n`, "\n", -1))
	if err == nil {
		var moves [][]Move
		data = data[20 : len(data)-2]
		data[0] = []byte("[")[0]

		err = json.Unmarshal(data, &moves)
		h.Moves = moves
	}
	return err
}

func (a *Action) UnmarshalJSON(b []byte) error {
	var data []interface{}
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	a.Action = data[0].(string)
	if len(data) == 2 && data[1] != nil {
		d := data[1].([]interface{})
		a.Location = make([]int, len(d))
		for i, pos := range data[1].([]interface{}) {
			a.Location[i] = int(pos.(float64))
		}
	}
	return nil
}

func (m *Match) GridString() string {
	var loc [100][17][17]int

	for roundNumber, round := range m.Data.History.Moves {
		for _, move := range round {
			pos := move.Location
			loc[roundNumber][pos[0]-1][pos[1]-1] = ((move.PlayerID * 2) - 1)
		}
	}

	var buffer bytes.Buffer

	for r := 0; r < 100; r++ {
		for i := 0; i < 17; i++ {
			for j := 0; j < 17; j++ {
				buffer.WriteString(strconv.Itoa(loc[r][i][j]))
				buffer.WriteString(" ")
			}
		}
	}
	buffer.WriteString("\n")

	return buffer.String()
}

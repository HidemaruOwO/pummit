package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/HidemaruOwO/nuts/log"
)

func Init(path string, isDebug bool) {
	path = filepath.Join(path)

	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		log.Debugf(isDebug, "File is exists")
	} else {
		// 初期化処理
	}
}

type AliasList struct {
	Alias [][]string `json:"alias"`
}

func GetAliasList() [][]string {
	s := `{
  "alias": [["s","sparkles"],["t","tada"]]
}`

	var aliasList AliasList
	err := json.Unmarshal([]byte(s), &aliasList)
	if err != nil {
		log.ErrorExit(fmt.Errorf("JSON encoding failed"))
	}
	return aliasList.Alias
}

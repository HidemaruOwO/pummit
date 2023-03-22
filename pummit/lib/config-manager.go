package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/HidemaruOwO/nuts/log"
	"github.com/HidemaruOwO/pummit/pummit/config"
	"github.com/tucnak/store"
)

type Gitmoji struct {
	Schema   string `json:"$schema"`
	Gitmojis []struct {
		Emoji       string `json:"emoji"`
		Entity      string `json:"entity"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Name        string `json:"name"`
		Semver      string `json:"semver"`
	} `json:"gitmojis"`
}

type Alias struct {
	WriteEmojiPrefix bool       `json:"writeEmojiPrefix"`
	UseAlias         bool       `json:"useAlias"`
	Alias            [][]string `json:"alias"`
}

func Init(path string, isDebug bool) {
	path = filepath.Join(path)

	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		log.Debugf(isDebug, "File is exists\n")
	} else {
		store.Init("pummit")
		url := "https://raw.githubusercontent.com/carloscuesta/gitmoji/master/packages/gitmojis/src/gitmojis.json"
		res, err := http.Get(url)
		if err != nil {
			log.Criticalf("http get failed\n")
			log.ErrorExit(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Criticalf("read geted data failed\n")
			log.ErrorExit(err)
		}

		var gitmoji Gitmoji
		if err := json.Unmarshal([]byte(body), &gitmoji); err != nil {
			log.Criticalf("JSON encoding failed\n")
			log.ErrorExit(err)
		}
		log.Debugf(isDebug, "Create gitimoji.json\n")
		store.Save("gitimoji.json", gitmoji)

		var alias Alias
		if err := json.Unmarshal([]byte(config.BaseJsonData), &alias); err != nil {
			log.Criticalf("JSON encoding failed\n")
			log.ErrorExit(err)
		}
		log.Debugf(isDebug, "Create config.json\n")
		store.Save("config.json", alias)
	}
}

func GetGitmoji() Gitmoji {
	var gitimoji Gitmoji
	store.SetApplicationName("pummit")
	if err := store.Load("gitimoji.json", &gitimoji); err != nil {
		log.Criticalf("Loading gitimoji list failed\n")
		log.ErrorExit(err)
	}
	return gitimoji
}

func GetAlias() Alias {
	var alias Alias
	store.SetApplicationName("pummit")
	if err := store.Load("config.json", &alias); err != nil {
		log.Criticalf("Loading alias list failed\n")
		log.ErrorExit(err)
	}
	return alias
}

func GetAliasList() [][]string {
	var alias Alias
	store.SetApplicationName("pummit")
	if err := store.Load("config.json", &alias); err != nil {
		log.Criticalf("Loading alias list failed\n")
		log.ErrorExit(err)
	}
	return alias.Alias
}

func WriteConfig(configData Alias) {
	store.SetApplicationName("pummit")
	if err := store.Save("config.json", configData); err != nil {
		log.ErrorExit(err)
	}
}

func PlatformPath(path string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s\\%s", os.Getenv("APPDATA"), path)
	}

	var unixConfigDir string
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		unixConfigDir = xdg
	} else {
		unixConfigDir = os.Getenv("HOME") + "/.config"
	}

	return fmt.Sprintf("%s/%s", unixConfigDir,
		path)
}

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
	"github.com/HidemaruOwO/pummit/src/config"
	"github.com/tucnak/store"
)

type Gitmoji struct {
	Schema   string     `json:"$schema"`
	Gitmojis []Gitmojis `json:"gitmojis"`
}

type Gitmojis struct {
	Emoji       string `json:"emoji"`
	Entity      string `json:"entity"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Semver      string `json:"semver"`
}
type AppConfig struct {
	WriteEmoji           *bool       `json:"writeEmoji"`
	UseAlias             *bool       `json:"useAlias"`
	UseLimitPathesLength *bool       `json:"useLimitPathesLength"`
	LimitPathesLength    *int        `json:"limitPathesLength"`
	Alias                *[][]string `json:"alias"`
}

func Init(path string) {
	path = filepath.Join(path)

	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		log.Debugf(config.IS_DEBUG, "File is exists\n")
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
		log.Debugf(config.IS_DEBUG, "Create gitimoji.json\n")
		store.Save("gitimoji.json", gitmoji)

		var appConfig AppConfig
		if err := json.Unmarshal([]byte(config.BASE_JSON_DATA), &appConfig); err != nil {
			log.Criticalf("JSON encoding failed\n")
			log.ErrorExit(err)
		}
		log.Debugf(config.IS_DEBUG, "Create config.json\n")
		store.Save("config.json", appConfig)
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

func GetAppConfig() AppConfig {
	var appConfig AppConfig
	store.SetApplicationName("pummit")
	if err := store.Load("config.json", &appConfig); err != nil {
		log.Criticalf("Loading alias list failed\n")
		log.ErrorExit(err)
	}
	checkConfig(appConfig)
	return appConfig
}

func resetConfig() {
	var appConfig AppConfig
	store.SetApplicationName("pummit")
	if err := json.Unmarshal([]byte(config.BASE_JSON_DATA), &appConfig); err != nil {
		log.Criticalf("JSON encoding failed\n")
		log.ErrorExit(err)
	}
	log.Debugf(config.IS_DEBUG, "Reset config.json\n")
	if err := store.Save("config.json", appConfig); err != nil {
		log.ErrorExit(err)
	}
}

func GetAliasList() [][]string {
	var config = GetAppConfig()
	return *config.Alias
}

func WriteConfig(configData AppConfig) {
	store.SetApplicationName("pummit")
	if err := store.Save("config.json", configData); err != nil {
		log.ErrorExit(err)
	}
}

// configの値が正常が見るわ！
func checkConfig(appConfig AppConfig) {
	log.Debugf(config.IS_DEBUG, "%+v\n", appConfig)
	// TODO configにプロパティを追加する度に書く必要があるわ！
	if appConfig.UseAlias == nil {
		defaultValue := true
		appConfig.UseAlias = &defaultValue
	}
	if appConfig.WriteEmoji == nil {
		defaultValue := true
		appConfig.WriteEmoji = &defaultValue
	}
	if appConfig.UseLimitPathesLength == nil {
		defaultValue := true
		appConfig.UseLimitPathesLength = &defaultValue
	}
	if appConfig.LimitPathesLength == nil {
		defaultValue := 50
		appConfig.LimitPathesLength = &defaultValue
	}
	if appConfig.Alias == nil {
		defaultValue := [][]string{}
		appConfig.Alias = &defaultValue
	}

	WriteConfig(appConfig)
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

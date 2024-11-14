package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/hectorgimenez/d2go/pkg/data"

	"os"
	"strings"

	"github.com/hectorgimenez/d2go/pkg/data/area"
	"github.com/hectorgimenez/d2go/pkg/data/difficulty"
	"github.com/hectorgimenez/d2go/pkg/data/item"
	"github.com/hectorgimenez/d2go/pkg/data/stat"
	cp "github.com/otiai10/copy"

	"github.com/hectorgimenez/d2go/pkg/nip"

	"gopkg.in/yaml.v3"
)

var (
	Koolo      *KooloCfg
	Characters map[string]*CharacterCfg
	Version    = "dev"
)

type KooloCfg struct {
	ClassicMode bool `yaml:"classicMode"` // Global Classic mode flag

	Debug struct {
		Log         bool `yaml:"log"`
		Screenshots bool `yaml:"screenshots"`
		RenderMap   bool `yaml:"renderMap"`
	} `yaml:"debug"`
	FirstRun              bool   `yaml:"firstRun"`
	UseCustomSettings     bool   `yaml:"useCustomSettings"`
	GameWindowArrangement bool   `yaml:"gameWindowArrangement"`
	LogSaveDirectory      string `yaml:"logSaveDirectory"`
	D2LoDPath             string `yaml:"D2LoDPath"`
	D2RPath               string `yaml:"D2RPath"`
	Discord               struct {
		Enabled                      bool     `yaml:"enabled"`
		EnableGameCreatedMessages    bool     `yaml:"enableGameCreatedMessages"`
		EnableNewRunMessages         bool     `yaml:"enableNewRunMessages"`
		EnableRunFinishMessages      bool     `yaml:"enableRunFinishMessages"`
		EnableDiscordChickenMessages bool     `yaml:"enableDiscordChickenMessages"`
		BotAdmins                    []string `yaml:"botAdmins"`
		ChannelID                    string   `yaml:"channelId"`
		Token                        string   `yaml:"token"`
	} `yaml:"discord"`
	Telegram struct {
		Enabled bool   `yaml:"enabled"`
		ChatID  int64  `yaml:"chatId"`
		Token   string `yaml:"token"`
	}
}

type CharacterCfg struct {
	MaxGameLength   int    `yaml:"maxGameLength"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	AuthMethod      string `yaml:"authMethod"`
	AuthToken       string `yaml:"authToken"`
	Realm           string `yaml:"realm"`
	CharacterName   string `yaml:"characterName"`
	CommandLineArgs string `yaml:"commandLineArgs"`
	KillD2OnStop    bool   `yaml:"killD2OnStop"`
	ClassicMode     bool   `yaml:"classicMode"` // Character-specific Classic mode flag
	CloseMiniPanel  bool   `yaml:"closeMiniPanel"`

	Scheduler Scheduler `yaml:"scheduler"`
	// Other existing fields remain unchanged...
}

// Load reads the config.yaml file and returns a Config struct filled with data from the yaml file
func Load() error {
	Characters = make(map[string]*CharacterCfg)

	// Get the absolute path of the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}

	// Function to get absolute path
	getAbsPath := func(relPath string) string {
		return filepath.Join(cwd, relPath)
	}

	kooloPath := getAbsPath("config/koolo.yaml")
	r, err := os.Open(kooloPath)
	if err != nil {
		return fmt.Errorf("error loading koolo.yaml: %w", err)
	}
	defer r.Close()

	d := yaml.NewDecoder(r)
	if err = d.Decode(&Koolo); err != nil {
		return fmt.Errorf("error reading config %s: %w", kooloPath, err)
	}

	configDir := getAbsPath("config")
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("error reading config directory %s: %w", configDir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		charCfg := CharacterCfg{}
		charConfigPath := getAbsPath(filepath.Join("config", entry.Name(), "config.yaml"))
		r, err = os.Open(charConfigPath)
		if err != nil {
			return fmt.Errorf("error loading config.yaml: %w", err)
		}
		defer r.Close()

		d := yaml.NewDecoder(r)
		if err = d.Decode(&charCfg); err != nil {
			return fmt.Errorf("error reading %s character config: %w", charConfigPath, err)
		}

		// Ensure the Classic mode flag is correctly set
		if Koolo.ClassicMode {
			charCfg.ClassicMode = true
		}

		Characters[entry.Name()] = &charCfg
	}
	for _, charCfg := range Characters {
		charCfg.Validate()
	}

	return nil
}

// Validate ensures the configuration is compatible with the Classic/Expansion modes
func (c *CharacterCfg) Validate() {
	if c.ClassicMode {
		// Ensure compatibility with Classic mode
		c.Character.StashToShared = false // Classic mode doesn't support shared stash
	}
}

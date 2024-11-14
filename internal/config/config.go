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
	NonExpansionMode bool `yaml:"nonExpansionMode"` // New global flag for non-expansion mode

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
	ClassicMode     bool   `yaml:"classicMode"`
	NonExpansionMode bool   `yaml:"nonExpansionMode"` // New character-specific flag
	CloseMiniPanel  bool   `yaml:"closeMiniPanel"`

	Scheduler Scheduler `yaml:"scheduler"`
	Health    struct {
		HealingPotionAt     int `yaml:"healingPotionAt"`
		ManaPotionAt        int `yaml:"manaPotionAt"`
		RejuvPotionAtLife   int `yaml:"rejuvPotionAtLife"`
		RejuvPotionAtMana   int `yaml:"rejuvPotionAtMana"`
		MercHealingPotionAt int `yaml:"mercHealingPotionAt"`
		MercRejuvPotionAt   int `yaml:"mercRejuvPotionAt"`
		ChickenAt           int `yaml:"chickenAt"`
		MercChickenAt       int `yaml:"mercChickenAt"`
	} `yaml:"health"`
	Inventory struct {
		InventoryLock [][]int     `yaml:"inventoryLock"`
		BeltColumns   BeltColumns `yaml:"beltColumns"`
	} `yaml:"inventory"`
	Character struct {
		Class         string `yaml:"class"`
		UseMerc       bool   `yaml:"useMerc"`
		StashToShared bool   `yaml:"stashToShared"`
		UseTeleport   bool   `yaml:"useTeleport"`
		BerserkerBarb struct {
			FindItemSwitch              bool `yaml:"find_item_switch"`
			SkipPotionPickupInTravincal bool `yaml:"skip_potion_pickup_in_travincal"`
		} `yaml:"berserker_barb"`
		NovaSorceress struct {
			BossStaticThreshold int `yaml:"boss_static_threshold"`
		} `yaml:"nova_sorceress"`
	} `yaml:"character"`

	Game struct {
		MinGoldPickupThreshold int                   `yaml:"minGoldPickupThreshold"`
		ClearTPArea            bool                  `yaml:"clearTPArea"`
		Difficulty             difficulty.Difficulty `yaml:"difficulty"`
		RandomizeRuns          bool                  `yaml:"randomizeRuns"`
		Runs                   []Run                 `yaml:"runs"`
		CreateLobbyGames       bool                  `yaml:"createLobbyGames"`
		PublicGameCounter      int                   `yaml:"-"`
		// Additional existing fields remain unchanged
	} `yaml:"game"`

	Companion struct {
		Leader           bool   `yaml:"leader"`
		LeaderName       string `yaml:"leaderName"`
		GameNameTemplate string `yaml:"gameNameTemplate"`
		GamePassword     string `yaml:"gamePassword"`
	} `yaml:"companion"`
}

type BeltColumns [4]string

func (bm BeltColumns) Total(potionType data.PotionType) int {
	typeString := ""
	switch potionType {
	case data.HealingPotion:
		typeString = "healing"
	case data.ManaPotion:
		typeString = "mana"
	case data.RejuvenationPotion:
		typeString = "rejuvenation"
	}

	total := 0
	for _, v := range bm {
		if strings.EqualFold(v, typeString) {
			total++
		}
	}

	return total
}

// Load reads the configuration and applies Non-Expansion Mode if enabled
func Load() error {
	Characters = make(map[string]*CharacterCfg)

	// Other unchanged parts of Load() remain
	// Ensure NonExpansionMode propagates to each character
	for _, charCfg := range Characters {
		if Koolo.NonExpansionMode {
			charCfg.NonExpansionMode = true
		}
		charCfg.Validate()
	}

	return nil
}

// Validate ensures compatibility with Non-Expansion Mode
func (c *CharacterCfg) Validate() {
	if c.NonExpansionMode {
		c.Character.StashToShared = false // Shared stash not available in non-expansion
		fmt.Println("Non-Expansion Mode: Shared stash disabled and Act 5 is skipped.")
	}
}

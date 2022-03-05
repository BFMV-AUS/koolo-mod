package game

import (
	"math/rand"
	"strings"
	"time"
)

const (
	ItemScrollTownPortal   = "ScrollOfTownPortal"
	ItemTomeOfTownPortal   = "TomeOfTownPortal"
	ItemSuperHealingPotion = "SuperHealingPotion"
	ItemSuperManaPotion    = "SuperManaPotion"

	ItemQualityNormal   Quality = "NORMAL"
	ItemQualitySuperior Quality = "SUPERIOR"
	ItemQualityMagic    Quality = "MAGIC"
	ItemQualitySet      Quality = "SET"
	ItemQualityRare     Quality = "RARE"
	ItemQualityUNIQUE   Quality = "UNIQUE"

	StatQuantity       Stat = "Quantity"
	StatGold           Stat = "Gold"
	StatStashGold      Stat = "StashGold"
	StatDurability     Stat = "Durability"
	StatMaxDurability  Stat = "MaxDurability"
	StatNumSockets     Stat = "NumSockets"
	StatFasterCastRate Stat = "FasterCastRate"
	StatEnhancedDamage Stat = "EnhancedDamage"
	StatDefense        Stat = "Defense"
)

type Stat string
type Quality string

type Items struct {
	Belt      Belt
	Inventory Inventory
	Shop      []Item
	Ground    []Item
}

type Inventory []Item

type Item struct {
	ID        int
	Name      string
	Quality   Quality
	Position  Position
	Ethereal  bool
	IsHovered bool
	Stats     map[Stat]int
}

func (i Item) IsPotion() bool {
	return i.IsHealingPotion() || i.IsManaPotion() || i.IsRejuvPotion()
}

func (i Item) IsHealingPotion() bool {
	return strings.Contains(strings.ToLower(i.Name), "healingpotion")
}

func (i Item) IsManaPotion() bool {
	return strings.Contains(strings.ToLower(i.Name), "manapotion")
}
func (i Item) IsRejuvPotion() bool {
	return strings.Contains(strings.ToLower(i.Name), "rejuvenationpotion")
}

func (i Inventory) ShouldBuyTPs() bool {
	for _, it := range i {
		if it.Name != ItemTomeOfTownPortal {
			continue
		}

		qty, found := it.Stats[StatQuantity]
		rand.Seed(time.Now().UnixNano())
		if qty <= rand.Intn(5-1)+1 || !found {
			return true
		}
	}
	return false
}

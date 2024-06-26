package run

import (
	"time"

	"github.com/hectorgimenez/d2go/pkg/data"
	"github.com/hectorgimenez/d2go/pkg/data/area"
	"github.com/hectorgimenez/d2go/pkg/data/object"	
	"github.com/hectorgimenez/d2go/pkg/data/npc"		
	"github.com/hectorgimenez/koolo/internal/action"
	"github.com/hectorgimenez/koolo/internal/game"
	"github.com/hectorgimenez/koolo/internal/action/step"	
)

func (a Rushing) rushAct2() action.Action {
	running := false
	return action.NewChain(func(d game.Data) []action.Action {
		if running || d.PlayerUnit.Area != area.LutGholein {
			return nil
		}

		running = true
		
		if a.CharacterCfg.Game.Rushing.GiveWPs {
			return []action.Action{
				a.builder.VendorRefill(true, false),
				a.GiveAct2WPs(),
				a.getHoradricCube(),
				a.getStaff(),
				a.getAmulet(),
				a.killSummonerQuest(),
				a.killDurielQuest(),	
			}
		}
		
		return []action.Action{
			a.builder.VendorRefill(true, false),
			a.getHoradricCube(),
			a.getStaff(),
			a.getAmulet(),
			a.killSummonerQuest(),
			a.killDurielQuest(),			
		}
	})
}

func (a Rushing) GiveAct2WPs() action.Action {
	areas := []area.ID{
		area.HallsOfTheDeadLevel2,
		area.FarOasis,
		area.LostCity,
		area.ArcaneSanctuary,
		area.CanyonOfTheMagi,
	}

	return action.NewChain(func(d game.Data) []action.Action {
		actions := []action.Action{}

		for _, areaID := range areas {
			actions = append(actions,
				a.builder.WayPoint(areaID),
				a.builder.ClearAreaAroundPlayer(20, data.MonsterAnyFilter()),
				a.builder.OpenTP(),
				a.builder.Wait(time.Second*5),
			)
		}

		return actions
	})
}

func (a Rushing) getHoradricCube() action.Action {
	return action.NewChain(func(d game.Data) []action.Action {
		return []action.Action{
			a.builder.WayPoint(area.HallsOfTheDeadLevel2),
			a.builder.Buff(),
			a.builder.MoveToArea(area.HallsOfTheDeadLevel3),
			a.builder.MoveTo(func(d game.Data) (data.Position, bool) {
				a.logger.Info("Horadric Cube chest found, moving to that room")
				chest, found := d.Objects.FindOne(object.HoradricCubeChest)

			return chest.Position, found
			}),
			a.builder.ClearAreaAroundPlayer(20, data.MonsterAnyFilter()),
			a.builder.OpenTP(),
//			a.waitForParty(d),
			a.builder.ReturnTown(),
		}
	})
}

func (a Rushing) getStaff() action.Action {
	return action.NewChain(func(d game.Data) []action.Action {
		return []action.Action{
			a.builder.WayPoint(area.FarOasis),
			a.builder.Buff(),			
			a.builder.MoveToArea(area.MaggotLairLevel1),
			a.builder.MoveToArea(area.MaggotLairLevel2),
			a.builder.MoveToArea(area.MaggotLairLevel3),
			a.builder.MoveTo(func(d game.Data) (data.Position, bool) {
				a.logger.Info("Staff Of Kings chest found, moving to that room")
				chest, found := d.Objects.FindOne(object.StaffOfKingsChest)

				return chest.Position, found
			}),
			a.builder.ClearAreaAroundPlayer(20, data.MonsterAnyFilter()),
			a.builder.OpenTP(),			
//			a.waitForParty(d),
			a.builder.ReturnTown(),
		}
	})
}

func (a Rushing) getAmulet() action.Action {
	return action.NewChain(func(d game.Data) []action.Action {
		return []action.Action{
			a.builder.WayPoint(area.LostCity),
			a.builder.Buff(),
			a.builder.MoveToArea(area.ValleyOfSnakes),
			a.builder.MoveToArea(area.ClawViperTempleLevel1),
			a.builder.MoveToArea(area.ClawViperTempleLevel2),
			a.builder.MoveTo(func(d game.Data) (data.Position, bool) {
				a.logger.Info("Altar found, moving closer")
				chest, found := d.Objects.FindOne(object.TaintedSunAltar)

				return chest.Position, found
			}),
			a.builder.ClearAreaAroundPlayer(20, data.MonsterAnyFilter()),
			a.builder.OpenTP(),
//			a.waitForParty(d),
			a.builder.ReturnTown(),			
		}
	})
}

func (a Rushing) killSummonerQuest() action.Action {
	return action.NewChain(func(d game.Data) []action.Action {
		return []action.Action{
			a.builder.WayPoint(area.ArcaneSanctuary), 
			a.builder.OpenTP(),			
			a.builder.Buff(),

			a.builder.MoveTo(func(d game.Data) (data.Position, bool) {
				if summoner, found := d.NPCs.FindOne(npc.Summoner); found {
					return summoner.Positions[0], true
				}
				return data.Position{}, false
			}, step.StopAtDistance(80)),

			a.builder.OpenTP(),			
//			a.waitForParty(d),
			a.char.KillSummoner(),
			a.builder.ReturnTown(),			
		}
	})
}

func (a Rushing) killDurielQuest() action.Action {
	var realTomb area.ID

	for _, tomb := range talRashaTombs {
		_, _, objects, _ := a.Reader.CachedMapData.NPCsExitsAndObjects(data.Position{}, tomb)
		for _, obj := range objects {
			if obj.Name == object.HoradricOrifice {
				realTomb = tomb
				break
			}
		}
	}

	if realTomb == 0 {
		a.logger.Info("Could not find the real tomb :(")
		return nil
	}

	return action.NewChain(func(d game.Data) []action.Action {
		var actions []action.Action

		actions = append(actions,
			a.builder.WayPoint(area.CanyonOfTheMagi),
			a.builder.OpenTP(),
			a.builder.Buff(),
			a.builder.MoveToArea(realTomb),
			a.builder.MoveTo(func(d game.Data) (data.Position, bool) {
				orifice, found := d.Objects.FindOne(object.HoradricOrifice)
				if found {
					return orifice.Position, true
				}
				return data.Position{}, false
			}),
			a.builder.ClearAreaAroundPlayer(15, data.MonsterAnyFilter()),
			a.builder.OpenTP(),
			// a.waitForParty(d),
			a.builder.Buff(),
		)

		actions = append(actions,
		action.NewChain(func(d game.Data) []action.Action {			
				_, found := d.Objects.FindOne(object.DurielsLairPortal)
				if found {
					return []action.Action{
						a.builder.InteractObject(object.DurielsLairPortal, func(d game.Data) bool {
							return d.PlayerUnit.Area == area.DurielsLair
						}),
					}
				}
				return nil
			}),
		)

		actions = append(actions,
			a.builder.MoveToArea(area.DurielsLair),
			a.char.KillDuriel(),
			a.builder.ReturnTown(),
			a.builder.WayPoint(area.KurastDocks),
		)

		return actions
	})
}

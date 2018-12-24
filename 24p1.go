package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// splits a string, trims spaces on every element
func splitandclean(in, sep string, n int) []string {
	v := strings.SplitN(in, sep, n)
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}
	return v
}

// convert string to integer
func atoi(in string) int {
	n, err := strconv.Atoi(in)
	must(err)
	return n
}

var unitre = regexp.MustCompile(`(\d+) units each with (\d+) hit points (\(.*?\))? ?with an attack that does (\d+) ([^\s]*?) damage at initiative (\d+)`)

type Army struct {
	name   string
	groups []Group
}

type Group struct {
	idx         int
	nunits      int
	starthp     int
	attackPower int
	attackType  string
	initiative  int
	weakTo      []string
	immuneTo    []string
	parentArmy  *Army

	tgt      *Group
	selected bool
}

var ImmuneSystem Army
var Infection Army

func main() {
	ImmuneSystem.name = "Immune System"
	Infection.name = "Infection"
	buf, err := ioutil.ReadFile("24.txt")
	must(err)
	var curArmy *Army
	for _, line := range strings.Split(string(buf), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "Immune System:" {
			curArmy = &ImmuneSystem
			continue
		}
		if line == "Infection:" {
			curArmy = &Infection
			continue
		}
		v := unitre.FindStringSubmatch(line)
		v = v[1:]
		fmt.Printf("%q\n", v[1:])
		nunits := atoi(v[0])
		starthp := atoi(v[1])
		attackPower := atoi(v[3])
		attackType := v[4]
		initiative := atoi(v[5])

		weakTo := []string{}
		immuneTo := []string{}

		if v[2] != "" {
			v[2] = v[2][1 : len(v[2])-1]
			for _, str := range splitandclean(v[2], ";", -1) {
				const (
					weakToPrefix   = "weak to "
					immuneToPrefix = "immune to "
				)
				var curKind *[]string
				if strings.HasPrefix(str, weakToPrefix) {
					curKind = &weakTo
					str = str[len(weakToPrefix):]
				} else if strings.HasPrefix(str, immuneToPrefix) {
					curKind = &immuneTo
					str = str[len(immuneToPrefix):]
				} else {
					curKind = nil
				}

				*curKind = append(*curKind, splitandclean(str, ",", -1)...)
			}
		}

		fmt.Printf("\tweakTo: %q\n\timmuneTo: %q\n", weakTo, immuneTo)

		curArmy.groups = append(curArmy.groups, Group{
			idx:         len(curArmy.groups) + 1,
			nunits:      nunits,
			starthp:     starthp,
			attackPower: attackPower,
			attackType:  attackType,
			initiative:  initiative,
			weakTo:      weakTo,
			immuneTo:    immuneTo,
			parentArmy:  curArmy,
		})
	}

	fmt.Printf("\n")

	battle()
}

const debugPrintStatus = true
const debugTargetSelection = true
const debugAttack = true

func battle() {
	for {
		if debugPrintStatus {
			fmt.Printf("Immune System:\n")
			ImmuneSystem.Describe()
			fmt.Printf("Infection:\n")
			Infection.Describe()
			fmt.Printf("\n")
		}

		selectTarget()
		attack()

		if victory() {
			break
		}
	}

	fmt.Printf("\n")

	fmt.Printf("Immune System: %d units\n", ImmuneSystem.totalUnits())
	fmt.Printf("infection: %d units\n", Infection.totalUnits())
}

func selectTarget() {
	allGroups := make([]*Group, 0, len(ImmuneSystem.groups)+len(Infection.groups))
	for i := range ImmuneSystem.groups {
		allGroups = append(allGroups, &ImmuneSystem.groups[i])
	}
	for i := range Infection.groups {
		allGroups = append(allGroups, &Infection.groups[i])
	}

	for _, group := range allGroups {
		group.selected = false
		group.tgt = nil
	}

	sort.Slice(allGroups, func(i, j int) bool {
		epi := allGroups[i].effectivePower()
		epj := allGroups[j].effectivePower()
		if epi == epj {
			return allGroups[i].initiative > allGroups[j].initiative
		}
		return epi > epj
	})

	for _, group := range allGroups {
		if group.nunits <= 0 {
			continue
		}
		maxDmg := 0
		var tgtGroup *Group
		for _, tgt := range allGroups {
			if tgt.selected || tgt.parentArmy.name == group.parentArmy.name || tgt.nunits <= 0 {
				continue
			}
			dmg := group.damageFor(tgt)
			if dmg == 0 {
				continue
			}
			if debugTargetSelection {
				fmt.Printf("%s group %d would deal defending group %d %d damage\n", group.parentArmy.name, group.idx, tgt.idx, dmg)
			}
			if dmg == maxDmg {
				if tgt.effectivePower() == tgtGroup.effectivePower() {
					if tgt.initiative > tgtGroup.initiative {
						tgtGroup = tgt
					}
				} else if tgt.effectivePower() > tgtGroup.effectivePower() {
					tgtGroup = tgt
				}
			} else if dmg > maxDmg {
				tgtGroup = tgt
				maxDmg = dmg
			}
		}

		if tgtGroup != nil {
			tgtGroup.selected = true
		}
		group.tgt = tgtGroup
	}

	if debugTargetSelection {
		fmt.Printf("\n")
	}
}

func attack() {
	allGroups := make([]*Group, 0, len(ImmuneSystem.groups)+len(Infection.groups))
	for i := range ImmuneSystem.groups {
		allGroups = append(allGroups, &ImmuneSystem.groups[i])
	}
	for i := range Infection.groups {
		allGroups = append(allGroups, &Infection.groups[i])
	}

	sort.Slice(allGroups, func(i, j int) bool {
		return allGroups[i].initiative > allGroups[j].initiative
	})

	for _, group := range allGroups {
		if group.nunits <= 0 {
			continue
		}
		if group.tgt == nil {
			continue
		}
		kills := 0
		dmg := group.damageFor(group.tgt)
		for dmg >= group.tgt.starthp && group.tgt.nunits > 0 {
			dmg -= group.tgt.starthp
			group.tgt.nunits--
			kills++
		}

		if debugAttack {
			fmt.Printf("%s group %d attacks defending group %d, killing %d units\n", group.parentArmy.name, group.idx, group.tgt.idx, kills)
		}
	}

	if debugAttack {
		fmt.Printf("\n")
	}
}

func victory() bool {
	return (ImmuneSystem.totalUnits() == 0) || (Infection.totalUnits() == 0)
}

func (army *Army) totalUnits() int {
	r := 0
	for i := range army.groups {
		r += army.groups[i].nunits
	}
	return r
}

func (group *Group) effectivePower() int {
	return group.nunits * group.attackPower
}

func (group *Group) damageFor(tgt *Group) int {
	if tgt.isImmune(group.attackType) {
		return 0
	}
	ep := group.effectivePower()
	if tgt.isWeak(group.attackType) {
		return ep * 2
	}
	return ep
}

func (group *Group) isImmune(attackType string) bool {
	for _, imm := range group.immuneTo {
		if imm == attackType {
			return true
		}
	}
	return false
}

func (group *Group) isWeak(attackType string) bool {
	for _, weak := range group.weakTo {
		if weak == attackType {
			return true
		}
	}
	return false
}

func (army *Army) Describe() {
	for i := range army.groups {
		group := &army.groups[i]
		fmt.Printf("Group %d contains %d units\n", group.idx, group.nunits)
	}
}

/*
Immune System:
Group 1 91'630
Group 2 1'259'986

Infection:
Group 1 3'769'506
Group 2 13'280'085
*/

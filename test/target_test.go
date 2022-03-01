package main

import (
	"afkl/fumofinger/model"
	"fmt"
	"testing"
)

func TestLoadTargetFromFile(t *testing.T) {
	targets := new(model.TargetList)
	testTargets := [3]string{
		"http://bilibili.com/",
		"https://fumo.website/",
		"http://thwiki.cc",
	}
	targetNum, err := targets.LoadTargetFromFile("./assets/targets.txt")
	fmt.Println(targets)
	if targetNum != len(testTargets) || err != nil {
		t.Errorf("Loaded Error: %s", err.Error())
	} else {
		for id, target := range testTargets {
			if target != targets.List[id] {
				t.Error("Loaded for different target")
			}
			fmt.Println(target, targets.List[id])
		}
	}
}

func TestSplitTarget(t *testing.T) {
	targetStrings := [4]string{
		"http://bilibili.com/",
		"https://fumo.website/",
		"https://thwiki.cc/",
		"https://www.baidu.com/",
	}
	targets := &model.TargetList{
		List: targetStrings[:],
		Len:  len(targetStrings),
	}

	fmt.Println(targets.SplitTargetList(5))
}

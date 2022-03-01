package model

import (
	"afkl/fumofinger/assets"
	"encoding/json"
	"strings"
)

type FoFaRule struct {
	Match   string `json:"match"`
	Content string `json:"content"`
}

type FoFaFingerBlock struct {
	RuleID         string       `json:"rule_id"`
	Level          string       `json:"level"`
	Softhard       string       `json:"softhard"`
	Product        string       `json:"product"`
	Company        string       `json:"company"`
	Category       string       `json:"category"`
	ParentCategory string       `json:"parent_category"`
	Rules          [][]FoFaRule `json:"rules"`
}

type FoFaFinger []FoFaFingerBlock

var FoFa FoFaFinger

func init() {
	err := json.Unmarshal(assets.FoFaFingerData, &FoFa)
	if err != nil {
		panic(err.Error())
	}
}

/**
 * captureType:
 * 0 => banner
 * 1 => title
 */
func captureBody(resp MatchResponse, rule FoFaRule, captureType int) bool {
	var re string
	switch captureType {
	case 0:
		re = `(?im)<\s*banner.*>(.*?)<\s*/\s*banner>`
	case 1:
		re = `(?im)<\s*title.*>(.*?)<\s*/\s*title>`
	}

	matchResults := ReMatchBody(resp, re)
	if len(matchResults) == 0 {
		return false
	}

	for _, matchResult := range matchResults {
		if !strings.Contains(
			strings.ToLower(matchResult),
			strings.ToLower(rule.Content),
		) {
			return false
		}
	}

	return true
}

func Capture(resp MatchResponse, splitFoFaFinger FoFaFinger) ([]FoFaFingerBlock, bool) {
	matchFinger := []FoFaFingerBlock{}
	isSubRuleMatched := true
	isAllSubRulesMatched := true
	for _, block := range splitFoFaFinger {
		for _, rule := range block.Rules {
			for _, subRule := range rule {
				switch strings.Split(subRule.Match, "_")[0] {
				case "banner":
					isSubRuleMatched = captureBody(resp, subRule, 0)
				case "title":
					isSubRuleMatched = captureBody(resp, subRule, 1)
				case "body":
					if !strings.Contains(
						strings.ToLower(string(resp.Body)),
						strings.ToLower(subRule.Content),
					) {
						isSubRuleMatched = false
					}
				case "header":
					if resp.Headers.Get(subRule.Content) == "" {
						isSubRuleMatched = false
					}
				case "server":
					if !strings.Contains(
						strings.ToLower(resp.Headers.Get("Server")),
						strings.ToLower(subRule.Content),
					) {
						isSubRuleMatched = false
					}
				case "cert":
					if (resp.Cert == nil) || (resp.Cert != nil && !strings.Contains(string(resp.Cert), subRule.Content)) {
						isSubRuleMatched = false
					}
				case "protocol":
					// TODO
					isSubRuleMatched = false
				default:
					// TODO
					isSubRuleMatched = false
				}
				/** 所有子规则为 与 条件，那么这里就是不断与上次的结果进行与运算
				 *	只要有一个子规则不符合，那么 `isAllSubRulesMatched` 就会为
				 *	false 从而跳出子规则的判断
				 */
				isAllSubRulesMatched = isAllSubRulesMatched && isSubRuleMatched
				if !isAllSubRulesMatched {
					break
				}
			}
			/** 同理，主规则和主规则直接为 或 条件，只要上面的子规则
			 *	全部通过，`isAllSubRulesMatched` 就为 true。说明找
			 *	到了对应的规则，直接返回
			 */
			if isAllSubRulesMatched {
				matchFinger = append(matchFinger, block)
			}
			isSubRuleMatched = true
			isAllSubRulesMatched = true
		}
	}
	if len(matchFinger) != 0 {
		return matchFinger, true
	}
	return nil, false
}

func (fingerList FoFaFinger) SplitTargetList(blockNum int) []FoFaFinger {
	max := len(fingerList)
	if max < blockNum {
		blockNum = max
	}
	var segmens = make([]FoFaFinger, 0)
	quantity := max / blockNum
	end := 0
	for i := 1; i <= blockNum; i++ {
		qu := i * quantity
		if i != blockNum {
			segmens = append(segmens, fingerList[i-1+end:qu])
		} else {
			segmens = append(segmens, fingerList[i-1+end:])
		}
		end = qu - i
	}
	return segmens
}

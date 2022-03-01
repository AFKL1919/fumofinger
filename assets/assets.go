package assets

import (
	_ "embed"
)

//go:embed fofa.json
var FoFaFingerData []byte

//go:embed koishi.txt
var NobodySeeingKoishi string

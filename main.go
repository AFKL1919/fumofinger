package main

import (
	"afkl/fumofinger/assets"
	"afkl/fumofinger/core"
	"afkl/fumofinger/model"
	"log"
	"os"

	cli "github.com/jawher/mow.cli"
)

var VERSION string = "FUMO-0.0.1"

func main() {
	app := cli.App("FumoFinger", "OMG! The FUMO! ᗜˬᗜ")
	app.Version("v version", VERSION)
	var (
		targets      = app.StringsOpt("t target", []string{}, "Set the target url.")
		targetsFiles = app.StringsOpt("r file", []string{}, "Set the target list from file.")
		thread       = app.IntOpt("threads", 10, "Threads number.")
		// timeout    = app.IntOpt("timeout", 20, "Request timeout.")
		fumo = app.StringOpt("fumo", "fumo", "fumo")
	)
	// app.Spec = "-v | (-t=<target> | -r=<target file>) [--threads=<threads>] [--timeout=<timeout>] | --fumo=<fumo>"
	app.Spec = "-v | --fumo=<fumo> | (-t=<target> | -r=<target file>)... [--threads=<threads>]"
	app.Action = func() {
		if *fumo == "fumo" {
			targetNum := 0
			execTargets := new(model.TargetList)
			if len(*targets) != 0 {
				for _, target := range *targets {
					execTargets.LoadTargetFromString(target)
					targetNum += 1
				}
			}

			if len(*targetsFiles) != 0 {
				for _, file := range *targetsFiles {
					num, err := execTargets.LoadTargetFromFile(file)
					if err != nil {
						log.Fatalln(err.Error())
					}
					targetNum += num
				}
			}

			if targetNum == 0 {
				log.Fatalln("target number is 0")
			}

			core.Start(*execTargets, *thread)
			core.Wait()
			core.End()
		} else {
			log.Println(assets.NobodySeeingKoishi)
		}
	}
	app.Run(os.Args)
}

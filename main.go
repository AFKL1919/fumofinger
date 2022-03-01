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
		target     = app.StringOpt("t target", "", "Set the target url.")
		targetFile = app.StringOpt("r file", "", "Set the target list from file.")
		thread     = app.IntOpt("threads", 10, "Threads number.")
		// timeout    = app.IntOpt("timeout", 20, "Request timeout.")
		fumo = app.StringOpt("fumo", "fumo", "fumo")
	)
	// app.Spec = "-v | (-t=<target> | -r=<target file>) [--threads=<threads>] [--timeout=<timeout>] | --fumo=<fumo>"
	app.Spec = "-v | --fumo=<fumo> | (-t=<target> | -r=<target file>) [--threads=<threads>]"
	app.Action = func() {
		if *fumo == "fumo" {
			targetNum := 0
			targets := new(model.TargetList)
			if len(*target) != 0 {
				targets.LoadTargetFromString(*target)
				targetNum += 1
			}

			if len(*targetFile) != 0 {
				num, err := targets.LoadTargetFromFile(*targetFile)
				if err != nil {
					log.Fatalln(err.Error())
				}
				targetNum += num
			}

			if targetNum == 0 {
				log.Fatalln("target number is 0")
			}

			core.Start(*targets, *thread)
			core.Wait()
			core.End()
		} else {
			log.Println(assets.NobodySeeingKoishi)
		}
	}
	app.Run(os.Args)
}

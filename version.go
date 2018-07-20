package main

import (
	"log"

	"github.com/logrusorgru/aurora"
)

var version = "0.1.5"

func banner() {
	log.Println("\n\n  ██████╗██████╗ ███████╗██╗    ██╗    ███████╗███████╗██████╗ ██╗   ██╗███████╗██████╗")
	log.Println(" ██╔════╝██╔══██╗██╔════╝██║    ██║    ██╔════╝██╔════╝██╔══██╗██║   ██║██╔════╝██╔══██╗")
	log.Println(" ██║     ██████╔╝█████╗  ██║ █╗ ██║    ███████╗█████╗  ██████╔╝██║   ██║█████╗  ██████╔╝")
	log.Println(" ██║     ██╔══██╗██╔══╝  ██║███╗██║    ╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝██╔══╝  ██╔══██╗")
	log.Println(" ╚██████╗██║  ██║███████╗╚███╔███╔╝    ███████║███████╗██║  ██║ ╚████╔╝ ███████╗██║  ██║")
	log.Println("  ╚═════╝╚═╝  ╚═╝╚══════╝ ╚══╝╚══╝     ╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝")
	log.Printf("\n                        This is %s version %s \n", aurora.BgBlue("Crew Server").Bold(), aurora.BgBlue(version).Bold())
	log.Printf("                               Made with %s by %s\n\n", aurora.Red("♥").Bold(), "Steps")
}

package main

import (
	"fmt"
	"os"

	"github.com/lutfiharidha/pokemon/app"
	"github.com/lutfiharidha/pokemon/db"
)

func main() {
	m := app.NewLogModels()
	m.Init()
	dbnew := db.NewSQL().SetupDatabaseConnection() //setup database connection
	db.Migrator(dbnew)

	args := os.Args
	if len(args) > 1 {

		switch args[1] {
		case "battle":
			fmt.Print("How much games? ")
			var times int
			fmt.Scanln(&times)
			if times != 0 {
				app.GoBattle(times)
			}
		case "pokedex-s":
			fmt.Print("Name Pokedex: ")
			var name string
			fmt.Scanln(&name)
			app.GetPokedexInfo(name)
		case "pokedex":
			app.GetPokedexInfo("")
		case "log-i":
			fmt.Print("Interval: ")
			var times string
			fmt.Scanln(&times)
			m.IntervalLog(times)
		case "log-d":
			fmt.Print("Date (yyyy/mm/dd): ")
			var times string
			fmt.Scanln(&times)
			m.SpecificLog(times)
		case "log-si":
			fmt.Print("From Date (yyyy/mm/dd): ")
			var from string
			fmt.Scanln(&from)

			fmt.Print("To Date (yyyy/mm/dd): ")
			var to string
			fmt.Scanln(&to)
			m.SpecificIntervalLog(from, to)
		case "help":
			fmt.Printf("\nPokémon Battle Royale\n\n")
			fmt.Printf("Usage:\n")
			fmt.Printf("make battle:\n")
			fmt.Printf("This command to start the Pokémon Battle Royale.\n\n")
			fmt.Printf("make pokedex:\n")
			fmt.Printf("This command is to view pokedex data that specifies how many pokemon are in pokedex.\n\n")
			fmt.Printf("make pokedex-s:\n")
			fmt.Printf("This command is to see all pokedex data .\n\n")
			fmt.Printf("make log-i:\n")
			fmt.Printf("This command for creating an interval log from this day\n")
			fmt.Printf("there required to input the time interval for an example\n")
			fmt.Printf("the input 1d which means 1 day for the month change \"d\" to \"m\" and the year use \"y\".\n\n")
			fmt.Printf("make log-d:\n")
			fmt.Printf("This command for creating a specific date log.\n\n")
			fmt.Printf("make log-si:\n")
			fmt.Printf("This command for creating a specific interval log.\n\n")
		}
	} else {
		fmt.Println("are you lost? you can try \"make help\"")
	}

}

package app

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lutfiharidha/pokemon/app/types"
	"github.com/mtslzr/pokeapi-go"
	"github.com/xujiajun/nutsdb"
)

type FirstBattle struct {
	Source *PokemonList
	Target *PokemonList
}
type PokemonList struct {
	name        string
	hp          int64
	move        []types.Moves
	numOfBattle int64
	point       int64
}

var listPokemon = []*PokemonList{}

var pokelist = make(map[int]bool)
var mu sync.RWMutex

var pokeMemberBattle []*PokemonList
var playerOut = []*PokemonList{}
var logging = types.Logging{}
var winner *PokemonList

var pokedexList []types.Pokedex

func GoBattle(times int) {
	//start the battle!
	for i := 0; i < times; i++ {
		showLog("\n======== GAME " + fmt.Sprint(i+1) + " ========\n")
		run()
		listPokemon = []*PokemonList{}
		pokelist = make(map[int]bool)
		pokeMemberBattle = []*PokemonList{}
		playerOut = []*PokemonList{}
		winner = &PokemonList{}
	}

}

func run() {
	if len(pokedexList) == 0 {
		//open connection to nutsDB
		db, err := nutsdb.Open(
			nutsdb.DefaultOptions,
			nutsdb.WithDir("./cache"),
		)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		//check if the pokedex data is already stored in the db
		if err := db.View(
			func(tx *nutsdb.Tx) error {
				g, _ := pokeapi.Resource("pokedex", 0, 30)
				for _, v := range g.Results {
					bucket := "pokedex"
					entries, err := tx.Get(bucket, []byte(v.Name))
					if err != nil {
						return err
					}
					pokemons := []types.Pokemon{}
					json.Unmarshal(entries.Value, &pokemons)
					pokedexList = append(pokedexList, types.Pokedex{
						Name:    v.Name,
						Pokemon: pokemons,
						Stock:   len(pokemons),
					})
				}
				return nil
			}); err != nil {
			//get data pokedex from api if the pokedek data isn't already stored in the db
			pokedexList, _ = getDataFromApi(db)
		}
	}

	rand.Seed(time.Now().UnixNano())

	battle(5)

	time.Sleep(1 * time.Second)
	getWinner()

}

func battle(noOfPoke int) {

	//start using goroutines to speed up the process of finding pokemon data
	var wg sync.WaitGroup
	wg.Add(noOfPoke)
	showLog("\nGET DATA POKEMONS... \n\n")
	for i := 0; i < noOfPoke; i++ {
		go func(wg *sync.WaitGroup) {
			findPoketoGame()
			wg.Done()
		}(&wg)
	}
	wg.Wait()

	//anyway showLog() is a customized function to display a log and mapping a log
	for _, v := range listPokemon {
		showLog(strings.ToUpper(v.name) + " with hp " + fmt.Sprint(v.hp) + "\n")
	}

	//map the pokemon data that is still in the arena
	for k := range pokelist {
		ll := listPokemon[k]
		if ll == nil {
			panic("pokemon not found")
		}
		pokeMemberBattle = append(pokeMemberBattle, listPokemon[k])
	}

	//the fight process
	showLog("\n======== FIGHT! ========\n\n")
	for { // loop until the game has a winner
		var firstb []*FirstBattle
		for _, s := range pokeMemberBattle {
			mu.Lock()
			enemy := FindEnemy(s.name)                                      //finding enemies
			firstb = append(firstb, &FirstBattle{Source: s, Target: enemy}) //mapping the data of who is fighting
			mu.Unlock()

		}
		// Start using goroutine so that 5 pokemon can fight at once
		var wgs sync.WaitGroup
		wgs.Add(len(pokeMemberBattle))

		for _, s := range firstb {
			go func(wgs *sync.WaitGroup, s *FirstBattle) {
				runPoke(s.Source, s.Target) //math function (hp-power)
				wgs.Done()
			}(&wgs, s)
		}
		wgs.Wait()

		//if there is only one pokemon left in the arena, he/she is the winner
		if len(pokeMemberBattle) == 1 {
			last := pokeMemberBattle[0]
			atomic.AddInt64(&last.point, 5)
			winner = last
			break
		}

		//if there is not a single pokemon left in the arena, take the last one out of the arena and he/she is the winner
		if len(pokeMemberBattle) == 0 {
			last := playerOut[len(playerOut)-1]
			atomic.AddInt64(&last.point, 5)
			winner = last
			break
		}

		//but if there are still more than 1 pokemon in the arena, the game continues with the pokemon looping to find who is still in the arena.
		var tmplist []*PokemonList
		for _, s := range pokeMemberBattle {
			if !findPlayer(playerOut, s.name) {
				tmplist = append(tmplist, s)
			}
		}

		pokeMemberBattle = tmplist //update pokemon list in arena

		if len(pokeMemberBattle) > 1 {
			showLog(fmt.Sprint(len(pokeMemberBattle)) + " pokemon left!\n")
		}

		time.Sleep(100 * time.Millisecond)

	}

}

// just function to search pokemon in list
func findPlayer(list []*PokemonList, player string) bool {
	for _, v := range list {
		if player == v.name {
			return true
		}
	}
	return false
}

func runPoke(p *PokemonList, enemy *PokemonList) {

	if enemy == nil {
		return
	}

	if enemy != nil {
		//record how many pokemon movements.
		atomic.AddInt64(&p.numOfBattle, 1)
		atomic.AddInt64(&enemy.numOfBattle, 1)

		mu.Lock()
		//pick random data move from pokemon
		ra := p.move[randomNum(0, len(p.move)-1)]
		ea := enemy.move[randomNum(0, len(enemy.move)-1)]

		//this condition just for log:)
		if ra.Power == 0 || ea.Power == 0 {
			powerName := ea.Name
			pokeName := enemy.name
			pokePower := ea.Power
			if ra.Power == 0 {
				powerName = ra.Name
				pokeName = p.name
				pokePower = ra.Power

			}
			showLog(strings.ToUpper(pokeName) + " using " + powerName + " power " + fmt.Sprint(pokePower) + "\n")

		} else {
			showLog(strings.ToUpper(p.name) + " using " + ra.Name + " with power " + fmt.Sprint(ra.Power) + " to hit " + strings.ToUpper(enemy.name) + "\n")
			showLog(strings.ToUpper(enemy.name) + " using " + ea.Name + " with power " + fmt.Sprint(ea.Power) + " to hit " + strings.ToUpper(p.name) + "\n")

			//Reduce pokemon HP as much as the enemy power received by it
			atomic.StoreInt64(&p.hp, (p.hp - int64(ea.Power)))
			atomic.StoreInt64(&enemy.hp, (enemy.hp - int64(ra.Power)))
			showLog(strings.ToUpper(p.name) + " hp left " + fmt.Sprint(p.hp-int64(ea.Power)) + " and " + strings.ToUpper(enemy.name) + " hp left " + fmt.Sprint(p.hp-int64(ea.Power)) + "\n")

		}
		mu.Unlock()

		mu.Lock()
		//this condition is to map which pokemon data has run out of HP.
		if p.hp <= 0 {
			if !findPlayer(playerOut, p.name) {
				playerOut = append(playerOut, p) //if HP is up, register to the playerout list.
				showLog("[" + strings.ToUpper(p.name) + " out from arena]\n")
			}

		}
		//same condition but this for the opponent
		if enemy.hp <= 0 {
			if !findPlayer(playerOut, enemy.name) {
				playerOut = append(playerOut, enemy)
				showLog("[" + strings.ToUpper(enemy.name) + " out from arena]\n")
			}
		}
		mu.Unlock()

	}
}

func getWinner() error {
	showLog("\n======== GAME DONE! ========\n\n")
	showLog("\nResults:\n")

	//give losers a point to appreciate their effort:)
	if len(listPokemon) > 1 {
		for _, last := range playerOut {
			if last.name != winner.name {
				atomic.AddInt64(&last.point, 1)
			}
		}
	}

	//sort by point to log a result
	sort.SliceStable(listPokemon, func(i, j int) bool {
		mi, mj := listPokemon[i], listPokemon[j]
		return mi.point > mj.point
	})
	for _, s := range listPokemon {
		showLog(strings.ToUpper(s.name) + " with hp " + fmt.Sprint(s.hp) + " left " + fmt.Sprint(s.numOfBattle) + " move(s) and with point " + fmt.Sprint(s.point) + "\n")
	}

	showLog("\nTHE WINNER IS" + strings.ToUpper(winner.name) + "!\n")

	logging.Winner = strings.ToUpper(winner.name)
	if winner.point != 5 {
		log.Panic("wrong ....")
	}
	logging.DataLog = dataLog
	m := NewLogModels()
	m.Init()
	m.SaveLog(logging)

	return nil

}

func FindEnemy(name string) *PokemonList {
	if len(pokeMemberBattle) == 1 {
		return nil
	}

	if len(pokeMemberBattle) == 0 {
		return nil
	}

	n := randomNum(0, len(pokeMemberBattle)-1)
	//pick a random pokemon
	dd := pokeMemberBattle[n]

	if dd == nil {
		return nil
	}
	//if the pokemon finds an opponent himself/herself search until he/she finds it.
	if dd.name == name {
		return FindEnemy(name)
	}
	return pokeMemberBattle[n]
}

func findPoketoGame() bool {
	if len(pokedexList) < 1 {
		return false
	}
	n := randomNum(0, len(pokedexList)-1)

	if pokelist[n] {
		return findPoketoGame()
	}

	if !pokelist[n] {
		randNum := randomNum(0, len(pokedexList)-1)
		move, err := getMoves(pokedexList[randNum].Pokemon[randomNum(0, len(pokedexList[randNum].Pokemon)-1)].Name) //mapping pokemon's move
		//if error then the process will be repeated
		if err != nil {
			return findPoketoGame()
		}

		//if there is no move data at all, REPEAT UNTIL FOUND! :)
		if len(move.Moves) == 0 {
			return findPoketoGame()
		}
		//if there are pokemon movement data, then map the data of the pokemon list that will play.
		if len(move.Moves) != 0 {
			mu.Lock()
			listPokemon = append(listPokemon, &PokemonList{
				name: move.Name,
				hp:   int64(move.Hp),
				move: move.Moves,
			})
			pokelist[len(pokelist)] = true
			mu.Unlock()
			return true
		}
	}

	return false
}

func randomNum(min, max int) int {
	return rand.Intn(max-min+1) + min
}

package app

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/lutfiharidha/pokemon/app/types"
	"github.com/mtslzr/pokeapi-go"
	"github.com/xujiajun/nutsdb"
)

// function to get detail pokedex
func GetPokedexInfo(data string) {
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
					Name:  v.Name,
					Stock: len(pokemons),
				})
			}
			return nil
		}); err != nil {
		//get data pokedex from api if the pokedek data isn't already stored in the db
		pokedexList, _ = getDataFromApi(db)
	}
	for _, v := range pokedexList {
		if data == "" {
			fmt.Printf("Pokedex %s has %d pokemons\n", v.Name, v.Stock)
		} else {
			if v.Name == data {
				fmt.Printf("Pokedex %s has %d pokemons\n", v.Name, v.Stock)
				break
			}

		}
	}
}

func getDataFromApi(db *nutsdb.DB) ([]types.Pokedex, error) {

	var pokedex []types.Pokedex

	g, _ := pokeapi.Resource("pokedex", 0, 30) //get all pokedex data from api

	for _, v := range g.Results {
		//get pokedex's name
		pokedex = append(pokedex, types.Pokedex{
			Name: v.Name,
		})
	}

	tx, err := db.Begin(true) //open transaction to nutsDB
	if err != nil {
		log.Fatal(err)
	}
	//start using goroutine to make the process faster
	var wg sync.WaitGroup
	wg.Add(len(pokedex)) //add workers up to the amount of pokedex data
	for i, v := range pokedex {
		var key, val []byte
		go func(wg *sync.WaitGroup, i int, v types.Pokedex) {
			pokedexData, _ := pokeapi.Pokedex(v.Name) //get detail pokedex

			for _, b := range pokedexData.PokemonEntries {
				//get pokemon name from pokedex
				pokedex[i].Pokemon = append(pokedex[i].Pokemon, types.Pokemon{
					Name: b.PokemonSpecies.Name,
				})

			}

			key = []byte(v.Name)
			val, _ = json.Marshal(pokedex[i].Pokemon)

			//start to put data to nutsDB with bucket name "pokedex", key name is pokedex name, and the value is pokemon in each pokedex
			if err := tx.Put("pokedex", key, val, nutsdb.Persistent); err != nil {
				//rollback the transaction.
				tx.Rollback()
			}

			pokedex[i].Stock = len(pokedex[i].Pokemon) //don't forget to save the number of stock pokemon each Pokedex has.
			wg.Done()
		}(&wg, i, v)
	}
	wg.Wait()

	//commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		tx.Rollback()
	}

	return pokedex, nil
}

func getMoves(pokemon string) (res types.Pokemon, err error) {
	pokemonData, err := pokeapi.Pokemon(pokemon) //get pokemon's detail from api
	if err != nil {
		return res, err
	}

	var moveData []types.Moves

	var errors []error
	//start using goroutine to make the process faster
	var wg sync.WaitGroup
	wg.Add(len(pokemonData.Moves)) //add workers up to the amount of pokemon's move data
	for _, move := range pokemonData.Moves {
		go func(wg *sync.WaitGroup, name string) {
			m, err := pokeapi.Move(name) //get move detail from api
			if err != nil {
				mu.Lock() //lock process to mapping data error
				errors = append(errors, err)
				mu.Unlock() //unlock process
			}

			if m.Accuracy != 0 || m.Power != 0 {
				mu.Lock() //lock process to mapping data move
				//take move name and move power
				moveData = append(moveData, types.Moves{
					Name:  m.Name,
					Power: m.Power,
				})
				mu.Unlock() //unlock process
			}
			wg.Done()
		}(&wg, move.Move.Name)
	}
	wg.Wait()

	if len(errors) != 0 {
		return res, errors[0]
	}

	var hp int
	//take pokemon's HP
	for _, v := range pokemonData.Stats {
		if v.Stat.Name == "hp" {
			hp = v.BaseStat
		}
	}

	res = types.Pokemon{
		Name:  pokemon,
		Hp:    hp,
		Moves: moveData,
	}

	return res, nil
}

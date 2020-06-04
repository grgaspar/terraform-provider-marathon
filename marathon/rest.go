package marathon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*
func main() {

fmt.Println(FindPoke(15))
}
*/
func call() []byte {
	response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseData
	//fmt.Println(string(responseData))
}

// Given an int, will go out and find a pokemon and return the Pokemon name
func FindPoke(i int) string {

	var responseData []byte
	responseData = call()

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	//fmt.Println(responseObject.Name)
	//fmt.Println(len(responseObject.Pokemon))

	pokeName := responseObject.Pokemon[i].Species.Name
	//for i := 0; i < len(responseObject.Pokemon); i++ {
	//  fmt.Println(pokeName)
	//}s.Name
	return pokeName
}

// A Response struct to map the Entire Response
type Response struct {
	Name    string    `json:"name"`
	Pokemon []Pokemon `json:"pokemon_entries"`
}

// A Pokemon Struct to map every pokemon to.
type Pokemon struct {
	EntryNo int            `json:"entry_number"`
	Species PokemonSpecies `json:"pokemon_species"`
}

// A struct to map our Pokemon's Species which includes it's name
type PokemonSpecies struct {
	Name string `json:"name"`
}

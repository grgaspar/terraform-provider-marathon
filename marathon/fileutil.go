package marathon

//package main

import (
	"fmt"
	"os"

	"github.com/magiconair/properties"
	// "io"
	// "bufio"
)

/*
func main() {

var pokemon string
pokemon = "fred"

directory := "foo"

Write(directory + "/" + "configuration.properties", "pokemon", pokemon)
//Show(directory + "/" + "configuration.properties")

var configValue string

configValue = Read(directory + "/" + "configuration.properties", "pokemon")
fmt.Println(configValue)
}

*/

// Read a property from a properties file
func Read(filename string, key string) string {
	var props *properties.Properties

	var loader properties.Loader

	props, err1 := loader.LoadFile(filename)

	if err1 != nil {
		return ""
	}

	ret, ok1 := props.Get(key)

	fmt.Println(ok1)

	return ret
}

// Display all key value pairs in a properties file.
func Show(filename string) {
	var props *properties.Properties

	var loader properties.Loader

	props, err1 := loader.LoadFile(filename)

	if err1 != nil {
		return
	}

	l := props.Len()

	fmt.Println("Length = ", l)

	keys := props.Keys()

	for i, key := range keys {
		ret, _ := props.Get(key)
		fmt.Println(i, ret)
	}

	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Writes a property (key/value pair) to the given filename
func Write(filename string, key string, value string) string {
	var props *properties.Properties

	var loader properties.Loader

	props, err1 := loader.LoadFile(filename)

	if err1 != nil {
		return ""
	}

	ret, ok1 := props.Get(key)

	fmt.Println(ok1)
	fmt.Println(ret)

	props.Set(key, value)

	f, err := os.Create(filename)
	check(err)

	defer f.Close()

	keys := props.Keys()

	for i, key := range keys {
		fmt.Println(i)
		ret, _ := props.Get(key)
		_, err := f.WriteString(key + "=" + ret + "\n")
		check(err)
	}

	f.Sync()

	return ret
}

func printOut(properties *properties.Properties) {

	l := properties.Len()

	fmt.Println("Length = ", l)

	keys := properties.Keys()

	for i, key := range keys {
		ret, _ := properties.Get(key)
		fmt.Println(i, ret)
	}
}

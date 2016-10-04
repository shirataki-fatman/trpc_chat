package main

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/robertkrimen/otto"
)

type PluginProperty struct {
	Message, Name string
	Result        string
}

func (pp *PluginProperty) RollDice(diceNum, diceMin, diceMax int) []string {
	return RollDice(diceNum, diceMin, diceMax)
}

var pluginManager PluginManager

type PluginManager struct {
	plugins map[string]string
}

func (pm *PluginManager) List() []string {
	var result []string
	for k := range pm.plugins {
		result = append(result, k)
	}
	return result
}

func (pm *PluginManager) Load() {
	reg := regexp.MustCompile("(.+)\\.js$")

	files, err := ioutil.ReadDir("./plugin")
	if err != nil {
		panic(err)
	}

	pm.plugins = map[string]string{}
	for f := range files {
		if reg.MatchString(files[f].Name()) {
			submatch := reg.FindSubmatch([]byte(files[f].Name()))
			pm.plugins[string(submatch[1])] = string(submatch[0])
		}
	}

	fmt.Println(pm.plugins)
}

func (pm *PluginManager) Exec(args *ChatArgs) string {
	fileName, ok := pm.plugins[args.Plugin]
	if !ok {
		return ""
	}

	prop := PluginProperty{
		Message: args.Message,
		Name:    args.Name,
		Result:  "",
	}

	vm := otto.New()
	script, err := vm.Compile("./plugin/"+fileName, nil)
	if err != nil {
		panic(err)
	}
	vm.Set("prop", &prop)
	vm.Run(script)

	return prop.Result
}

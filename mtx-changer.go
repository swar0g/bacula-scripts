package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"bytes"
	"time"
	"regexp"
)

func pipe_commands(commands ...*exec.Cmd) ([]byte, error) {
	for i, command := range commands[:len(commands) - 1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		command.Start()
		commands[i + 1].Stdin = out
	}
	final, err := commands[len(commands) - 1].Output()
	if err != nil {
		return nil, err
	}
	return final, nil
}


func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 4 {
		fmt.Println("Too few arguments")
		return
	}
	ctl := args[0]
	cmd := args[1]
	slot := args[2]
	device := args[3]
	drive := args[4]

	var out bytes.Buffer

	switch cmd {
		case "listall":
			cmd := exec.Command("mtx", "-f", ctl, "status")
			cmd.Stdout = &out

			cmd.Run()

			oDTFullRE := regexp.MustCompile(`Data Transfer Element (\d+):Full \(Storage Element (\d+) Loaded\)(:VolumeTag =\s*(.+))`)
			aMatchesDTFull := oDTFullRE.FindAllStringSubmatch(out.String(), -1)

			for _, v := range aMatchesDTFull{
				fmt.Println("D:" + v[1] + ":F:" + v[2] + ":" + v[4])
			}

			oDTEmptyRE := regexp.MustCompile(`Data Transfer Element (\d+):Empty`)
			aMatchesDTEmpty := oDTEmptyRE.FindAllStringSubmatch(out.String(), -1)

			for _, v := range aMatchesDTEmpty{
				fmt.Println("D:" + v[1] + ":E")
			}

			oSEEmptyRE := regexp.MustCompile(`Storage Element (\d+):Empty`)
			aMatchesSEEmpty := oSEEmptyRE.FindAllStringSubmatch(out.String(), -1)

			for _, v := range aMatchesSEEmpty{
				fmt.Println("S:" + v[1] + ":E")
			}

			oSEFullRE := regexp.MustCompile(`Storage Element (\d+):Full( :VolumeTag=(.+))`)
			aMatchesSEFull := oSEFullRE.FindAllStringSubmatch(out.String(), -1)

			for _, v := range aMatchesSEFull{
				fmt.Println("S:" + v[1] + ":F:" + v[3])
			}

			oSEImExERE := regexp.MustCompile(`Storage Element (\d+) IMPORT.EXPORT:Empty`)
			aMatchesSEImExE := oSEImExERE.FindAllStringSubmatch(out.String(), -1)

			for _, v := range aMatchesSEImExE{
				fmt.Println("I:" + v[1] + ":E")
			}

			oSEImExFRE := regexp.MustCompile(`Storage Element (\d+) IMPORT.EXPORT:Full( :VolumeTag=(.+))`)
			aMatchesSEImExF := oSEImExFRE.FindAllStringSubmatch(out.String(), -1)

			for _, v := range aMatchesSEImExF{
				fmt.Println("I:" + v[1] + ":F:" + v[3])
			}

			os.Exit(0)
		case "unload":
			cmd := exec.Command("mtx", "-f", ctl, "unload", slot, drive)
			cmd.Run()

			time.Sleep(30 * time.Second)
			os.Exit(0)
		case "load":
			cmd := exec.Command("mtx", "-f", ctl, "load", slot, drive)
			cmd.Run()
			time.Sleep(30 * time.Second)
			os.Exit(0)
		case "list":
			cmd := exec.Command("mtx", "-f", ctl, "status")
			cmd.Stdout = &out

			cmd.Run()

			oSElemRE := regexp.MustCompile(`Storage Element ([0-9]*):.*Full\s*:VolumeTag\=(.*)`)
			oDTElemRE := regexp.MustCompile(`Data Transfer Element [0-9]*:.*Full\s*\(Storage Element ([0-9]*) Loaded\)\:VolumeTag \= (.*)`)

			aMatchesSE := oSElemRE.FindAllStringSubmatch(out.String(), -1)
			aMatchesDT := oDTElemRE.FindAllStringSubmatch(out.String(), -1)

			for _, v := range aMatchesSE{
				fmt.Println(v[1] + ":" + v[2])
			}

			for _, v := range aMatchesDT{
				fmt.Println(v[1] + ":" + v[2])
			}
			os.Exit(0)
		case "transfer":
			cmd := exec.Command("mtx", "-f", ctl, "transfer", slot, device)
			cmd.Stdout = &out

			cmd.Run()

			time.Sleep(60 * time.Second)
			os.Exit(0)
		case "slots":
			cmd := exec.Command("mtx", "-f", ctl, "status")
			cmd.Stdout = &out

			cmd.Run()

			oSlotsRE 	:= regexp.MustCompile(`Storage Changer ` + ctl + `:[0-9]* Drives, ([0-9]*) Slots`)

			aMatchesSlots := oSlotsRE.FindAllStringSubmatch(out.String(), -1)

			for _, v := range aMatchesSlots{
				fmt.Println(v[1])
			}

			os.Exit(0)
		case "loaded":
			cmd := exec.Command("mtx", "-f", ctl, "status")
			cmd.Stdout = &out

			cmd.Run()

			oDTFullRE 	:= regexp.MustCompile(`Data Transfer Element ` + drive + `:.*Full\s*\(Storage Element ([0-9]*)`)
			oDTEmptyRE 	:= regexp.MustCompile(`Data Transfer Element ` + drive + `:.*Empty\s*\(Storage Element ([0-9]*)`)

			aMatchesDTFull := oDTFullRE.FindAllStringSubmatch(out.String(), -1)
			aMatchesDTEmpty := oDTEmptyRE.FindAllStringSubmatch(out.String(), -1)

			for _, v := range aMatchesDTFull{
				fmt.Println(v[1])
			}

			for i, _ := range aMatchesDTEmpty{
				if i >= 0 {
					fmt.Println(0)
				}
			}
			os.Exit(0)
	}
}


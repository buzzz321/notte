package main

import (
	"fmt"
	"strings"
	"time"
)

type Cmds struct {
	selector string
	cmd      string
	help     string
	//Call     func(p string, cmds []Cmds) (bool, string)
	Call func(p string, pp string) (bool, string)
}

func showAll(_ string, _ string) (bool, string) {
	var retval string
	for _, item := range commands {
		retval += item.help + "\n"
	}
	return true, retval
}

func whatIsClock(_ string, _ string) (bool, string) {
	datum := time.Now()
	fmt.Println(datum.Format("02 Jan 06 15:04 MST"))
	return true, datum.Format("02 Jan 06 15:04 MST")
}

func whatIsCrafts(username string, line string) (bool, string) {
	tokens := strings.Split(line, " ")

	if len(tokens) == 0 {
		crafts := showCrafting("", username)
		return true, crafts[:]
	}
	if len(tokens) >= 1 {
		var crafts string
		var craftsArr []string

		if strings.TrimSpace(tokens[0]) == "all" {
			chars := strings.Split(showChars(username), "\n")
			for _, char := range chars[:len(chars)-1] {
				craftsArr = append(craftsArr, "**"+char+"**")
				craftsArr = append(craftsArr, showCrafting(char, username))
			}
			crafts = strings.Join(craftsArr, "\n")
		} else {
			crafts = showCrafting(strings.Join(tokens[1:], " "), username)
		}
		return true, crafts[:]
	}
	return false, ""
}

/*
func findMyItem(username string, cmdline string) (bool, string) {
	items := gw2util.GetItems(gw2, GetItemIdsFromBags(jsonParsed, username))
	return gw2util.FindItem(items, "food")
}
*/
func whatIsChars(username string, _ string) (bool, string) {
	return true, showChars(username)
}

func searchInBags(username string, itemName string) (bool, string) {
	return true, findItem(username, itemName)
}

func removeCmd(cmd Cmds, cmdline string) string {
	retVal := strings.Replace(cmdline, cmd.selector, "", 1)
	return strings.TrimSpace(strings.Replace(retVal, cmd.cmd, "", 1))
}

func Process(username string, cmdline string) (bool, string) {

	for _, cmd := range commands {
		if strings.Contains(cmdline, cmd.selector) && strings.Contains(cmdline, cmd.cmd) {
			return cmd.Call(username, removeCmd(cmd, cmdline))
		}
	}

	return false, ""
}

func initCommands() []Cmds {
	var cmds = []Cmds{
		Cmds{
			selector: "what is",
			cmd:      "time",
			help:     "<what is time> Shows the time and date",
			Call:     whatIsClock,
		},
		Cmds{
			selector: "show my",
			cmd:      "crafts",
			help:     "<show my crafts> Show all crafts for a char",
			Call:     whatIsCrafts,
		},
		Cmds{
			selector: "show my",
			cmd:      "commands",
			help:     "<show my commands> Show all commands",
			Call:     showAll,
		},
		Cmds{
			selector: "show my",
			cmd:      "chars",
			help:     "<show my chars> Show all chars connected to api key",
			Call:     whatIsChars,
		},
		Cmds{
			selector: "search in",
			cmd:      "bags",
			help:     "<search in bags> Search thru all bags for an item",
			Call:     searchInBags,
		},
	}

	return cmds
}

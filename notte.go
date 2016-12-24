package main

import (
	"fmt"
	"bufio"
	"os"
	"log"

	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
	"gw2util"
)


// Variables used for command line parameters
var (
	BotID string
)

func readkey(filename string) (string) {
	inputFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Scan()
	key := scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return key
}

func whatIs(username string, line string) (bool, string) {
	if (strings.Contains(strings.ToLower(line), "vad är") ||
		strings.Contains(strings.ToLower(line), "show my")) {
		tokens := strings.Split(line, " ")
		datum := time.Now()
		switch tokens[2]{
		case "klockan":
			fmt.Println(datum.Format("02 Jan 06 15:04 MST"))
			return true, datum.Format("02 Jan 06 15:04 MST")
		case "crafts":
			crafts := showCrafting(username)

			return true, crafts[:]
		}
	}

	return false, ""
}
/*

	gw2 := gw2util.Gw2Api{"https://api.guildwars2.com/v2/", gw2util.GetKey("../../../gw2/test.key")}
	body := gw2util.QueryAnetAuth(gw2, "characters")

 */
func showCrafting(name string) string {
	var retVal[] string

	gw2 := gw2util.Gw2Api{"https://api.guildwars2.com/v2/", gw2util.GetKey("../../../gw2/test.key")}
	crafts := gw2util.GetCrafting(gw2, name)
	for _, craft:= range crafts{
		tmp := craft.String()
		retVal = append(retVal, tmp)
	}

	return strings.Join(retVal,"")
}

// This function will be called (due to AddHandler below) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	whatis, answer := whatIs(m.Author.Username, strings.TrimSpace(m.Content))
	if (whatis) {
		_, _ = s.ChannelMessageSend(m.ChannelID, answer)
	}
}

func main() {
	discordKey := readkey("../../../discord/disc.key")
	//fmt.Printf("Hello, 世界 key=[%s]\n", key)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + discordKey)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	/*	channels, err := dg.UserChannels()

		if err != nil {
			fmt.Println("error obtaining user channels details,", err)
		} else {
			for i := 0; i < len(channels); i++ {
				fmt.Println("channel: ", channels[i].Name)
			}
		}

		guilds, err := dg.UserGuilds()
		if err != nil {
			fmt.Println("error obtaining user guilds details,", err)
		} else {
			for i := 0; i < len(guilds); i++ {
				fmt.Println("Guild name:", guilds[i].Name)
				fmt.Println("channals: ", len(guilds[0].Channels))
				for _, channel := range guilds[i].Channels {
					fmt.Println("channel: ", channel.Name)
				}
			}
		}*/
	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})

	fmt.Println("Exiting...")
	return

}

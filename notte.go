package main

import (
    "bufio"
    "fmt"
    "log"
    "os"

    "gw2util"
    "strings"

    "github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
    BotID    string
    userData gw2util.UserDataSlice
    commands []Cmds
)

func readkey(filename string) string {
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

func chunkString(message string, chunkSize int) []string {
    var retVal []string
    chunks := len(message) / chunkSize
    remainder := len(message) % chunkSize

    for i := 0; i < chunks; i++ {
        retVal = append(retVal, message[i*chunkSize:chunkSize+i*chunkSize])
    }
    retVal = append(retVal, message[chunks*chunkSize:remainder+chunks*chunkSize])
    return retVal
}

func setMy(username string, line string) (bool, string) {
    if strings.Contains(strings.ToLower(line), "spara min") ||
        strings.Contains(strings.ToLower(line), "set my") {
        tokens := strings.Split(line, " ")
        fmt.Println(tokens)
        if len(tokens) < 3 {
            return false, ""
        }

        switch tokens[2] {
        case "apikey":
            gameID := strings.Join(tokens[3:len(tokens)-1], " ")
            userData = gw2util.UpsertUserData(userData, gw2util.UserData{ChatName: username, GameID: gameID, Key: tokens[len(tokens)-1]})
            fmt.Printf("username = %s GameId = %s Key = %s\n", username, gameID, tokens[len(tokens)-1])
            gw2util.SaveUserData(userData)
        }
    }
    return false, ""
}

/*

   gw2 := gw2util.Gw2Api{"https://api.guildwars2.com/v2/", gw2util.GetKey("../../../gw2/test.key")}
   body := gw2util.QueryAnetAuth(gw2, "characters")

*/

func findItem(chatName string, itemName string) string {
    var retVal string

    userData := gw2util.GetUserData(userData, chatName)
    gw2 := gw2util.Gw2Api{BaseURL: "https://api.guildwars2.com/v2/", Key: userData.Key}

    items := gw2util.FindItem(gw2, userData.GameID, itemName)

    for _, item := range items {
        retVal += item.String() + "\n"
    }
    return retVal
}

func showCrafting(charName string, chatName string) string {
    var retVal []string
    
    userData := gw2util.GetUserData(userData, chatName)

    gw2 := gw2util.Gw2Api{BaseURL: "https://api.guildwars2.com/v2/", Key: userData.Key}
    if charName == "" {
        charName = userData.GameID
    }
    crafts := gw2util.GetCrafting(gw2, charName)
    for _, craft := range crafts {
        tmp := craft.String()
        retVal = append(retVal, tmp)
    }

    return strings.Join(retVal, "")
}

func showChars(chatName string) string {
    var retVal []string

    fmt.Println(chatName)
    gw2 := gw2util.Gw2Api{BaseURL: "https://api.guildwars2.com/v2/", Key: gw2util.GetUserData(userData, chatName).Key}
    fmt.Println(gw2)
    chars := gw2util.GetCharacterNames(gw2)
    fmt.Println(chars)
    for _, char := range chars {
        if char != "" {
            retVal = append(retVal, char)
            retVal = append(retVal, "\n")
        }
    }

    return strings.Join(retVal, "")
}

func getWvWvWKD(chatName string) string {
    var kdRed, kdGreen, kDBlue float64 = 0.0, 0.0, 0.0

    gw2 := gw2util.Gw2Api{BaseURL: "https://api.guildwars2.com/v2/", Key: gw2util.GetUserData(userData, chatName).Key}
    stats := gw2util.GetWWWStats(gw2, "2007")

    if stats.Kills.Blue > 0 {
        kDBlue = stats.Kills.Blue / stats.Deaths.Blue
    }
    if stats.Kills.Red > 0 {
        kdRed = stats.Kills.Red / stats.Deaths.Red
    }
    if stats.Kills.Green > 0 {
        kdGreen = stats.Kills.Green / stats.Deaths.Green
    }

    return fmt.Sprintf("K/D \n Blue: %6.2f\n Red: %6.2f\n Green: %6.2f\n", kDBlue, kdRed, kdGreen)
}

func showWvWvWstats(chatName string) string {
    gw2 := gw2util.Gw2Api{BaseURL: "https://api.guildwars2.com/v2/", Key: gw2util.GetUserData(userData, chatName).Key}
    stats := gw2util.GetWWWStats(gw2, "2007")

    return stats.String()
}

/*
func chunkMessage(message string, chunksize uint64) []string {
    var retVal []string
}
*/
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

    //    whatis, answer := whatIs(m.Author.Username, strings.TrimSpace(m.Content))
    whatis, answer := Process(m.Author.Username, strings.TrimSpace(m.Content))

    if whatis {
        if len(answer) < 2000 {
            _, _ = s.ChannelMessageSend(m.ChannelID, answer)
        } else {
            chunks := chunkString(answer, 2000)

            for _, chunk := range chunks {
                _, _ = s.ChannelMessageSend(m.ChannelID, chunk)
            }
        }

    }

    setmy, answer := setMy(m.Author.Username, strings.TrimSpace(m.Content))
    if setmy {
        _, _ = s.ChannelMessageSend(m.ChannelID, answer)
    }

}

func main() {
    discordKey := readkey("../../../discord/disc.key")
    // Create a new Discord session using the provided bot token.
    dg, err := discordgo.New("Bot " + discordKey)
    if err != nil {
        fmt.Println("Error creating Discord session: ", err)
        return
    }

    commands = initCommands()
    userData = gw2util.ReadUserData("data.dat")
    // Get the account information.
    u, err := dg.User("@me")
    if err != nil {
        fmt.Println("error obtaining account details,", err)
    }
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

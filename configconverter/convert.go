package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/goccy/go-yaml"
)

func main() {
	setupConfig()

	if err := createIfDoesntExist("new-configs/"); err != nil {
		log.Fatalf("error creating new folder")
	}

	// for _, oldCommand := range getCommands() {
	// 	oldCommand = strings.TrimPrefix(oldCommand, "command.")
	// 	oldCommand = strings.TrimSuffix(oldCommand, ".response")
	// 	fmt.Printf(oldCommand + "\n")
	// }

	// for _, oldKeyword := range getKeywords() {
	// 	oldKeyword = strings.TrimPrefix(oldKeyword, "keyword.")
	// 	oldKeyword = strings.TrimSuffix(oldKeyword, ".response")
	// 	oldKeyword = strings.TrimSuffix(oldKeyword, ".reaction")
	// 	fmt.Printf(oldKeyword + "\n")
	// }

	// fmt.Printf("%s\n", getParsingImageFiletypes())

	// fmt.Printf("%s\n", getParsingPasteKeys())

	// for _, parseKey := range strings.Split(getParsingPasteKeys(), ", ") {
	// 	parseKey = strings.TrimSuffix(parseKey, ".url")
	// 	parseKey = strings.TrimSuffix(parseKey, ".format")
	// 	parseKey = strings.TrimSuffix(parseKey, ".append")
	// 	if parseKey == "parse.image.filetype" {

	// 	} else {
	// 		newParse := parsingConfig{
	// 			Name:   parseKey,
	// 			URL:    getParsingPasteString(parseKey + ".url"),
	// 			Format: getParsingPasteString(parseKey + ".format"),
	// 		}
	// 		fmt.Printf("%+v\n", newParse)
	// 	}
	// }

	for _, service := range getBotServices() {
		switch service {
		case "discord":
			createDiscordConf()
		}
	}
}

// Exists reports whether the named file or directory exists.
func createIfDoesntExist(name string) (err error) {
	path, file := path.Split(name)

	// if confdir exists carry on
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			if _, err = os.Stat(name); err != nil {
				if file == "" {
					if err = os.Mkdir(path, 0755); err != nil {
					}
				} else {
					if fileCheck, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644); err != nil {
					} else {
						fileCheck.Close()
					}
				}
			}
		}
	}
	return
}

// created configs for discord bot.
func createDiscordConf() {
	newDiscordBot := discordBot{}

	newDiscordBot.Config.Token = getDiscordConfigString("token")
	newDiscordBot.Config.Game = getDiscordConfigString("game")

	fmt.Printf("%+v\n", newDiscordBot)

	newDiscordServer := discordServer{}

	// set up discord config
	fmt.Printf("loading bot configs\n")
	newDiscordServer.ServerID = "set and ID here"
	newDiscordServer.Config.Prefix = getDiscordConfigString("prefix")
	newDiscordServer.Config.Clear = getDiscordConfigBool("command.remove")
	newDiscordServer.Config.Webhooks.Logs = getDiscordConfigString("webhook")

	// set up channel groups
	fmt.Printf("loading channel groups\n")
	newChanGroup := channelGroup{}

	for _, channel := range Discord.GetStringSlice("discord.channels.listening") {
		newChanGroup.ChannelIDs = append(newChanGroup.ChannelIDs, channel)
	}

	// set up commands
	fmt.Printf("loading commands\n")
	newCommands := []command{}

	fmt.Printf("There are %d commands\n", len(getCommands()))

	for _, oldCommand := range getCommands() {
		oldCommand = strings.TrimPrefix(oldCommand, "command.")
		oldCommand = strings.TrimSuffix(oldCommand, ".response")
		newCommands = append(newCommands, command{
			Command:  fmt.Sprintf("%s", oldCommand),
			Reaction: getCommandReaction(oldCommand),
			Response: getCommandResonse(oldCommand),
		})
	}

	newChanGroup.Commands = newCommands

	// setup keywords
	fmt.Printf("loading keywords\n")
	newKeywords := []keyword{}

	for _, oldKeyword := range getKeywords() {
		oldKeyword = strings.TrimPrefix(oldKeyword, "keyword.")
		oldKeyword = strings.TrimSuffix(oldKeyword, ".response")
		oldKeyword = strings.TrimSuffix(oldKeyword, ".reaction")
		if strings.HasPrefix(oldKeyword, "exact.") {
			fmt.Printf("found exact prefix for %s\n", oldKeyword)
			newKeywords = append(newKeywords, keyword{
				Keyword:  fmt.Sprintf("\"%s\"", strings.TrimPrefix(oldKeyword, "exact.")),
				Reaction: getKeywordReaction(oldKeyword),
				Response: getKeywordResponse(oldKeyword),
				Exact:    true,
			})
		} else {
			newKeywords = append(newKeywords, keyword{
				Keyword:  fmt.Sprintf("\"%s\"", oldKeyword),
				Reaction: getKeywordReaction(oldKeyword),
				Response: getKeywordResponse(oldKeyword),
			})
		}

	}

	newChanGroup.Keywords = newKeywords

	newParsing := parsing{}
	fmt.Printf("loading parsing\n")
	fmt.Printf("%s\n", getParsingPasteKeys())
	newParsing.Image.Filetypes = getParsingImageFiletypes()
	for _, parseKey := range strings.Split(getParsingPasteKeys(), ", ") {
		// lot's of dumb stuff from viper
		if strings.HasSuffix(parseKey, ".format") || strings.HasSuffix(parseKey, ".append") || parseKey == "parse.image.filetype" {
			continue
		}

		parseKey = strings.TrimSuffix(parseKey, ".url")

		newParse := parsingConfig{
			Name:   parseKey,
			URL:    getParsingPasteString(parseKey + ".url"),
			Format: getParsingPasteString(parseKey + ".format"),
		}
		fmt.Printf("checking on %s\n", newParse.Name)
		if len(newParsing.Paste.Sites) == 0 {
			fmt.Printf("no parsing sites loaded\n")
		} else {
			for _, site := range newParsing.Paste.Sites {
				// fmt.Printf("found config for %s\n", site.Name)
				if site.Name == newParse.Name {
					fmt.Printf("A config for %s already exists\n", newParse.Name)
					continue
				}
			}
		}

		fmt.Printf("adding a config for %s\n", newParse.Name)
		newParsing.Paste.Sites = append(newParsing.Paste.Sites, newParse)

		fmt.Printf("finished checking on %s\n", newParse.Name)
	}

	newChanGroup.Parsing = newParsing

	// set up permissions
	fmt.Printf("loading permissions\n")
	newPerms := []permission{}
	for _, group := range getDiscordGroups() {
		newPerms = append(newPerms, permission{
			Group: group,
			Users: getDiscordGroupUsers(group),
			Roles: getDiscordGroupRoles(group),
		})
	}

	newDiscordServer.ChanGroups = append(newDiscordServer.ChanGroups, newChanGroup)

	newDiscordServer.Permissions = newPerms

	newFileFolders := []string{"new-configs/discord/", "new-configs/discord/server/", "new-configs/discord/bot.yml", "new-configs/discord/server/server.yml"}

	for _, folder := range newFileFolders {
		if err := createIfDoesntExist(folder); err != nil {
			log.Fatalf("error creating new folder")
		}
	}

	if err := writeYamlToFile("new-configs/discord/bot.yml", newDiscordBot); err != nil {
		log.Fatalf("error writing to new config file")
	}

	if err := writeYamlToFile("new-configs/discord/server/server.yml", newDiscordServer); err != nil {
		log.Fatalf("error writing to new config file")
	}

}

func writeYamlToFile(file string, iface interface{}) (err error) {
	ydata, err := yaml.Marshal(iface)
	if err != nil {
		return
	}

	// create a file with a supplied name
	yamlFile, err := os.Create(file)
	if err != nil {
		return
	}

	if _, err = yamlFile.Write(ydata); err != nil {
		return
	}

	return
}

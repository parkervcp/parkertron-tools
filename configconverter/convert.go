package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
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

		var newCommandReaction []string
		for _, line := range getCommandReaction(oldCommand) {
			newCommandReaction = append(newCommandReaction, fmt.Sprintf("'%s'", line))
		}

		var newCommandResponse []string
		for _, line := range getCommandResonse(oldCommand) {
			newCommandResponse = append(newCommandResponse, fmt.Sprintf("'%s'", line))
		}

		newCommands = append(newCommands, command{
			Command:  fmt.Sprintf("%s", oldCommand),
			Reaction: newCommandReaction,
			Response: newCommandResponse,
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
			// add ' marks around old reaction strings
			var newKeywordReaction []string
			for _, line := range getKeywordReaction(oldKeyword) {
				newKeywordReaction = append(newKeywordReaction, fmt.Sprintf("'%s'", line))
			}
			// add ' marks around old response strings
			var newKeywordResponse []string
			for _, line := range getKeywordResponse(oldKeyword) {
				newKeywordResponse = append(newKeywordResponse, fmt.Sprintf("'%s'", line))
			}
			newKeywords = append(newKeywords, keyword{
				Keyword:  fmt.Sprintf("'%s'", strings.TrimPrefix(oldKeyword, "exact.")),
				Reaction: newKeywordReaction,
				Response: newKeywordResponse,
				Exact:    true,
			})
		} else {
			var newKeywordReaction []string
			// add ' marks around old reaction strings
			for _, line := range getKeywordReaction(oldKeyword) {
				newKeywordReaction = append(newKeywordReaction, fmt.Sprintf("'%s'", line))
			}
			// add ' marks around old response strings
			var newKeywordResponse []string
			for _, line := range getKeywordResponse(oldKeyword) {
				newKeywordResponse = append(newKeywordResponse, fmt.Sprintf("'%s'", line))
			}
			newKeywords = append(newKeywords, keyword{
				Keyword:  fmt.Sprintf("'%s'", oldKeyword),
				Reaction: newKeywordReaction,
				Response: newKeywordResponse,
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
			Name:   fmt.Sprintf("'%s'", parseKey),
			URL:    fmt.Sprintf("'%s'", getParsingPasteString(parseKey+".url")),
			Format: fmt.Sprintf("'%s'", getParsingPasteString(parseKey+".format")),
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

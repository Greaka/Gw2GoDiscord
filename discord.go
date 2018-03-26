package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/greaka/gw2godiscord/commands"
	"github.com/greaka/gw2godiscord/language"
)

var (
	Prefix string = ";"
)

func InitializeDiscord(token string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageReceived)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection, ", err)
		return
	}
	defer dg.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageReceived(session *discordgo.Session, message *discordgo.MessageCreate) {
	go handleMessage(session, message)
}

func handleMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	if message.Content == "ping" {
		if _, err := session.ChannelMessageSend(message.ChannelID, "pong!"); err != nil {
			fmt.Println("%v", err)
		}
		return
	}

	if !strings.HasPrefix(message.Content, Prefix) {
		checkForDM(session, message)
		return
	}
	stripped := strings.TrimPrefix(message.Content, Prefix)
	splitted := strings.Fields(stripped)
	if len(splitted) < 1 {
		return
	}

	go session.ChannelTyping(message.ChannelID)

	switch splitted[0] {
	case commands.CommandKey:
		commands.HandleKey(splitted[1:])
	default:
		checkForDM(session, message)
		return
	}
}

func checkForDM(session *discordgo.Session, message *discordgo.MessageCreate) {
	channel, err := session.Channel(message.ChannelID)
	if err != nil {
		return
	}
	switch channel.Type {
	case discordgo.ChannelTypeDM:
		session.ChannelMessageSend(message.ChannelID, language.NotACommand(language.English))
	case discordgo.ChannelTypeGroupDM:
		session.ChannelMessageSend(message.ChannelID, language.NotACommand(language.English))
	}
}

func sendMessage(session *discordgo.Session, channelId, content string) {
	perm, e := session.UserChannelPermissions(session.State.User.ID, channelId)
	if e != nil {
		Log(e)
		return
	}
	if perm&discordgo.PermissionSendMessages != 0 {
		if _, err := session.ChannelMessageSend(channelId, content); err != nil {
			Log(err)
		}
	}
}

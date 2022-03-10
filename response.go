package slashlib

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// 返すタイプ
type ReturnType int8

const (
	// Reply for Message
	ReplyMessage ReturnType = 4
	// Reply for Thiking
	// pls leter Edit()
	ReplyThiking ReturnType = 5
	// Interaction In Modal Window
	Window ReturnType = 9
)

// チェーンメゾット用の型
type InteractionResponse struct {
	Discord     *discordgo.Session
	Interaction *discordgo.Interaction
	Response    *discordgo.InteractionResponse
}

// 返すタイプ
type CommandType int8

const (
	// Slash Command (MessageCommand)
	SlashCommand CommandType = 1
	// Comporment Command (MessageComponent)
	ComponentCommand CommandType = 2
	// Submit (ModalSubmit)
	SubmitCommand CommandType = 3
)

// 整形用構造体
type InteractionStruct struct {
	GuildID     string
	GuildName   string
	GuildData   *discordgo.Guild
	ChannelID   string
	ChannelName string
	ChannelData *discordgo.Channel
	UserID      string
	UserNum     string
	UserName    string
	UserData    *discordgo.User
	// TypeResult:
	// discordgo.InteractionApplicationCommand
	// discordgo.InteractionMessageComponent
	// discordgo.InteractionModalSubmit
	Type      discordgo.InteractionType
	Check     CommandType
	Command   discordgo.ApplicationCommandInteractionData
	Component discordgo.MessageComponentInteractionData
	Submit    discordgo.ModalSubmitInteractionData
}

// InteractionResponseのFlag
const Invisible uint64 = 1 << 6

// InteractionCreate 整形
func InteractionViewAndEdit(discord *discordgo.Session, i *discordgo.InteractionCreate) (iData InteractionStruct) {
	var err error
	cmdData := i.Interaction
	iData.GuildID = cmdData.GuildID
	iData.GuildData, err = discord.Guild(iData.GuildID)
	if err == nil {
		iData.GuildName = iData.GuildData.Name
	} else {
		iData.GuildName = "DirectMessage"
	}
	iData.ChannelID = cmdData.ChannelID
	iData.ChannelData, _ = discord.Channel(iData.ChannelID)
	iData.ChannelName = iData.ChannelData.Name
	// DMならばUser じゃ無ければMember
	if cmdData.User != nil {
		iData.UserNum = cmdData.User.Discriminator
		iData.UserName = cmdData.User.Username
		iData.UserData = cmdData.User
	} else {
		iData.UserNum = cmdData.Member.User.Discriminator
		iData.UserName = cmdData.Member.User.Username
		iData.UserData = cmdData.Member.User
	}
	iData.Type = cmdData.Type
	switch iData.Type {
	case discordgo.InteractionApplicationCommand:
		iData.Check = SlashCommand
		iData.Command = cmdData.ApplicationCommandData()
		//表示
		log.Print("Guild:\"" + iData.GuildName + "\"  Channel:\"" + iData.ChannelName + "\"  [" + iData.UserName + "#" + iData.UserNum + "] Slash /" + iData.Command.Name)
	case discordgo.InteractionMessageComponent:
		iData.Check = ComponentCommand
		iData.Component = cmdData.MessageComponentData()
		//表示
		log.Print("Guild:\"" + iData.GuildName + "\"  Channel:\"" + iData.ChannelName + "\"  [" + iData.UserName + "#" + iData.UserNum + "] Component ID:" + iData.Component.CustomID)
	case discordgo.InteractionModalSubmit:
		iData.Check = SubmitCommand
		iData.Submit = cmdData.ModalSubmitData()
		//表示
		log.Print("Guild:\"" + iData.GuildName + "\"  Channel:\"" + iData.ChannelName + "\"  [" + iData.UserName + "#" + iData.UserNum + "] Submit ID:" + iData.Submit.CustomID)
	}
	return
}

// Interaction Return Message
// Flags Usual: Invisible
func (i *InteractionResponse) Return(resType ReturnType, resData *discordgo.InteractionResponseData) error {
	i.Response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(resType),
		Data: resData,
	}
	return i.Discord.InteractionRespond(i.Interaction, i.Response)
}

// Interaction Window Message
// Component doesn't usual: AddButton(),AddMenu
// Result Type is Submit, doesn't Component,SlashCommand.
func (i *InteractionResponse) Window(title, customID string, comp *Component) error {
	i.Response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(Window),
		Data: &discordgo.InteractionResponseData{
			Title:      title,
			CustomID:   customID,
			Components: comp.Parse(),
		},
	}
	return i.Discord.InteractionRespond(i.Interaction, i.Response)
}

// Interaction Edit Message
// Flags Usual: Invisible
func (i *InteractionResponse) Edit(newData *discordgo.WebhookEdit) error {
	appID := i.Discord.State.User.ID
	_, err := i.Discord.InteractionResponseEdit(appID, i.Interaction, newData)
	return err
}

package slashlib

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// 最後に呼ばれた関数
type lastCallFunc int8

// 関数ナンバー
const (
	// AddCommand
	addCommand lastCallFunc = 0
	// AddOption
	addOption lastCallFunc = 1
	// AddChoice
	addChoice lastCallFunc = 2
)

// Optionタイプ
type OptionType int8

const (
	// String
	TypeString OptionType = 3
	// Int
	TypeInt OptionType = 4
	// Bool
	TypeBool OptionType = 5
	// @User
	TypeUser OptionType = 6
	// #Channnel
	TypeChannel OptionType = 7
	// @Role
	TypeRole OptionType = 8
	// Mentionable (@User,@Role)
	TypeMention OptionType = 9
	// Float
	TypeFloat OptionType = 10
	// Files
	TypeFile OptionType = 11
)

// チェーンメゾット用の型
type Command struct {
	Discord  []*discordgo.ApplicationCommand
	lastCall lastCallFunc
}

// コマンド生成
func (c *Command) AddCommand(cmd, description string) *Command {
	c.Discord = append(c.Discord, &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        cmd,
		Description: description,
	})
	c.lastCall = addCommand
	return c
}

// コマンドのオプションの追加
// valueMin,valueMax request cmdType= int or float
func (c *Command) AddOption(cmdType OptionType, name, description string, request bool, valueMin, valueMax float64) *Command {
	data := &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionType(cmdType),
		Name:        name,
		Description: description,
		Required:    request,
		MinValue:    &valueMin,
		MaxValue:    valueMax,
	}
	c.Discord[c.CommandLast()].Options = append(c.Discord[c.CommandLast()].Options, data)
	c.lastCall = addOption
	return c
}

// コマンドの選択を追加
// max 25 Choices
func (c *Command) AddChoice(name string, value interface{}) *Command {
	if c.lastCall == addCommand {
		panic("AddChoice() Request before AddChoice() or AddOption()")
	}

	switch {
	case c.Discord[c.CommandLast()].Options[c.OptionsLast()].Type == 3: // string
		fallthrough
	case c.Discord[c.CommandLast()].Options[c.OptionsLast()].Type == 4: // int
		fallthrough
	case c.Discord[c.CommandLast()].Options[c.OptionsLast()].Type == 10: // float
		c.Discord[c.CommandLast()].Options[c.OptionsLast()].Choices = append(c.Discord[c.CommandLast()].Options[c.OptionsLast()].Choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: value,
		})
		c.lastCall = addChoice
	default:
		panic("AddChoice() Request OptionType string,int or float")
	}
	return c
}

// Command[-1] 参照用
func (c *Command) CommandLast() int {
	return len(c.Discord) - 1
}

// Command[-1].Options[-1] 参照用
func (c *Command) OptionsLast() int {
	return len(c.Discord[c.CommandLast()].Options) - 1
}

// Command生成
// all guild when guildID == ""
func (c *Command) CommandCreate(discord *discordgo.Session, guildID string) {
	for _, command := range c.Discord {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, guildID, command)
		if err != nil {
			fmt.Printf("Failed Create Command \"%s\"\n", command.Name)
			panic(err)
		}
	}
}

// Command削除
// all guild when guildID == ""
func CommandDelete(discord *discordgo.Session, guildID string, name string) error {
	cmd, err := discord.ApplicationCommands(discord.State.User.ID, guildID)
	if err != nil {
		return err
	}
	for _, cmdData := range cmd {
		if cmdData.Name == name {
			return discord.ApplicationCommandDelete(discord.State.User.ID, guildID, cmdData.ID)
		}
	}
	return fmt.Errorf("not found \"%s\"", name)
}
func FindData(interaction *discordgo.InteractionCreate, key string) (value []interface{}) {
	for _, data := range interaction.ApplicationCommandData().Options {
		if data.Name == key {
			value = append(value, data.Value)
		}
	}
	return
}

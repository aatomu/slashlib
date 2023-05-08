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
// perm is command run user permission by https://github.com/bwmarrin/discordgo/blob/v0.27.1/structs.go#L2092-L2179
func (c *Command) AddCommand(cmd, description string, perm int64) *Command {
	c.Discord = append(c.Discord, &discordgo.ApplicationCommand{
		Type:                     discordgo.ChatApplicationCommand,
		Name:                     cmd,
		DefaultMemberPermissions: &perm,
		Description:              description,
	})
	c.lastCall = addCommand
	return c
}

// コマンドのオプションの追加
// how to use Minvalue,MinLength
// x :=1.1
// ops.MinValue = &x
func (c *Command) AddOption(ops *discordgo.ApplicationCommandOption) *Command {
	c.Discord[c.CommandLast()].Options = append(c.Discord[c.CommandLast()].Options, ops)
	c.lastCall = addOption
	return c
}

// コマンドの選択を追加
// max 25 Choices
func (c *Command) AddChoice(name string, value interface{}) *Command {
	if c.lastCall == addCommand {
		panic("AddChoice() Request before AddChoice() or AddOption()")
	}

	switch c.Discord[c.CommandLast()].Options[c.OptionsLast()].Type {
	case 3, 4, 10: // string,int,float
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

func FindValue(interaction *discordgo.InteractionCreate, key string) (value []interface{}) {
	for _, data := range interaction.ApplicationCommandData().Options {
		if data.Name == key {
			value = append(value, data.Value)
		}
	}
	return
}

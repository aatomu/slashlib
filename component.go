package slashlib

import (
	"github.com/bwmarrin/discordgo"
)

// チェーンメゾット用の型
type Component struct {
	Discord []discordgo.ActionsRow
}

// 行の追加
// pls before run AddButton,AddMenu or AddInput
func (c *Component) AddLine() *Component {
	c.Discord = append(c.Discord, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{},
	})
	return c
}

// ボタンの追加
// Request: Label, Style.
// if Style:5 => Request: URL
// unles Style:5 => Request: "CustomID"
// max 5 buttons in 1line
func (c *Component) AddButton(data discordgo.Button) *Component {
	if len(c.Discord[c.CompArrayLast(1)].Components) > 5 {
		panic("Up to 5 buttons per line")
	}
	c.Discord[c.CompArrayLast(1)].Components = append(c.Discord[c.CompArrayLast(1)].Components, data)
	return c
}

// 選択メニューの追加
// FuncReq: AddLine()
// Request: CustomID, Options.
// SelectMenuOption Request: Label,Value.
// if MinValue or MaxValue != 0 multi select
func (c *Component) AddMenu(data discordgo.SelectMenu) *Component {
	if len(c.Discord[c.CompArrayLast(1)].Components) != 0 {
		panic("AddMenu() Request AddLine() before")
	}
	c.Discord[c.CompArrayLast(1)].Components = append(c.Discord[c.CompArrayLast(1)].Components, data)
	return c
}

// 入力の追加
// FuncReq: AddLine()
// Request: CustomID,Label, Style.
// Interaction Response Only
func (c *Component) AddInput(data discordgo.TextInput) *Component {
	c.Discord[c.CompArrayLast(1)].Components = append(c.Discord[c.CompArrayLast(1)].Components, data)
	return c
}

// Componentの最後の要素 参照用
func (c *Component) CompArrayLast(n int) int {
	return len(c.Discord) - n
}

// Componentをdiscordgoで使えるように
func (c *Component) Parse() (result []discordgo.MessageComponent) {
	for _, line := range c.Discord {
		result = append(result, line)
	}
	return
}

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aatomu/slashlib"
	"github.com/bwmarrin/discordgo"
)

var (
	token   = flag.String("token", "", "Please Bot Token")
	guildID = flag.String("guild", "", "Please GuildID")
)

func main() {
	//flag入手
	flag.Parse()
	fmt.Println("BotToken   :", *token)
	fmt.Println("GuildID    :", *guildID)

	//bot起動準備
	discord, err := discordgo.New("Bot " + *token)
	Error2Panic("Failed Bot Setup", err)
	//eventトリガー設定
	discord.AddHandler(onReady)
	discord.AddHandler(onInteractionCreate)

	//起動
	err = discord.Open()
	Error2Panic("Failed Session Start", err)
	defer func() {
		err = discord.Close()
		Error2Panic("Failed Session Close", err)
	}()

	//bot停止対策
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func onReady(discord *discordgo.Session, r *discordgo.Ready) {
	//起動メッセージ
	fmt.Println("Bot is OnReady now!")
	cmd := slashlib.Command{}

	RandMin := 1.0
	cmd.
		AddCommand("button", "Generate Button", 0).
		AddCommand("rand", "Generate Rand", 0).
		AddOption(&discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "n",
			Description: "Random % n",
			Required:    true,
			MinValue:    &RandMin,
			MaxValue:    100,
		}).
		CommandCreate(discord, *guildID)
}

// メッセージが送られたときにCall
func onInteractionCreate(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	iData := slashlib.InteractionViewAndEdit(discord, i)

	res := slashlib.InteractionResponse{
		Discord:     discord,
		Interaction: i.Interaction,
	}

	switch iData.Check {
	case slashlib.SlashCommand:
		switch iData.Command.Name {
		case "button":
			err := res.Thinking(false)
			ErrorCheck("Failed ", err)
			time.Sleep(5 * time.Second)
			_, err = res.Follow(&discordgo.WebhookParams{
				Content: "It is Button?",
				Components: new(slashlib.Component).AddLine().
					AddButton(discordgo.Button{
						Label:    "It is Night",
						Style:    1,
						CustomID: "sw1",
					}).
					AddButton(discordgo.Button{
						Label:    "Code Block",
						Style:    1,
						CustomID: "sw2",
					}).
					AddButton(discordgo.Button{
						Label: "Library Link",
						Style: 5,
						URL:   "http://github.com/aatomu/slashlib",
					}).Parse(),
			})
			ErrorCheck("Failed ", err)
		case "rand":
			rand.Seed(time.Now().UnixNano())
			random := rand.Intn(int(iData.Command.Options[0].Value.(float64)))
			err := res.Reply(&discordgo.InteractionResponseData{
				Content: "Rand = " + fmt.Sprint(random),
			})
			ErrorCheck("Failed Send", err)
		}
	case slashlib.ComponentCommand:
		switch iData.Component.CustomID {
		case "sw1":
			res.Reply(&discordgo.InteractionResponseData{
				Content: "Is that true?",
			})
		case "sw2":
			res.Reply(&discordgo.InteractionResponseData{
				Content: "```xl\n'hello? world!'```",
			})
		}
	case slashlib.SubmitCommand:
	}
}

func Error2Panic(comment string, err error) {
	if err != nil {
		log.Println(comment)
		panic(err)
	}
}
func ErrorCheck(comment string, err error) {
	if err != nil {
		log.Println(comment)
		log.Println(err)
	}
}

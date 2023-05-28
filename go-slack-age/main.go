package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
)
func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent){
   for event := range analyticsChannel{
	fmt.Println("Command Events")
	fmt.Println(event.Timestamp)
	fmt.Println(event.Command)
	fmt.Println(event.Parameters)
	fmt.Println(event.Event)
	fmt.Println()
   }
}
func main(){
	 os.Setenv("SLACK_BOT_TOKEN","xoxb-4994700834663-5009190862995-wunkUq5fBNYebgysjwzALs7L")
	 os.Setenv("SLACK_APP_TOKEN","xapp-1-A050BP015B6-5021933123185-7e629f69318fa8c5a225cfeb2fc2d58f0809b7466891b3ce06ae2d30229e114a")
	 bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	 go printCommandEvents(bot.CommandEvents())
	 bot.Command("my yob is <year>", &slacker.CommandDefinition{
		Description: "yob calculator",
		Examples:    []string{"my yob is 2010"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter){
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil{
				println("error")
			}
			t = time.Time
			age := t.Month()-yob
			r := fmt.Sprintf("age is %d", age)
			response.Reply(r)
		},
	 })
	 ctx, cancel := context.WithCancel(context.Background())
	 defer cancel()
	 err := bot.Listen(ctx)
	 if err != nil{
            log.Fatal(err)
	 }
}
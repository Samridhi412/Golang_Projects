package main

import (
	"fmt"
	"os"
	"github.com/slack-go/slack"
)
func main(){
	os.Setenv("SLACK_BOT_TOKEN", "xoxb-4994700834663-4994711749831-NXzxEnOg5qr6mIDe2XL0pa7u")
	os.Setenv("CHANNEL_ID","C050BLML25A")
   api := slack.New(os.Getenv("SLACK_BOT_TOKEN"))
   channelArr := []string{os.Getenv("CHANNEL_ID")}
   fileArr := []string{"categories_202302281316.csv"}
   for i := 0; i<len(fileArr); i++{
	params := slack.FileUploadParameters{
		Channels: channelArr,
		File: fileArr[i],
	}
	file, err := api.UploadFile(params)
	if err != nil{
		fmt.Printf("%s\n", err)
	    return
	}
	fmt.Printf("Name: %s, URL: %s\n", file.Name, file.URL)
   }
}
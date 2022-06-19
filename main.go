package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var lastTimestamp time.Time = time.Now()

func main() {
	config := oauth1.NewConfig("", "")
	token := oauth1.NewToken("", "")
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	for {
		fmt.Println("Executed!")
		mentions, _, err := client.Timelines.MentionTimeline(&twitter.MentionTimelineParams{
			Count: 20,
		})
		if err != nil {
			panic(err)
		}
		for _, tweet := range mentions {
			ts, err := tweet.CreatedAtTime()
			if err != nil {
				fmt.Println("Error getting timestamp of tweet with ID: ", tweet.ID, err)
				continue
			}
			if ts.Before(lastTimestamp) || ts.Equal(lastTimestamp) {
				continue
			}
			if tweet.InReplyToStatusID == 0 {
				fmt.Println("This tweet is not a reply: ", tweet.ID)
				_, _, err = client.Statuses.Update(`@`+tweet.User.ScreenName+` Hola! Por favor, etiquetame debajo de un tweet que contenga un video.`, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
				if err != nil {
					fmt.Println("Error sending a response to the user: ", err)
				}
				continue
			}
			replyTweet, _, err := client.Statuses.Show(tweet.InReplyToStatusID, nil)
			if err != nil {
				fmt.Println("Error getting original tweet info: ", err)
				continue
			}
			if len(replyTweet.Entities.Media) == 0 {
				fmt.Println("This tweet does not contain any media: ", replyTweet.ID)

				_, _, err = client.Statuses.Update(`@`+tweet.User.ScreenName+` Hola! Lamentablemente, no puedo encontrar el vídeo que intentas descargar :(`, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
				if err != nil {
					fmt.Println("Error sending a response to the user: ", err)
				}
				continue
			}
			if !strings.Contains(replyTweet.Entities.Media[0].URLEntity.ExpandedURL, "video") {
				fmt.Println("This tweet does not contain a video: ", replyTweet.ID)

				_, _, err = client.Statuses.Update(`@`+tweet.User.ScreenName+` Hola! Lamentablemente, no puedo encontrar el vídeo que intentas descargar :(`, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
				if err != nil {
					fmt.Println("Error sending a response to the user: ", err)
				}
				continue
			}

			_, _, err = client.Statuses.Update(`@`+tweet.User.ScreenName+` Hola! Acá tenés tu link: 
		
		https://ssstwitter.com/`+tweet.InReplyToScreenName+`/status/`+replyTweet.IDStr, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
			if err != nil {
				fmt.Println("Error sending a response to the user: ", err)
			}
		}
		lastTimestamp, err = mentions[0].CreatedAtTime()
		if err != nil {
			lastTimestamp = time.Now()
			fmt.Println("Error getting last timestamps of mentions: ", err)
		}
		time.Sleep(40 * time.Second)
	}

}

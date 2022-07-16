package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

var lastTimestamp time.Time = time.Now()

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	config := oauth1.NewConfig(os.Getenv("consumer_key"), os.Getenv("consumer_secret"))
	token := oauth1.NewToken(os.Getenv("access_token"), os.Getenv("token_secret"))
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)
	var wg sync.WaitGroup
	for {
		fmt.Println("Executed!")
		mentions, _, err := client.Timelines.MentionTimeline(&twitter.MentionTimelineParams{
			Count: 1,
		})
		if err != nil {
			panic(err)
		}
		for _, tweet := range mentions {
			wg.Add(1)
			go func(tweet twitter.Tweet) {
				defer wg.Done()
				ts, err := tweet.CreatedAtTime()
				if err != nil {
					fmt.Println("Error getting timestamp of tweet with ID: ", tweet.ID, err)
					return
				}
				if ts.Before(lastTimestamp) || ts.Equal(lastTimestamp) {
					return
				}

				if tweet.InReplyToStatusID == 0 {
					fmt.Println("This tweet is not a reply: ", tweet.ID)
					_, _, err = client.Statuses.Update(`@`+tweet.User.ScreenName+` Hola! Por favor, etiquetame debajo de un tweet que contenga un video.`, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
					if err != nil {
						fmt.Println("Error sending a response to the user: ", err)
					}
					return
				}
				replyTweet, _, err := client.Statuses.Show(tweet.InReplyToStatusID, nil)
				if err != nil {
					fmt.Println("Error getting original tweet info: ", err)
					return
				}
				if len(replyTweet.Entities.Media) == 0 {
					_, _, err = client.Statuses.Update(`@`+tweet.User.ScreenName+` Hola! Acá tenés tu link: 
			
					https://ssstwitter.com/`+tweet.InReplyToScreenName+`/status/`+replyTweet.IDStr, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
					if err != nil {
						fmt.Println("Error sending a response to the user: ", err)
					}
					return
				}
				if !strings.Contains(replyTweet.Entities.Media[0].URLEntity.ExpandedURL, "video") {
					fmt.Println("This tweet does not contain a video: ", replyTweet.ID)

					_, _, err = client.Statuses.Update(`@`+tweet.User.ScreenName+` Hola! Lamentablemente, no puedo encontrar el vídeo que intentas descargar :(`, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
					if err != nil {
						fmt.Println("Error sending a response to the user: ", err)
					}
					return
				}

				_, _, err = client.Statuses.Update(`@`+tweet.User.ScreenName+` Hola! Acá tenés tu link: 
			
			https://ssstwitter.com/`+tweet.InReplyToScreenName+`/status/`+replyTweet.IDStr, &twitter.StatusUpdateParams{InReplyToStatusID: tweet.ID})
				if err != nil {
					fmt.Println("Error sending a response to the user: ", err)
				}
			}(tweet)

		}
		wg.Wait()
		lastTimestamp, err = mentions[0].CreatedAtTime()
		if err != nil {
			lastTimestamp = time.Now()
			fmt.Println("Error getting last timestamps of mentions: ", err)
		}
		time.Sleep(40 * time.Second)
	}

}

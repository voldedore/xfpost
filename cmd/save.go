package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

// pageParam is default XF2 param for number of page
const pageParam = "page-"
const defaultOutputFile = "result.json"

// breakDuration is set in order to prevent DOS
const breakDuration = 3

var pageFrom int
var pageTo int
var outputFile string

// User is a post author
type User struct {
	ID   int64  `json:"user_id"`
	Name string `json:"user_name"`
}

// Message is a post, literally
type Message struct {
	CreatedBy *User     `json:"posted_by"`
	Time      time.Time `json:"message_time"`
	Body      string    `json:"message_body"`
	URL       string    `json:"message_url"`
}

// saveCmd represents the get command
var saveCmd = &cobra.Command{
	Use:   "get",
	Short: "Save all posts of a thread",
	Long:  `Usage: xfpost get https://xf.com/t/thread-url.123/`,
	Run: func(cmd *cobra.Command, args []string) {
		mainProcess(args[0])
	},
	Args: cobra.ExactArgs(1),
}

var messages = []*Message{}

func init() {
	saveCmd.Flags().IntVarP(&pageFrom, "from", "f", 1, "fetch from page")
	saveCmd.Flags().IntVarP(&pageTo, "to", "t", 1, "fetch to page")
	saveCmd.Flags().StringVarP(&outputFile, "output", "o", defaultOutputFile, "output file")

	rootCmd.AddCommand(saveCmd)
}

func mainProcess(url string) {
	for i := pageFrom; i <= pageTo; i++ {
		var exactURL *goquery.Document
		if strings.LastIndex(url, "/")+1 == len(url) {
			exactURL = getDocument(url + pageParam + strconv.Itoa(i))
		} else {
			exactURL = getDocument(url + "/" + pageParam + strconv.Itoa(i))
		}
		parseHTML(exactURL)
		time.Sleep(breakDuration * time.Second) // Prevent DOS
	}
	writeToFile()
}

func writeToFile() {
	clientsFile, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()
	str, _ := json.Marshal(messages)
	clientsFile.Write(str)
	log.Printf("Write successfully to file %s.\n", outputFile)
}

// parseHTML receives Document, then process and output data
func parseHTML(doc *goquery.Document) {

	doc.Find("article.message").Each(func(i int, s *goquery.Selection) {
		userIDStr, _ := s.Find("h4.message-name a").Attr("data-user-id")
		uid, _ := strconv.ParseInt(userIDStr, 0, 64)
		user := &User{ID: uid}
		user.Name = s.Find("h4.message-name a").Text()

		message := &Message{Body: s.Find("article.message-body .bbWrapper").Text(),
			CreatedBy: user}
		messageTimeRaw, _ := s.Find(".message-attribution time").Attr("data-time")
		messageTime, _ := strconv.ParseInt(messageTimeRaw, 0, 64)
		message.Time = time.Unix(messageTime, 0)

		messagePermalink, _ := s.Find(".message-attribution a").Attr("href")
		message.URL = messagePermalink
		messages = append(messages, message)
		// fmt.Printf("Post %d on %s by %s (%d): %s\n", i+1,
		// 	message.Time.Format(time.UnixDate), user.Name, user.ID, message.Body)
		// fmt.Printf("Details: %s\n", messagePermalink)
	})

}

// getDocument accepts a thread URL, then fetches the page and returns a goQuery document
func getDocument(url string) *goquery.Document {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil
	}

	req.Header.Set("User-Agent", "telegram-bot:xfpost:v1.0.0")

	resp, err := client.Do(req)

	if err != nil {
		log.Print(err)
	} else {

		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("GET %s successfully\n", url)
		return doc
	}

	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
)

var (
	jsonFilePath string
)

type User struct {
	ID              int    `json:"id"`
	URLName         string `json:"url_name"`
	ProfileImageUrl string `json:"profile_image_url"`
}

type Tag struct {
	Name string `json:"name"`
}

type Comment struct {
	ID      int    `json:"id"`
	UUID    string `json:"uuid"`
	User    User   `json:"user"`
	RawBody string `json:"raw_body"`
}

type Article struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	User      User      `json:"user"`
	Title     string    `json:"title"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Tags      []Tag     `json:"tags"`
	URL       string    `json:"url"`
	RawBody   string    `json:"raw_body"`
	Comments  []Comment `json:"comments"`
}

func setAppInfo(app *cli.App) {
	app.Name = "qt2md - convert from Qiita::Team backup json to markdowns"
	app.Usage = "load your Qiita::Team backup json and then markdowns will be created in current directory."
	app.Version = "0.0.1"
}

func setFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "export.json",
			Usage: "Qiita::Team export json path eg. ~/Downloda/exprot.json",
		},
	}

}

func loadJson(jsonFilePath string) ([]Article, error) {
	fmt.Println(jsonFilePath)
	b, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}
	var articles []Article
	json.Unmarshal(b, &articles)
	return articles, nil
}

func writeMarkDown(filePath string, article Article) error {
	title := fmt.Sprintf("# %s\n", article.Title)
	err := ioutil.WriteFile(filePath, []byte(title+article.RawBody), 0644)
	if err != nil {
		return err
	}
	return nil
}

func convertToMarkDown(articles []Article) error {
	for _, article := range articles {
		fmt.Println(article.Title)
		err := writeMarkDown(article.Title+".md", article)
		if err != nil {
			msg := fmt.Sprintf("convert failed! %s %s", article.Title, err)
			coloredMsg := ansi.Color(msg, "red")
			fmt.Println(coloredMsg)
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	setAppInfo(app)
	setFlags(app)
	app.Action = func(c *cli.Context) {
		jsonFilePath = c.String("file")

		articles, err := loadJson(jsonFilePath)
		if err != nil {
			panic(err)
		}

		err = convertToMarkDown(articles)
		if err != nil {
			panic(err)
		}

	}
	app.Run(os.Args)
}

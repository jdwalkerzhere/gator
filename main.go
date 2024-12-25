package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jdwalkerzhere/gator/internal/config"
	"github.com/jdwalkerzhere/gator/internal/database"
	_ "github.com/lib/pq"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	registry map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.registry == nil {
		c.registry = make(map[string]func(*state, command) error)
	}
	c.registry[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	function, ok := c.registry[cmd.name]
	if !ok {
		return fmt.Errorf("Command not registered, please select from the valid commands")
	}
	// We've got a registered cli command, now run it
	if err := function(s, cmd); err != nil {
		return err
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Must provide arguments to `login` command, None provided")
	}
	newUser := cmd.arguments[0]
	if _, err := s.db.GetUser(context.Background(), newUser); err != nil {
		fmt.Printf("User [%s] doesn't exists, register username first", newUser)
		os.Exit(1)
	}
	if err := config.SetUser(s.config, newUser); err != nil {
		return err
	}
	fmt.Printf("User [%s] set as current user\n", newUser)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Must provide argument to `register` command, None provided")
	}
	if user, err := s.db.GetUser(context.Background(), cmd.arguments[0]); err == nil {
		fmt.Printf("User [%s] already exists, register with different username", user.Name)
		os.Exit(1)
	}
	now := time.Now()
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      cmd.arguments[0],
	}
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	config.SetUser(s.config, user.Name)
	fmt.Printf("New user [%s] created", user.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return err
	}
	return nil
}

func handlerGetUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user == s.config.CurrentUser {
			fmt.Printf("* %s (current)\n", user)
			continue
		}
		fmt.Printf("* %s\n", user)
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(*feed)
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	ctxRequest, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	ctxRequest.Header.Set("User-Agent", "gator")
	client := http.Client{}
	resp, err := client.Do(ctxRequest)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	feed := RSSFeed{}
	if err := xml.Unmarshal(body, &feed); err != nil {
		return &RSSFeed{}, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Link = html.UnescapeString(feed.Channel.Link)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Link = html.UnescapeString(item.Link)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
	return &feed, nil
}

func main() {
	cnfg, err := config.Read()
	if err != nil {
		log.Panic(err)
	}
	db, err := sql.Open("postgres", cnfg.DbUrl)
	if err != nil {
		log.Panic(err)
	}
	dbQueries := database.New(db)
	state := state{dbQueries, &cnfg}
	cmds := commands{}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	arguments := os.Args

	if len(arguments) < 2 {
		fmt.Println("No arguments given, please run with arguments")
		os.Exit(1)
	}
	cmdName := arguments[1]
	cmdArgs := arguments[2:]
	userCommand := command{cmdName, cmdArgs}
	if err := cmds.run(&state, userCommand); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

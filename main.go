package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jdwalkerzhere/gator/internal/config"
	"github.com/jdwalkerzhere/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No username provided, please provide one\n")
	}
	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("No User [%s] Exists, Must be registered first\n", username)
		os.Exit(1)
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Current User Set To [%s]\n", s.cfg.CurrentUser)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No username provided, please provide one\n")
	}
	username := cmd.args[0]
	timeNow := time.Now()
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Name:      username,
	}
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		fmt.Printf("User [%s] Already Exists, Provide new username\n", username)
		os.Exit(1)
	}
	s.cfg.SetUser(user.Name)
	fmt.Printf("User [%s] Created: %v\n", s.cfg.CurrentUser, user)
	return nil
}

func handlerReset(s *state, _ command) error {
	s.db.Reset(context.Background())
	return nil
}

func handlerGetUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return &RSSFeed{}, err
	}
	return &rssFeed, nil
}

func handlerAgg(s *state, _ command) error {
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(*rss)
	return nil
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	command, ok := c.commandMap[cmd.name]
	if !ok {
		return fmt.Errorf("No command [%s] registered in the CLI\n", cmd.name)
	}
	err := command(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// initializing state and command registry
	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		panic(err)
	}
	dbQueries := database.New(db)

	stateNew := state{dbQueries, &cfg}
	cmds := commands{make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)

	// fetching user cli args
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No arguments provided, please run again")
		os.Exit(1)
	}

	cmd := command{args[1], args[2:]}
	err = cmds.run(&stateNew, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

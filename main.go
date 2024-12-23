package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jdwalkerzhere/gator/internal/config"
)

type state struct {
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
	if err := config.SetUser(s.config, newUser); err != nil {
		return err
	}
	fmt.Printf("User [%s] set as current user\n", newUser)
	return nil
}

func main() {
	cnfg, err := config.Read()
	if err != nil {
		log.Panic(err)
	}
	state := state{&cnfg}
	cmds := commands{}
	cmds.register("login", handlerLogin)
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

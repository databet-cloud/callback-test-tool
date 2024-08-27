package prompt

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

const (
	exitCommand          = "exit"
	defaultSelectionSize = 15
)

type Tree struct {
	Label             string
	Commands          func() []*Command
	ReturnAfterAction bool
}

type Command struct {
	Key    string
	Action func()
	Tree   *Tree
}

type Stringer interface {
	String() string
}

func Select[T Stringer](label string, values ...T) (value T, err error) {
	keys := make([]string, len(values))
	for i, v := range values {
		keys[i] = v.String()
	}

	selection := promptui.Select{Label: label, Items: keys, Size: defaultSelectionSize}

	i, _, err := selection.Run()
	if err != nil {
		return value, err
	}

	return values[i], nil
}

func ProcessCommands(tree *Tree) {
	list := tree.Commands()
	commands := make(map[string]*Command, len(list))
	keys := make([]string, 0, len(list)+1)

	for _, cmd := range list {
		commands[cmd.Key] = cmd
		keys = append(keys, cmd.Key)
	}

	commandPrompt := promptui.Select{Label: tree.Label, Items: append(keys, exitCommand), Size: defaultSelectionSize}

	// nolint:wsl // its ok, because that needed to re-render menu after action
	for {
		_, input, err := commandPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			continue
		}

		if input == exitCommand {
			break
		}

		cmd := commands[input]
		if cmd.Action != nil {
			cmd.Action()
		}

		if cmd.Tree != nil {
			ProcessCommands(cmd.Tree)
		}

		if tree.ReturnAfterAction {
			break
		}
	}
}

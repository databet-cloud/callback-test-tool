package prompt

import (
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

func Float(label string) float64 {
	prompt := &promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			_, err := strconv.ParseFloat(input, 64)
			if err != nil {
				return fmt.Errorf("fail to parse %s as float", input)
			}
			return nil
		},
	}

	valueStr, err := prompt.Run()
	if err != nil {
		return 0
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0
	}

	return value
}

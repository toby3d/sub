package domain

import (
	"errors"
	"fmt"
)

type Action struct {
	action string
}

var (
	ActionUnd      = Action{action: ""}         // "und"
	ActionBlock    = Action{action: "block"}    // "block"
	ActionChannels = Action{action: "channels"} // "channels"
	ActionFollow   = Action{action: "follow"}   // "follow"
	ActionMute     = Action{action: "mute"}     // "mute"
	ActionPreview  = Action{action: "preview"}  // "preview"
	ActionSearch   = Action{action: "search"}   // "search"
	ActionTimeline = Action{action: "timeline"} // "timeline"
	ActionUnfollow = Action{action: "unfollow"} // "unfollow"
)

var ErrActionSyntax = errors.New("unknown or unsupported action")

var stringsActions = map[string]Action{
	ActionBlock.action:    ActionBlock,
	ActionChannels.action: ActionChannels,
	ActionFollow.action:   ActionFollow,
	ActionMute.action:     ActionMute,
	ActionPreview.action:  ActionPreview,
	ActionSearch.action:   ActionSearch,
	ActionTimeline.action: ActionTimeline,
	ActionUnfollow.action: ActionUnfollow,
}

func ParseAction(src string) (Action, error) {
	if action, ok := stringsActions[src]; ok {
		return action, nil
	}

	return ActionUnd, fmt.Errorf("%w: %s", ErrActionSyntax, src)
}

func (a Action) String() string {
	if a.action != "" {
		return a.action
	}

	return "und"
}

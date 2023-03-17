package http

import (
	"fmt"
	"net/http"

	"source.toby3d.me/toby3d/sub/internal/domain"
)

type (
	RequestChannelsCreate struct {
		Action domain.Action // channels
		Name   string
	}

	RequestChannelsUpdate struct {
		Action  domain.Action // channels
		Channel string
		Name    string
	}

	RequestChannelsOrder struct {
		Action  domain.Action // channels
		Method  domain.Method // order
		Channel []string
	}

	RequestChannelsDelete struct {
		Action  domain.Action // channels
		Method  domain.Method // delete
		Channel string
	}

	ResponseChannel struct {
		UID  string `json:"uid"`
		Name string `json:"name"`
	}

	RequestTimelines struct {
		After   any
		Before  any
		Action  domain.Action
		Channel string
	}

	ResponseChannels struct {
		Channels []ResponseChannelsChannel `json:"channels"`
	}

	ResponseChannelsChannel struct {
		UID    string `json:"uid"`
		Name   string `json:"name"`
		Unread uint   `json:"unread"`
	}

	ResponseTimelines struct {
		Paging ResponsePaging  `json:"paging"`
		Items  []ResponseEntry `json:"items"`
	}

	ResponsePaging struct {
		After  string `json:"after"`
		Before string `json:"before"`
	}

	ResponseEntry struct {
		Checkin     CardPlace       `json:"checkin"`
		Author      CardPeople      `json:"author"`
		Content     ResponseContent `json:"content"`
		Type        string          `json:"type"`
		Published   string          `json:"published"`
		URL         string          `json:"url"`
		UID         string          `json:"uid"`
		Name        string          `json:"name"`
		Summary     string          `json:"summary"`
		ID          string          `json:"_id"`
		Video       []string        `json:"video"`
		Audio       []string        `json:"audio"`
		LikeOf      []string        `json:"like-of"`
		RepostOf    []string        `json:"repost-of"`
		BookmarkOf  []string        `json:"bookmark-of"`
		InReplyTo   []string        `json:"in-reply-of"`
		Syndication []string        `json:"syndication"`
		Photo       []string        `json:"photo"`
		Category    []string        `json:"category"`
		IsRead      bool            `json:"_is_read"`
	}

	ResponseAuthor struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
		Photo string `json:"photo"`
	}

	ResponseContent struct {
		Text string `json:"text"`
		HTML string `json:"html"`
	}

	CardPlace struct {
		Type          string `json:"type"`
		Name          string `json:"name"`
		URL           string `json:"url"`
		Latitude      string `json:"latitude"`
		Longitude     string `json:"longitude"`
		StreetAddress string `json:"street-address"`
		Locality      string `json:"locality"`
		Region        string `json:"region"`
		Country       string `json:"country"`
	}

	CardPeople struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
		Photo string `json:"photo"`
	}

	ResponseSource struct {
		URL   string `json:"url"`
		Name  string `json:"name"`
		Photo string `json:"photo"`
		ID    string `json:"_id"`
	}
)

func NewResponseChannels(channels ...domain.Channel) *ResponseChannels {
	out := &ResponseChannels{
		Channels: make([]ResponseChannelsChannel, len(channels)),
	}

	for i := range channels {
		out.Channels[i] = ResponseChannelsChannel{
			UID:    channels[i].UID,
			Name:   channels[i].Name,
			Unread: 0,
		}
	}

	return out
}

func NewResponseChannel(c *domain.Channel) *ResponseChannel {
	out := new(ResponseChannel)

	if c == nil {
		return out
	}
	out.UID = c.UID
	out.Name = c.Name

	return out
}

func (r *RequestChannelsCreate) bind(req *http.Request) error {
	var err error
	if r.Action, err = domain.ParseAction(req.PostFormValue("action")); err != nil {
		return fmt.Errorf("cannot decode channel delete request: %w", err)
	}

	if r.Action != domain.ActionChannels {
		return fmt.Errorf("expect '%s' action, got '%s'", domain.ActionChannels, r.Action)
	}

	if r.Name = req.PostFormValue("name"); r.Name == "" {
		return fmt.Errorf("expect channel name value, but it's not provided")
	}

	return nil
}

func (r *RequestChannelsUpdate) bind(req *http.Request) error {
	var err error
	if r.Action, err = domain.ParseAction(req.PostFormValue("action")); err != nil {
		return fmt.Errorf("cannot decode channel delete request: %w", err)
	}

	if r.Action != domain.ActionChannels {
		return fmt.Errorf("expect '%s' action, got '%s'", domain.ActionChannels, r.Action)
	}

	if r.Channel = req.PostFormValue("channel"); r.Channel == "" {
		return fmt.Errorf("expect channel UID value, but it's not provided")
	}

	if r.Name = req.PostFormValue("name"); r.Name == "" {
		return fmt.Errorf("expect channel name value, but it's not provided")
	}

	return nil
}

func (r *RequestChannelsOrder) bind(req *http.Request) error {
	var err error
	if r.Action, err = domain.ParseAction(req.PostFormValue("action")); err != nil {
		return fmt.Errorf("cannot decode channel delete request: %w", err)
	}

	if r.Action != domain.ActionChannels {
		return fmt.Errorf("expect '%s' action, got '%s'", domain.ActionChannels, r.Action)
	}

	if r.Method, err = domain.ParseMethod(req.PostFormValue("method")); err != nil {
		return fmt.Errorf("cannot decode channel delete request: %w", err)
	}

	if r.Method != domain.MethodOrder {
		return fmt.Errorf("expect '%s' method, got '%s'", domain.MethodOrder, r.Method)
	}

	println(req.PostFormValue("channels[]"))

	return nil
}

func (r *RequestChannelsDelete) bind(req *http.Request) error {
	var err error
	if r.Action, err = domain.ParseAction(req.PostFormValue("action")); err != nil {
		return fmt.Errorf("cannot decode channel delete request: %w", err)
	}

	if r.Action != domain.ActionChannels {
		return fmt.Errorf("expect '%s' action, got '%s'", domain.ActionChannels, r.Action)
	}

	if r.Method, err = domain.ParseMethod(req.PostFormValue("method")); err != nil {
		return fmt.Errorf("cannot decode channel delete request: %w", err)
	}

	if r.Method != domain.MethodDelete {
		return fmt.Errorf("expect '%s' method, got '%s'", domain.MethodDelete, r.Method)
	}

	if r.Channel = req.PostFormValue("channel"); r.Channel == "" {
		return fmt.Errorf("expect channel UID value, but it's not provided")
	}

	return nil
}

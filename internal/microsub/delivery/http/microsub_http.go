package http

import (
	"net/http"

	"github.com/goccy/go-json"

	"source.toby3d.me/toby3d/sub/internal/channel"
	"source.toby3d.me/toby3d/sub/internal/common"
	"source.toby3d.me/toby3d/sub/internal/domain"
)

type Handler struct {
	channels channel.UseCase
}

func NewHandler(channels channel.UseCase) *Handler {
	return &Handler{
		channels: channels,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(*domain.User)
	encoder := json.NewEncoder(w)

	switch r.Method {
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	case "", http.MethodGet:
		action, err := domain.ParseAction(r.URL.Query().Get("action"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		switch action {
		case domain.ActionChannels:
			channels, err := h.channels.Fetch(r.Context(), *user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			w.Header().Set(common.HeaderContentType, common.MIMEApplicationJSONCharsetUTF8)
			_ = encoder.Encode(NewResponseChannels(channels...))
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		switch {
		default:
			req := new(RequestChannelsCreate)
			if err := req.bind(r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}

			result, err := h.channels.Create(r.Context(), *user, req.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			w.Header().Set(common.HeaderContentType, common.MIMEApplicationJSONCharsetUTF8)
			_ = encoder.Encode(NewResponseChannel(result))
		case r.PostForm.Has("method"):
			req := new(RequestChannelsDelete)
			if err := req.bind(r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}

			if err := h.channels.Delete(r.Context(), *user, req.Channel); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			w.WriteHeader(http.StatusNoContent)
		case r.PostForm.Has("channels[]"):
			req := new(RequestChannelsOrder)
			if err := req.bind(r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}

			if err := h.channels.Order(r.Context(), *user, req.Channel); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			w.WriteHeader(http.StatusNoContent)
		case r.PostForm.Has("channel"):
			req := new(RequestChannelsUpdate)
			if err := req.bind(r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			}

			result, err := h.channels.Update(r.Context(), *user, req.Channel, req.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

			w.Header().Set(common.HeaderContentType, common.MIMEApplicationJSONCharsetUTF8)
			_ = encoder.Encode(NewResponseChannel(result))
		}
	}
}

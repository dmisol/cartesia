package pkg

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dmisol/cartesia/pkg/types"
	"github.com/gorilla/websocket"
)

type Session struct {
	conn   *websocket.Conn
	onDone func()
	onData func(id string, f32 []byte, final bool)
}

func NewSession(ctx context.Context, key string,
	onData func(id string, f32 []byte, fin bool), onDone func()) (*Session, error) {

	s := &Session{
		onData: onData,
		onDone: onDone,
	}
	c, r, err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://api.cartesia.ai/v0/audio/websocket?api_key=%s", key), nil)
	if err != nil {
		return nil, err
	}
	s.conn = c

	s.Println("dial status", r.Status)

	go s.run(ctx)
	return s, nil
}

func (s *Session) run(ctx context.Context) {
	defer s.conn.Close()
	defer func() {
		if s.onDone != nil {
			s.onDone()
		}
	}()

	done := make(chan bool)
	go func() {
		defer func() { done <- true }()
		for {
			mt, b, err := s.conn.ReadMessage()
			if err != nil {
				s.Println("read", err)
				return
			}
			if mt != websocket.TextMessage {
				continue
			}
			resp := &types.Response{}
			if err = json.Unmarshal(b, resp); err != nil {
				s.Println("unmarshal", err)
				continue
			}
			if s.onData != nil {
				if resp.Done || len(resp.Data) == 0 {
					s.onData(resp.ContextId, nil, resp.Done)
					continue
				}
				ba, err := base64.StdEncoding.DecodeString(resp.Data)
				if err != nil {
					s.Println("decoding", err)
					continue
				}
				s.onData(resp.ContextId, ba, false)
			}
		}
	}()

	select {
	case <-ctx.Done():
		s.Println("closed (ctx)")
		return
	case <-done:
		s.Println("closed (chan)")
		return
	}
}

func (s *Session) TTS(id string, text string, model types.Model, voice types.Voice) {
	r := &types.Request{

		ContextId: id,
		Data: types.RequestData{
			Text:  text,
			Model: string(model),
			Voice: voice,
		},
	}
	/*
		b, err := json.Marshal(r)
		if err != nil {
			s.Println("snd marshal", err)
		}
		s.Println(string(b))
	*/

	if err := s.conn.WriteJSON(r); err != nil {
		s.Println("ws send", err)
	}
}
func (s *Session) Close() {

}

func (s *Session) Println(i ...interface{}) {
	log.Println("sess", i)
}

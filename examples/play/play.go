package main

// https://github.com/gordonklaus/portaudio
// apt-get install portaudio19-dev

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"log"
	"math"
	"os"
	"time"

	"github.com/dmisol/cartesia/pkg"
	"github.com/dmisol/cartesia/pkg/model"
	"github.com/dmisol/cartesia/pkg/voice"
	"github.com/google/uuid"
	"github.com/gordonklaus/portaudio"
)

const chunkSize = 4608

type Config struct {
	Key string `json:"key"`
}

var (
	cancel context.CancelFunc
	ctx    context.Context

	stream *portaudio.Stream

	dataChan = make(chan []int16, 100)
	final    bool
)

func main() {
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	b, _ := os.ReadFile("conf.json")
	c := &Config{}
	if err := json.Unmarshal(b, c); err != nil {
		log.Fatal(err)
	}
	s, err := pkg.NewSession(ctx, c.Key, onData, nil)
	if err != nil {
		log.Fatal(err)
	}

	portaudio.Initialize()
	defer portaudio.Terminate()

	stream, err = portaudio.OpenDefaultStream(0, 1, 44100, chunkSize, feedInt16)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	if err = stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()

	s.TTS(uuid.NewString(), "If you try the API in Gitbook, you will need to manually specify the API key in the X-API-Key header field. The authentication pane does not currently work due to a GitBook bug. We're following up on that with them.", model.SonicTurboEnglish, voice.Elon)

	log.Println("waiting for ctx done")
	<-ctx.Done()
	log.Println("ctx done")

}

func onData(id string, data []byte, fin bool) {
	log.Println(len(data), fin)

	if !fin {
		b := make([]byte, 4)
		tail := data
		out := make([]int16, 0)
		for {
			b, tail = tail[:4], tail[4:]

			bits := binary.LittleEndian.Uint32(b)
			f := 32000 * math.Float32frombits(bits)

			i := int16(f)
			out = append(out, i)

			if len(tail) <= 3 {
				break
			}
		}

		log.Println("int16:", len(out))
		dataChan <- out
	}

	if fin {
		final = true
		// cancel()
	}

}

func feedInt16(out []int16) {
	in := <-dataChan
	copy(out, in)
	if (len(dataChan) == 0) && final {
		go func() {
			time.Sleep(time.Second)
			cancel()
			log.Println("cancel()")
		}()
	}
}

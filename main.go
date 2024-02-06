package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

type Config struct {
	Nsec   string   `json:"nsec"`
	Relays []string `json:"relays"`
}

type Article struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (s *Article) Publish(cfg *Config) error {

	log.Println(s)
	log.Println(cfg)

	ctx := context.Background()

	var sk string
	var pub string
	if _, s, err := nip19.Decode(cfg.Nsec); err == nil {
		sk = s.(string)
		if pub, err = nostr.GetPublicKey(s.(string)); err != nil {
			return err
		}
	} else {
		return err
	}

	var tags nostr.Tags
	tags = append(tags, nostr.Tag{"title", s.Title})
	for _, v := range s.Tags {
		tags = append(tags, nostr.Tag{"t", v})
	}

	e := nostr.Event{
		Kind:      nostr.KindArticle,
		Content:   s.Content,
		CreatedAt: nostr.Now(),
		PubKey:    pub,
		Tags:      tags,
	}

	err := e.Sign(sk)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, r := range cfg.Relays {
		wg.Add(1)

		go func(relayUrl string) {
			defer wg.Done()

			relay, err := nostr.RelayConnect(ctx, relayUrl)
			if err != nil {
				log.Println(err)
				return
			}
			defer relay.Close()

			err = relay.Publish(ctx, e)
			if err != nil {
				log.Println(err)
				return
			}
		}(r)
	}
	wg.Wait()

	return nil
}

func loadConfig() (*Config, error) {

	env := os.Getenv("NOSTR_ZET")

	data, err := os.ReadFile(env)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg, nil
}

func main() {

	args := os.Args[1:]

	var tags []string
	for _, v := range strings.Split(args[2], ",") {
		tags = append(tags, strings.Trim(v, " "))
	}

	content, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	a := Article{
		Title:   args[1],
		Content: string(content),
		Tags:    tags,
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = a.Publish(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

package zet

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
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	References []string `json:"references"`
}

func (s *Article) Publish(cfg *Config) error {

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
	for _, v := range s.References {
		tags = append(tags, nostr.Tag{"r", v})
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

func LoadConfig(path string) (*Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Config file: %v", err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg, nil
}

func Main() error {

	args := os.Args[1:]

	var tags []string
	for _, v := range strings.Split(args[2], ",") {
		tags = append(tags, strings.Trim(v, " "))
	}

	var refs []string
	for _, v := range strings.Split(args[3], ",") {
		refs = append(refs, strings.Trim(v, " "))
	}

	content, err := os.ReadFile(args[0])
	if err != nil {
        return err
	}

	a := Article{
		Title:      args[1],
		Content:    string(content),
		Tags:       tags,
		References: refs,
	}

	path, ok := os.LookupEnv("NOSTR_ZET")
	if !ok {
		log.Fatalln("NOSTR_ZET env var not set")
	}

	cfg, err := LoadConfig(path)
	if err != nil {
        return err
	}

	err = a.Publish(cfg)
	if err != nil {
        return err
	}

    return nil
}

package nip23

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

type Config struct {
	Nsec   string   `json:"nsec"`
	Relays []string `json:"relays"`
}

type Article struct {
	Identifier string   `json:"identifier"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Tags       []string `json:"tags"`
	Urls       []string `json:"urls"`
	Events     []string `json:"events"`
}

func (s *Article) Publish(cfg *Config) (string, error) {

	ctx := context.Background()

	var sk string
	var pub string
	if _, s, err := nip19.Decode(cfg.Nsec); err == nil {
		sk = s.(string)
		if pub, err = nostr.GetPublicKey(s.(string)); err != nil {
            return "", err
		}
	} else {
        return "", err
	}

	var tags nostr.Tags

	tags = append(tags, nostr.Tag{"title", s.Title})
    tags = append(tags, nostr.Tag{"d", s.Identifier})

	for _, v := range s.Tags {
		tags = append(tags, nostr.Tag{"t", v})
	}
	for _, v := range s.Urls {
		tags = append(tags, nostr.Tag{"r", v})
	}
	for _, v := range s.Events {
		tags = append(tags, nostr.Tag{"e", v})
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
        return "", err
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

	naddr, err := nip19.EncodeEntity(
        pub,
		nostr.KindArticle,
        s.Identifier,
        cfg.Relays,
    )
	if err != nil {
        return "", err
	}

	return naddr, nil
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

	// Create a flag set and define flags
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// String flag for title
    title := fs.String("title", "", "Title of the document")

    // String slices for tags, URLs, and events
    tags := make([]string, 0)
    urls := make([]string, 0)
    events := make([]string, 0)
    
	// Define a function to append the input string to the slice
	appendFlag := func(slice *[]string) func(string) error {
		return func(s string) error {
			*slice = append(*slice, s)
			return nil
		}
	}

	// Tags
	fs.Func("t", "Tags related to the document (can be specified multiple times)", appendFlag(&tags))

	// URLs
	fs.Func("r", "Reference URLs (can be specified multiple times)", appendFlag(&urls))

	// Events
	fs.Func("e", "Event identifiers (can be specified multiple times)", appendFlag(&events))

	// Parse the flags
	err := fs.Parse(os.Args[2:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	markdownFile := os.Args[1]

	content, err := os.ReadFile(markdownFile)
	if err != nil {
		return err
	}

    a := Article{
        Identifier: time.Now().Format("200601021504"),
        Title: *title,
        Content: string(content),
        Tags: tags,
        Urls: urls,
        Events: events,
    }

 	path, ok := os.LookupEnv("NOSTR")
 	if !ok {
 		log.Fatalln("NOSTR env var not set")
 	}
 
 	cfg, err := LoadConfig(path)
 	if err != nil {
 		return err
 	}

    naddr, err := a.Publish(cfg)
 	if err != nil {
 		return err
 	}

    fmt.Println("[*] Article published")
	fmt.Printf("Markdown File: %s\n", markdownFile)
	fmt.Printf("Identifier: %s\n", a.Identifier)
	fmt.Printf("Title: %s\n", a.Title)
	fmt.Printf("Tags: %v\n", a.Tags)
	fmt.Printf("URLs: %v\n", a.Urls)
	fmt.Printf("Events: %v\n", a.Events)
    fmt.Println(naddr)

	return nil
}

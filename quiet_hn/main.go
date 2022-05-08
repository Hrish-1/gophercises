package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"quiet-hn/hn"
	"sort"
	"strings"
	"sync"
	"time"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	var sc storyCache = storyCache{
		duration:   15 * time.Second,
		numStories: numStories,
	}
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			temp := storyCache{
				numStories: numStories,
				duration:   15 * time.Second,
			}
			temp.stories()
			sc.mutex.Lock()
			sc.cache = temp.cache
			sc.expiration = temp.expiration
			sc.mutex.Unlock()
			<-ticker.C
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := sc.stories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

type storyCache struct {
	cache      []item
	expiration time.Time
	numStories int
	duration   time.Duration
	mutex      sync.Mutex
}

func (s *storyCache) stories() ([]item, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if time.Now().Sub(s.expiration) < 0 {
		return s.cache, nil
	}
	stories, err := getTopStories(s.numStories)
	if err != nil {
		return nil, err
	}
	s.expiration = time.Now().Add(s.duration)
	s.cache = stories
	return s.cache, nil
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}
	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[at:at+need])...)
		at += need
	}
	return stories[:numStories], nil
}

func getStories(ids []int) []item {
	var stories []item
	type story struct {
		story item
		err   error
		idx   int
	}
	resultCh := make(chan story)
	for i := 0; i < len(ids); i++ {
		go func(idx int, id int) {
			var client hn.Client
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- story{err: err, idx: idx}
			}
			item := parseHNItem(hnItem)
			resultCh <- story{story: item, idx: idx}
		}(i, ids[i])
	}

	var results []story
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	for _, res := range results {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.story) {
			stories = append(stories, res.story)
		}
	}
	return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

package crawler

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
)

const LONGMAN_DICTIONARY_DOMAIN = "www.ldoceonline.com"

//go:generate mockery --name Spider
type Spider interface {
	FindWordMeaningsFromDictionary(
		word string,
	) ([]model.WordMeaning, error)
}

type spider struct{}

func NewSpider() Spider {
	return &spider{}
}

func (mySpider spider) FindWordMeaningsFromDictionary(
	word string,
) ([]model.WordMeaning, error) {
	wordMeangins := []model.WordMeaning{}

	// 排序用的編號
	var orderByNo int32 = 0

	c := colly.NewCollector(
		colly.AllowedDomains(LONGMAN_DICTIONARY_DOMAIN),
	)

	// 隨機設定 user agent，避免被網站認出是爬蟲而被網站擋住
	extensions.RandomUserAgent(c)

	var parseHtmlErr error

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("parse html error: %v", err)
		parseHtmlErr = fmt.Errorf(
			"Request URL: %s, failed with response: %v, \nError: %w",
			r.Request.URL,
			r,
			err,
		)
	})

	c.OnHTML("div.content", func(e *colly.HTMLElement) {
		pageTitleWord := strings.TrimSpace(e.DOM.Find("h1.pagetitle").Text())

		e.ForEachWithBreak("span.dictentry", func(i int, dictentry *colly.HTMLElement) bool {
			// 不要抓來自其他字典的解釋，因為只抓來自 Longman Dictionary of Contemporary 就很夠了
			if i > 0 && dictentry.DOM.Is(":has(.dictionary_intro)") {
				return false // break
			}

			dictlink := dictentry.DOM.Find("span.dictlink")
			senses := dictlink.Find("span.Sense:has(span.DEF)")

			if senses.Length() == 0 {
				return true // continue
			}

			partOfSpeech := strings.TrimSpace(dictlink.Find("span.Head span.POS").Text())
			headGram := strings.TrimSpace(dictlink.Find("span.Head span.GRAM").Text())

			// 音標與發音
			pronText := strings.TrimSpace(dictlink.Find("span.Head span.PronCodes").Text())
			ukAudioUrl, ukAudioUrlExists := dictlink.Find("span.speaker.brefile").
				Attr("data-src-mp3")
			usAudioUrl, usAudioUrlExists := dictlink.Find("span.speaker.amefile").
				Attr("data-src-mp3")

			if !ukAudioUrlExists || !usAudioUrlExists {
				return true // continue
			}

			// Find meanings
			senses.Each(func(senseIndex int, sense *goquery.Selection) {
				defGram := strings.TrimSpace(sense.Find("span.GRAM").Text())
				def := sense.Find("span.DEF")

				// 朗文網頁中會在某些單字右上角標注小數字，移除它
				def.Find("span.REFHOMNUM").Remove()
				definition := strings.TrimSpace(def.Text())
				orderByNo += 1

				var queryByWords []string
				if pageTitleWord == word {
					queryByWords = []string{word}
				} else {
					queryByWords = []string{pageTitleWord, word}
				}

				wordMeaning := model.WordMeaning{
					Word:         pageTitleWord,
					PartOfSpeech: partOfSpeech,
					Gram:         headGram,
					Pronunciation: model.Pronunciation{
						Text:       pronText,
						UkAudioUrl: ukAudioUrl,
						UsAudioUrl: usAudioUrl,
					},
					DefGram:      defGram,
					Definition:   definition,
					Examples:     []model.Example{},
					OrderByNo:    orderByNo,
					QueryByWords: queryByWords,
				}

				// Find examples
				sense.ChildrenFiltered("span.GramExa, span.EXAMPLE").
					Each(func(childIndex int, child *goquery.Selection) {
						var example model.Example
						pattern := strings.TrimSpace(
							child.Find("span.PROPFORMPREP, span.PROPFORM").Text(),
						)

						if child.Is(".GramExa") {
							example = model.Example{
								Pattern:  pattern,
								Examples: []model.Sentence{},
							}

							child.Find("span.EXAMPLE").
								Each(func(gramExaExampleIndex int, gramExaExample *goquery.Selection) {
									audioUrl, _ := gramExaExample.Find("span[data-src-mp3]").
										Attr("data-src-mp3")
									text := strings.TrimSpace(gramExaExample.Text())
									example.Examples = append(example.Examples, model.Sentence{
										AudioUrl: audioUrl,
										Text:     text,
									})
								})

						} else {
							audioUrl, _ := child.Find("span[data-src-mp3]").Attr("data-src-mp3")
							example = model.Example{
								Pattern: "",
								Examples: []model.Sentence{
									{
										AudioUrl: audioUrl,
										Text:     strings.TrimSpace(child.Text()),
									},
								},
							}
						}

						wordMeaning.Examples = append(wordMeaning.Examples, example)
					})

				wordMeangins = append(wordMeangins, wordMeaning)
			})

			return true
		})
	})

	// Start scraping
	c.Visit(fmt.Sprintf("https://%s/dictionary/%s", LONGMAN_DICTIONARY_DOMAIN, word))

	if parseHtmlErr != nil {
		return nil, parseHtmlErr
	}

	return wordMeangins, nil
}

package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const (
	kospiURL        = "https://finance.naver.com/sise/sise_market_sum.naver?sosok=0&page=%d"
	kosdaqURL       = "https://finance.naver.com/sise/sise_market_sum.naver?sosok=1&page=%d"
	etfAPIURL       = "https://finance.naver.com/api/sise/etfItemList.nhn?etfType=0&targetColumn=market_sum&sortOrder=desc"
	defaultMaxPages = 40 // Fallback if detection fails
)

// Scraper handles fetching the full list of tickers from Naver Finance.
type Scraper struct {
	client *http.Client
	re     *regexp.Regexp
	pgRe   *regexp.Regexp
}

func NewScraper() *Scraper {
	return &Scraper{
		client: &http.Client{Timeout: 10 * time.Second},
		re:     regexp.MustCompile(`href="/item/main.naver\?code=([A-Z0-9]+)" class="tltle">([^<]+)</a>`),
		pgRe:   regexp.MustCompile(`class="pgRR">\s*<a href=".*?page=(\d+)`),
	}
}

type etfResponse struct {
	Result struct {
		EtfItemList []struct {
			ItemCode string `json:"itemcode"`
			ItemName string `json:"itemname"`
		} `json:"etfItemList"`
	} `json:"result"`
}

func (s *Scraper) ScrapeAll() ([]Ticker, error) {
	var (
		tickers []Ticker
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	scrapeMarket := func(urlFmt, marketName string) {
		defer wg.Done()

		firstURL := fmt.Sprintf(urlFmt, 1)
		resp, err := s.client.Get(firstURL)
		if err != nil {
			return
		}
		rawBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return
		}

		decodedBody, err := DecodeEUCKR(rawBody)
		if err != nil {
			return
		}
		body := []byte(decodedBody)

		matches := s.re.FindAllSubmatch(body, -1)
		mu.Lock()
		for _, m := range matches {
			tickers = append(tickers, Ticker{
				Code:   string(m[1]),
				Name:   string(m[2]),
				Market: marketName,
			})
		}
		mu.Unlock()

		lastPage := defaultMaxPages
		if pgMatch := s.pgRe.FindSubmatch(body); len(pgMatch) > 1 {
			if lp, err := strconv.Atoi(string(pgMatch[1])); err == nil {
				lastPage = lp
			}
		}

		for page := 2; page <= lastPage; page++ {
			url := fmt.Sprintf(urlFmt, page)
			time.Sleep(50 * time.Millisecond)

			resp, err := s.client.Get(url)
			if err != nil {
				break
			}

			rawBody, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				break
			}

			decodedBody, err := DecodeEUCKR(rawBody)
			if err != nil {
				break
			}
			body := []byte(decodedBody)

			matches := s.re.FindAllSubmatch(body, -1)
			if len(matches) == 0 {
				break
			}

			mu.Lock()
			for _, m := range matches {
				tickers = append(tickers, Ticker{
					Code:   string(m[1]),
					Name:   string(m[2]),
					Market: marketName,
				})
			}
			mu.Unlock()
		}
	}

	scrapeETF := func() {
		defer wg.Done()

		resp, err := s.client.Get(etfAPIURL)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		rawBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		decodedBody, err := DecodeEUCKR(rawBody)
		if err != nil {
			return
		}

		var result etfResponse
		if err := json.Unmarshal([]byte(decodedBody), &result); err != nil {
			return
		}

		mu.Lock()
		for _, item := range result.Result.EtfItemList {
			tickers = append(tickers, Ticker{
				Code:   item.ItemCode,
				Name:   item.ItemName,
				Market: "KOSPI", // Most ETFs are listed on KOSPI
			})
		}
		mu.Unlock()
	}

	wg.Add(3)
	go scrapeMarket(kospiURL, "KOSPI")
	go scrapeMarket(kosdaqURL, "KOSDAQ")
	go scrapeETF()

	wg.Wait()

	if len(tickers) == 0 {
		return nil, fmt.Errorf("scraped 0 tickers; network or parsing error likely")
	}

	unique := make([]Ticker, 0, len(tickers))
	seen := make(map[string]bool)
	for _, t := range tickers {
		if !seen[t.Code] {
			seen[t.Code] = true
			unique = append(unique, t)
		}
	}

	return unique, nil
}

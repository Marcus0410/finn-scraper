package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"sort"
	"strings"
)

type Job struct {
	title, publishedDate, url, company, deadline string
}

const visitURL = "https://www.finn.no/job/fulltime/search.html?location=1.20001.20061&occupation=0.23&q=nyutdannet"

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.finn.no"),
	)

	deadlineCollector := c.Clone()

	// called before an HTTP request is triggered
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error while scraping:", err.Error())
	})

	jobs := []Job{}

	c.OnHTML("div.flex.flex-col",
		func(h *colly.HTMLElement) {
			job := Job{}
			job.title = h.ChildText("h2")
			job.company = h.DOM.Find("div.flex.flex-col.text-xs span").First().Text()
			job.publishedDate = strings.Split(h.ChildText(".s-text-subtle"), "|")[0]
			job.url = h.ChildAttr("h2 a", "href")

			if job.title != "" && job.url != "" {
				jobs = append(jobs, job)
				deadlineCollector.Visit(job.url)
			}
		})

	// scrape deadline
	deadlineCollector.OnHTML("li.flex.flex-col", func(h *colly.HTMLElement) {
		// if deadline has not been scraped
		if jobs[len(jobs)-1].deadline == "" {
			jobs[len(jobs)-1].deadline = h.DOM.Find("span.font-bold").First().Text()
		}
	})

	c.Visit(visitURL)

	// sort jobs based on publishedDate
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].deadline < jobs[j].deadline
	})

	fmt.Printf("Found %v jobs\n", len(jobs))
	for _, job := range jobs {
		fmt.Println("-----------------------")
		fmt.Println("Title:", job.title)
		fmt.Println("Company:", job.company)
		fmt.Println("Published:", job.publishedDate)
		fmt.Println("URL:", job.url)
		fmt.Println("Deadline:", job.deadline)
		fmt.Println("-----------------------")
		fmt.Println()
	}
}

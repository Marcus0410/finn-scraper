package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/gocolly/colly"
)

type Job struct {
	Title, PublishedDate, Url, Company, Deadline string
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
			job.Title = h.ChildText("h2")
			job.Company = h.DOM.Find("div.flex.flex-col.text-xs span").First().Text()
			job.PublishedDate = strings.Split(h.ChildText(".s-text-subtle"), "|")[0]
			job.Url = h.ChildAttr("h2 a", "href")

			if job.Title != "" && job.Url != "" {
				jobs = append(jobs, job)
				deadlineCollector.Visit(job.Url)
			}
		})

	// scrape deadline
	deadlineCollector.OnHTML("li.flex.flex-col", func(h *colly.HTMLElement) {
		// if deadline has not been scraped
		if jobs[len(jobs)-1].Deadline == "" {
			jobs[len(jobs)-1].Deadline = h.DOM.Find("span.font-bold").First().Text()
		}
	})

	c.Visit(visitURL)

	// sort jobs based on publishedDate
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Deadline < jobs[j].Deadline
	})

	fmt.Printf("Found %v jobs\n", len(jobs))
	for _, job := range jobs {
		fmt.Println("-----------------------")
		fmt.Println("Title:", job.Title)
		fmt.Println("Company:", job.Company)
		fmt.Println("Published:", job.PublishedDate)
		fmt.Println("URL:", job.Url)
		fmt.Println("Deadline:", job.Deadline)
		fmt.Println("-----------------------")
		fmt.Println()
	}

	err := saveJobs(jobs)
	if err != nil {
		log.Fatal(err)
	}
}

// save jobs to json file
func saveJobs(jobs []Job) error {
	file, err := os.Create("jobs.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(jobs)
}

// load previously stored jobs from json file
func loadJobs() ([]Job, error) {
	var jobs []Job
	file, err := os.Open("jobs.json")
	if err != nil {
		return jobs, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&jobs)

	return jobs, err
}

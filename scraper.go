package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"sort"
	"strings"
)

type Job struct {
	title, deadline, url, company string
}

const visitURL = "https://www.finn.no/job/fulltime/search.html?location=1.20001.20061&occupation=0.23&q=nyutdannet"

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.finn.no"),
	)
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
			// scraping logic
			selection := h.DOM
			// fmt.Println(selection.Find("div.s-text-subtle").Text())
			companyDiv := selection.Find("div.flex.flex-col.text-xs")

			job := Job{}
			job.title = h.ChildText("h2")
			job.company = companyDiv.Text()
			job.deadline = strings.Split(h.ChildText(".s-text-subtle"), "|")[0]
			job.url = ""

			if job.title != "" {
				jobs = append(jobs, job)
			}
		})

	c.Visit(visitURL)

	// sort jobs based on deadline
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].deadline < jobs[j].deadline
	})

	fmt.Printf("Found %v jobs\n", len(jobs))
	for _, job := range jobs {
		fmt.Println("-----------------------")
		fmt.Println("Title:", job.title)
		fmt.Println("Company:", job.company)
		fmt.Println("Deadline:", job.deadline)
		fmt.Println("URL:", job.url)
		fmt.Println("-----------------------")
		fmt.Println()
	}
}

func (j Job) String() string {
	return j.title + " " + j.company
}

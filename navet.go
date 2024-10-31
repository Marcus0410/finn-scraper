package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func scrapeNavet() []Job {
	c := colly.NewCollector(colly.AllowedDomains("ifinavet.no"))

	detailCollector := c.Clone()

	// called before an HTTP request is triggered
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error while scraping:", err.Error())
	})

	jobs := []Job{}

	c.OnHTML("div.job-card-bottom", func(h *colly.HTMLElement) {
		var job Job
		job.Title = h.ChildText("h3.job-title")
		job.Url = "https://ifinavet.no" + h.ChildAttr("a", "href")

		// get the rest of the details from url
		jobs = append(jobs, job)
		detailCollector.Visit(job.Url)
	})

	// get deadline
	detailCollector.OnHTML("div.event-meta.job-meta:nth-of-type(1) span:nth-of-type(3)", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Deadline = h.Text
	})

	// get company name
	detailCollector.OnHTML("div.company-info", func(h *colly.HTMLElement) {
		jobs[len(jobs)-1].Company = h.ChildText("h2")
		fmt.Println(jobs[len(jobs)-1].Company)
	})

	c.Visit("https://ifinavet.no/stillingsannonser/")

	return jobs
}

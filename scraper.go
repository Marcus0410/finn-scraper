package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"strconv"
	"strings"
)

type Job struct {
	Title, PublishedDate, Url, Company, Deadline string
	Id                                           int
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

	// jobsMap, err := loadJobs()
	jobsMap := make(map[int]Job)

	var lastAddedJobId int
	var err error

	c.OnHTML("div.flex.flex-col",
		func(h *colly.HTMLElement) {
			job := Job{}
			job.Title = h.ChildText("h2")
			job.Company = h.DOM.Find("div.flex.flex-col.text-xs span").First().Text()
			job.PublishedDate = strings.Split(h.ChildText(".s-text-subtle"), "|")[0]
			job.Url = h.ChildAttr("h2 a", "href")

			if job.Title != "" && job.Url != "" {
				// get uniqe id from url
				job.Id, err = strconv.Atoi(strings.Split(job.Url, "=")[1])
				if err != nil {
					log.Fatal(err)
				}

				jobsMap[job.Id] = job
				lastAddedJobId = job.Id

				deadlineCollector.Visit(job.Url)
			}
		})

	// scrape deadline
	deadlineCollector.OnHTML("li.flex.flex-col", func(h *colly.HTMLElement) {
		// if deadline has not been scraped
		if jobsMap[lastAddedJobId].Deadline == "" {
			job := jobsMap[lastAddedJobId]
			job.Deadline = h.DOM.Find("span.font-bold").First().Text()
			jobsMap[lastAddedJobId] = job
		}
	})

	// start scraping jobs
	c.Visit(visitURL)

	// send listings via smtp
	err = sendSmtp(jobsMap)
	if err == nil {
		fmt.Println("Successfully sent email!")
	} else {
		log.Fatal(err)
	}

	// save jobs to file
	saveJobs(jobsMap)
}

func sendSmtp(jobsMap map[int]Job) error {
	var sb strings.Builder

	sb.WriteString("Fant " + strconv.Itoa(len(jobsMap)) + " jobbannonser.\n\n")
	for _, job := range jobsMap {
		sb.WriteString("Tittel: " + job.Title + "\n")
		sb.WriteString("Bedrift: " + job.Company + "\n")
		sb.WriteString("Søknadsfrist: " + job.Deadline + "\n")
		sb.WriteString("URL: " + job.Url + "\n\n")
	}

	return sendNewListings(sb.String())
}

// save jobs to json file
func saveJobs(jobs map[int]Job) error {
	file, err := os.Create("jobs.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(jobs)
}

// load previously stored jobs from json file
func loadJobs() (map[int]Job, error) {
	jobs := make(map[int]Job)
	file, err := os.Open("jobs.json")
	if err != nil {
		return jobs, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&jobs)

	return jobs, err
}

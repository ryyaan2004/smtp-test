package main

import (
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"time"
)

var recipient = ""

func main() {
	var toSend int64
	var workers int
	var server, from, to string
	flag.Int64Var(&toSend, "count", 10, "The number of messages to send, default 10, max is 65,000")
	flag.IntVar(&workers, "workers", 10, "The number of worker go routines")
	flag.StringVar(&recipient, "recipient", "", "the recipient to use")
	flag.StringVar(&server, "server", "", "The mail server and port, like mail.place.com:587")
	flag.StringVar(&from, "from", "", "The from address, e.g. tester@place.com")
	flag.StringVar(&to, "to", "", "Email address to send test messages to")

	flag.Parse()
	if toSend > 65_000 {
		log.Fatalln("please use a value less than 65,000")
	}
	if recipient == "" {
		log.Fatalln("please provide a recipient")
	}

	senders := make(chan *WorkItem, toSend)
	results := make(chan string)
	var sendResults []string

	// set up simple load test
	for i := 0; i < workers; i++ {
		c, err := smtp.Dial(server)
		if err != nil {
			log.Printf("error connecting to %s %s", server, err.Error())
		}
		go worker(c, senders, results)
	}

	go func() {
		for i := int64(1); i <= toSend; i++ {
			senders <- &WorkItem{
				To:   to,
				Msg:  fmt.Sprintf("Test message %d of %d enqueued at %v", i, toSend, time.Now()),
				From: "",
			}
		}
	}()

	for i := int64(0); i < toSend; i++ {
		result := <-results
		sendResults = append(sendResults, result)
	}
	close(senders)
	close(results)
	for _, m := range sendResults {
		fmt.Println(m)
	}
}

type WorkItem struct {
	From string
	Msg  string
	To   string
}

func worker(c *smtp.Client, workItems chan *WorkItem, results chan string) {
	defer c.Quit()
	for s := range workItems {
		// all below errors are likely to be rate limiting or other temporary errors (e.g. 421)
		if err := c.Mail(s.From); err != nil {
			log.Printf("unable to set the sender: %s\n", err.Error())
		}
		if err := c.Rcpt(s.To); err != nil {
			log.Printf("unable to set the receiver: %s\n", err.Error())
		}

		wc, err := c.Data()
		if err != nil {
			log.Printf("unable to initiate the message body: %s\n", err.Error())
		}

		msg := fmt.Sprintf("%s\nsent at %v\n", s.Msg, time.Now())
		_, err = fmt.Fprint(wc, msg)
		err = wc.Close()
		results <- fmt.Sprintf("To:%s, From:%s, Msg:%s", s.To, s.From, msg)
	}
}

package cron

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/robfig/cron"
)

type EntryID int

type Execute func(EntryID)

type Scheduler struct {
	stop    chan struct{}
	add     chan *Entry
	update  chan *Entry
	remove  chan EntryID
	entries []*Entry
	logger  *log.Logger
}

func New(logger *log.Logger) *Scheduler {
	return &Scheduler{
		stop:    make(chan struct{}),
		add:     make(chan *Entry, 10),
		update:  make(chan *Entry, 10),
		remove:  make(chan EntryID, 10),
		entries: []*Entry{},
		logger:  logger,
	}
}

type Entry struct {
	id       EntryID
	schedule cron.Schedule
	next     time.Time
	execute  Execute
}

type byTime []*Entry

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	if s[i].next.IsZero() {
		return false
	}
	if s[j].next.IsZero() {
		return true
	}
	return s[i].next.Before(s[j].next)
}

func (c *Scheduler) ValidateSpec(spec string) error {
	_, err := cron.Parse(spec)
	return err
}

func (c *Scheduler) AddEntry(id EntryID, spec string, execute Execute) error {
	schedule, err := cron.Parse(spec)
	if err != nil {
		return err
	}

	now := time.Now()
	next := schedule.Next(now)

	c.add <- &Entry{
		id:       id,
		schedule: schedule,
		next:     next,
		execute:  execute,
	}

	return nil
}

func (c *Scheduler) Start() {
	c.logger.Println("Started cron scheduler")
	var wg sync.WaitGroup
Loop:
	for {
		nextCh := make(<-chan time.Time)
		if len(c.entries) > 0 {
			c.checker(&wg)
			sort.Sort(byTime(c.entries))
			durationTillNext := time.Until(c.entries[0].next)
			nextCh = time.After(durationTillNext)
		}

		select {
		case <-c.stop:
			break Loop
		case e := <-c.add:
			c.updateOrAddEntry(e)
		case id := <-c.remove:
			c.removeEntryByID(id)
		case <-nextCh:
		}
	}

	c.logger.Println("Scheduler waiting for jobs to finish...")
	wg.Wait()
	c.logger.Println("Stopped scheduler")
}

func (c *Scheduler) updateOrAddEntry(entry *Entry) {
	for i, e := range c.entries {
		if entry.id == e.id {
			c.entries[i] = entry
			return
		}
	}
	c.entries = append(c.entries, entry)
}

func (c *Scheduler) Stop() {
	c.logger.Println("Stopping scheduler...")
	c.stop <- struct{}{}
}

func (c *Scheduler) checker(wg *sync.WaitGroup) {
	now := time.Now()
	for _, e := range c.entries {
		entry := e
		if entry.next.After(now) || entry.next.IsZero() {
			continue
		}
		go func() {
			wg.Add(1)
			defer wg.Done()
			c.execute(entry)
		}()
		entry.next = entry.schedule.Next(now)
	}
}

func (c *Scheduler) execute(e *Entry) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Printf("Entry with id of %d failed due to: %v", e.id, r)
		}
	}()
	e.execute(e.id)
}

func (c *Scheduler) RemoveEntry(id EntryID) {
	c.remove <- id
}

func (c *Scheduler) removeEntryByID(id EntryID) {
	found := false
	foundID := 0

	for i, entry := range c.entries {
		if entry.id == id {
			found = true
			foundID = i
			break
		}
	}

	if found {
		c.entries = append(c.entries[:foundID], c.entries[foundID+1:]...)
	}
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/GitProger/go-telecom-2025/internal/config"
	"github.com/GitProger/go-telecom-2025/internal/model"
	"github.com/GitProger/go-telecom-2025/internal/monitor"
	"github.com/GitProger/go-telecom-2025/internal/provider"
)

func printNonNil(s any) {
	if !reflect.ValueOf(s).IsNil() {
		fmt.Println(s)
	}
}

var _ = mainSimple

func mainSimple() { // file only
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <config_file> <event_file>\n", os.Args[0])
	}

	config, err := config.LoadConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	events, err := provider.ScanFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	m := monitor.NewEventMonitor(config)

	for _, event := range events {
		out, err := m.DigestEvent(event)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(event)
		printNonNil(out)
	}

	for _, e := range m.Disqualified() {
		printNonNil(e)
	}

	fmt.Println("### Resulting Report ###")
	for _, c := range m.GetReport() {
		fmt.Println(c)
	}
}

func main() { // interactive
	// mainSimple()
	// return

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <config_file> [event_file]\n", os.Args[0])
	}
	configFile := os.Args[1]
	config, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	m := monitor.NewEventMonitor(config)
	digestLog := func(event *model.Event) {
		out, err := m.DigestEvent(event)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(event)
		printNonNil(out)
	}

	var source io.Reader
	if len(os.Args) == 3 {
		eventFile := os.Args[2]
		f, err := os.Open(eventFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		source = f
	} else { // EOF catches on Ctrl+D
		source = os.Stdin
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, syscall.SIGINT, syscall.SIGTERM)

	events, errs := provider.Scan(ctx, source)

rwLoop:
	for {
		select {
		case <-ctrlC:
			cancel()
			break rwLoop
		case event, ok := <-events:
			if !ok {
				events = nil // remove chan from select-case
			} else if event != nil {
				digestLog(event)
			}
		case err, ok := <-errs:
			if ok && err != nil {
				log.Fatalf("error during scan: %v", err)
			}
			errs = nil
		}
		if events == nil && errs == nil {
			break
		}
	}

	for _, e := range m.Disqualified() {
		printNonNil(e)
	}

	fmt.Println("### Resulting Report ###")
	for _, c := range m.GetReport() {
		fmt.Println(c)
	}
}

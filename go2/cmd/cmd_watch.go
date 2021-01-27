package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	p "path"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/zaydek/retro/loggers"
	"github.com/zaydek/retro/sse"
	"github.com/zaydek/retro/watcher"
)

func (r *Runtime) esbuildBuild() {
	results := api.Build(api.BuildOptions{
		Bundle:      true,
		EntryPoints: []string{path.Join(r.WatchCommand.Directory, "index.js")},
		Incremental: true,
		Outfile:     path.Join(r.Config.BuildDirectory, "app.js"),
		Write:       true,
	})
	r.esbuildResult = results
	r.esbuildWarnings = results.Warnings
	r.esbuildErrors = results.Errors
}

func (r *Runtime) esbuildRebuild() {
	// start := time.Now()
	results := r.esbuildResult.Rebuild()
	r.esbuildResult = results
	r.esbuildWarnings = results.Warnings
	r.esbuildErrors = results.Errors
	// fmt.Printf("⚡️ %0.3fs\n", time.Since(start).Seconds())
}

func (r Runtime) Watch() {
	serverSentEvents := make(chan sse.Event, 8)

	must(copyAssetDirectoryToBuildDirectory(r.Config))
	r.esbuildBuild()

	fmt.Printf("👾 http://localhost:%d\n", r.getPort())

	if len(r.esbuildWarnings) > 0 {
		loggers.Stderr.Println(formatEsbuildMessagesAsTermString(r.esbuildWarnings))
	}
	if len(r.esbuildErrors) > 0 {
		loggers.Stderr.Println(formatEsbuildMessagesAsTermString(r.esbuildErrors))
	}

	go func() {
		for range watcher.New(r.WatchCommand.Directory, r.WatchCommand.Poll) {
			r.esbuildRebuild()
			serverSentEvents <- sse.Event{Event: "reload"}
		}
	}()

	http.HandleFunc("/", func(wr http.ResponseWriter, req *http.Request) {
		if ext := path.Ext(req.URL.Path); ext == "" {
			r.esbuildRebuild()
			if len(r.esbuildWarnings) > 0 {
				// TODO
				loggers.Stderr.Println(formatEsbuildMessagesAsTermString(r.esbuildWarnings))
				data, _ := json.Marshal(formatEsbuildMessagesAsTermString(r.esbuildWarnings))
				defer func() {
					// Pause 100ms so the server-sent event does not drop on refresh:
					time.Sleep(100 * time.Millisecond)
					serverSentEvents <- sse.Event{Event: "warning", Data: string(data)}
				}()
			}
			if len(r.esbuildErrors) > 0 {
				loggers.Stderr.Println(formatEsbuildMessagesAsTermString(r.esbuildErrors))
				fmt.Fprintln(wr, esbuildMessagesAsHTMLDocument(r.esbuildErrors))
				return
			}
		}
		http.ServeFile(wr, req, p.Join(string(r.Config.BuildDirectory), req.URL.Path))
	})

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		flusher, ok := w.(http.Flusher)
		if !ok {
			loggers.Stderr.Println("Your browser does not support server-sent events (SSE).")
			os.Exit(1)
			return
		}
		for {
			select {
			case e := <-serverSentEvents:
				e.Write(w)
				flusher.Flush()
			case <-r.Context().Done():
				// No-op
				return
			}
		}
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", r.WatchCommand.Port), nil); err != nil {
		loggers.Stderr.Println(err)
		os.Exit(1)
	}
}

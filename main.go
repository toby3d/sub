//go:generate go install github.com/valyala/quicktemplate/qtc@latest
//go:generate qtc -dir=web
//go:generate go install golang.org/x/text/cmd/gotext@latest
//go:generate gotext -srclang=en update -out=catalog_gen.go -lang=en,ru
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)

var (
	cpuProfilePath, memProfilePath string
	enablePprof                    bool
)

func init() {
	flag.BoolVar(&enablePprof, "pprof", false, "enable pprof mode")
	flag.StringVar(&cpuProfilePath, "cpuprofile", "", "set path to saveing CPU memory profile")
	flag.StringVar(&memProfilePath, "memprofile", "", "set path to saveing pprof memory profile")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	server := http.Server{
		Addr:     ":3000",
		ErrorLog: logger,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, World!")
		}),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if cpuProfilePath != "" {
		cpuProfile, err := os.Create(cpuProfilePath)
		if err != nil {
			logger.Fatalln("could not create CPU profile:", err)
		}
		defer cpuProfile.Close()

		if err = pprof.StartCPUProfile(cpuProfile); err != nil {
			logger.Fatalln("could not start CPU profile:", err)
		}
		defer pprof.StopCPUProfile()
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalln("cannot listen and serve:", err)
		}
	}()

	<-done

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalln("failed shutdown of server:", err)
	}

	if memProfilePath == "" {
		return
	}

	memProfile, err := os.Create(memProfilePath)
	if err != nil {
		logger.Fatalln("could not create memory profile:", err)
	}
	defer memProfile.Close()

	runtime.GC()
	if err = pprof.WriteHeapProfile(memProfile); err != nil {
		logger.Fatalln("could not write memory profile:", err)
	}
}

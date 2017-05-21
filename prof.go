package yet_another_goprof

import (
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"time"
)

func ProfileCpu(cpuProfile string) {
	f, err := os.Create(cpuProfile)
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		pprof.StopCPUProfile()
		os.Exit(0)
	}()
}

func ProfileMemory(heapProfile string) {
	m := &runtime.MemStats{}
	tick1m := time.Tick(1 * time.Minute)
	tick5s := time.Tick(5 * time.Second)

	for {
		select {
		case <-tick5s:
			runtime.ReadMemStats(m)
			log.Println("")
			log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			log.Printf("~ Goroutines count %d\n", runtime.NumGoroutine())
			log.Printf("~ Alloc %dKB\n", m.Alloc/1024)
			log.Printf("~ TotalAlloc %dKB\n", m.TotalAlloc/1024)
			log.Printf("~ Sys (sum of XxxSys below) %dKB\n", m.Sys/1024)
			log.Printf("~ Lookups (number of pointer lookups) %d\n", m.Lookups)
			log.Printf("~ Mallocs %d\n", m.Mallocs)
			log.Printf("~ Frees %d\n", m.Frees)

			log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

			log.Printf("~ HeapAlloc %dKB\n", m.HeapAlloc/1024)
			log.Printf("~ HeapSys %dKB\n", m.HeapSys/1024)
			log.Printf("~ HeapIdle %dKB\n", m.HeapIdle/1024)
			log.Printf("~ HeapInuse %dKB\n", m.HeapInuse/1024)
			log.Printf("~ HeapReleased %dKB\n", m.HeapReleased/1024)
			log.Printf("~ HeapObjects %d\n", m.HeapObjects)

			log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

			log.Printf("~ NextGC %d\n", m.NextGC)
			log.Printf("~ LastGC %v\n", time.Unix(0, int64(m.LastGC)))
			log.Printf("~ PauseTotalNs %d\n", m.PauseTotalNs)
			log.Printf("~ NumGC %d\n", m.PauseTotalNs)
			log.Printf("~ GCCPUFraction %f\n", m.GCCPUFraction)
			log.Printf("~ EnableGC %v\n", m.EnableGC)
			log.Printf("~ DebugGC %v\n", m.DebugGC)
			log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			log.Println("")

		case <-tick1m:
			var fHeapProfiling io.Writer
			fHeapProfiling, _ = os.Create(heapProfile)
			pprof.WriteHeapProfile(fHeapProfiling)
			log.Println("~ head saved")
		}
	}
}

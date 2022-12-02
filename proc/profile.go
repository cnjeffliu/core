//go:build linux || darwin
// +build linux darwin

/*
 * @Author: Jeffrey Liu
 * @Date: 2022-11-29 13:54:49
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-12-02 20:58:14
 * @Description:
 */

package proc

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/cnjeffliu/gocore/timex"
)

// DefaultMemProfileRate is the default memory profiling rate.
// See also http://golang.org/pkg/runtime/#pkg-variables
const DefaultMemProfileRate = 4096

// started is non zero if a profile is running.
var started uint32

// Profile represents an active profiling session.
type Profile struct {
	// closers holds cleanup functions that run after each profile
	closers []func()

	// stopped records if a call to profile.Stop has been made
	stopped uint32

	// dumpdir records the profile saved in
	dumpdir string
}

func (p *Profile) close() {
	for _, closer := range p.closers {
		closer()
	}
}

func (p *Profile) startBlockProfile() {
	fn := createDumpFile(p.dumpdir, "block")
	f, err := os.Create(fn)
	if err != nil {
		fmt.Printf("profile: could not create block profile %q: %v", fn, err)
		return
	}

	runtime.SetBlockProfileRate(1)
	fmt.Printf("profile: block profiling enabled, %s", fn)
	p.closers = append(p.closers, func() {
		pprof.Lookup("block").WriteTo(f, 0)
		f.Close()
		runtime.SetBlockProfileRate(0)
		fmt.Printf("profile: block profiling disabled, %s", fn)
	})
}

func (p *Profile) startCpuProfile() {
	fn := createDumpFile(p.dumpdir, "cpu")
	f, err := os.Create(fn)
	if err != nil {
		fmt.Printf("profile: could not create cpu profile %q: %v", fn, err)
		return
	}

	fmt.Printf("profile: cpu profiling enabled, %s", fn)
	pprof.StartCPUProfile(f)
	p.closers = append(p.closers, func() {
		pprof.StopCPUProfile()
		f.Close()
		fmt.Printf("profile: cpu profiling disabled, %s", fn)
	})
}

func (p *Profile) startMemProfile() {
	fn := createDumpFile(p.dumpdir, "mem")
	f, err := os.Create(fn)
	if err != nil {
		fmt.Printf("profile: could not create memory profile %q: %v", fn, err)
		return
	}

	old := runtime.MemProfileRate
	runtime.MemProfileRate = DefaultMemProfileRate
	fmt.Printf("profile: memory profiling enabled (rate %d), %s", runtime.MemProfileRate, fn)
	p.closers = append(p.closers, func() {
		pprof.Lookup("heap").WriteTo(f, 0)
		f.Close()
		runtime.MemProfileRate = old
		fmt.Printf("profile: memory profiling disabled, %s", fn)
	})
}

func (p *Profile) startMutexProfile() {
	fn := createDumpFile(p.dumpdir, "mutex")
	f, err := os.Create(fn)
	if err != nil {
		fmt.Printf("profile: could not create mutex profile %q: %v", fn, err)
		return
	}

	runtime.SetMutexProfileFraction(1)
	fmt.Printf("profile: mutex profiling enabled, %s", fn)
	p.closers = append(p.closers, func() {
		if mp := pprof.Lookup("mutex"); mp != nil {
			mp.WriteTo(f, 0)
		}
		f.Close()
		runtime.SetMutexProfileFraction(0)
		fmt.Printf("profile: mutex profiling disabled, %s", fn)
	})
}

func (p *Profile) startThreadCreateProfile() {
	fn := createDumpFile(p.dumpdir, "threadcreate")
	f, err := os.Create(fn)
	if err != nil {
		fmt.Printf("profile: could not create threadcreate profile %q: %v", fn, err)
		return
	}

	fmt.Printf("profile: threadcreate profiling enabled, %s", fn)
	p.closers = append(p.closers, func() {
		if mp := pprof.Lookup("threadcreate"); mp != nil {
			mp.WriteTo(f, 0)
		}
		f.Close()
		fmt.Printf("profile: threadcreate profiling disabled, %s", fn)
	})
}

func (p *Profile) startTraceProfile() {
	fn := createDumpFile(p.dumpdir, "trace")
	f, err := os.Create(fn)
	if err != nil {
		fmt.Printf("profile: could not create trace output file %q: %v", fn, err)
		return
	}

	if err := trace.Start(f); err != nil {
		fmt.Printf("profile: could not start trace: %v", err)
		return
	}

	fmt.Printf("profile: trace enabled, %s", fn)
	p.closers = append(p.closers, func() {
		trace.Stop()
		fmt.Printf("profile: trace disabled, %s", fn)
	})
}

// Stop stops the profile and flushes any unwritten data.
func (p *Profile) Stop() {
	if !atomic.CompareAndSwapUint32(&p.stopped, 0, 1) {
		// someone has already called close
		return
	}
	p.close()
	atomic.StoreUint32(&started, 0)
	fmt.Printf("profile: dump profile file in %s", p.dumpdir)
}

// StartProfile starts a new profiling session.
// The caller should call the Stop method on the value returned
// to cleanly stop profiling.
//
//	go func() {
//		var profiler proc.Stopper
//		signals := make(chan os.Signal, 1)
//		signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)
//		for {
//			v := <-signals
//			switch v {
//			case syscall.SIGUSR2:
//				if profiler == nil {
//					profiler = proc.StartProfile()
//				} else {
//					profiler.Stop()
//					profiler = nil
//				}
//			}
//		}
//	}()
func StartProfile(dumpDir ...string) Stopper {
	if !atomic.CompareAndSwapUint32(&started, 0, 1) {
		fmt.Println("profile: Start() already called")
		return noopStopper
	}

	var prof Profile
	prof.dumpdir = os.TempDir()
	if len(dumpDir) > 0 && len(dumpDir[0]) > 0 {
		prof.dumpdir = dumpDir[0]
	}

	prof.startCpuProfile()
	prof.startMemProfile()
	prof.startMutexProfile()
	prof.startBlockProfile()
	prof.startTraceProfile()
	prof.startThreadCreateProfile()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		<-c

		fmt.Println("profile: caught interrupt, stopping profiles")
		prof.Stop()

		signal.Reset()
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	}()

	return &prof
}

func createDumpFile(saveDir string, kind string) string {
	command := path.Base(os.Args[0])
	pid := syscall.Getpid()
	return path.Join(saveDir, fmt.Sprintf("%s-%d-%s-%s.pprof",
		command, pid, kind, time.Now().Format(timex.TIME_LAYOUT_COMPACT_SECOND)))
}

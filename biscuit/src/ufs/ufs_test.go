package ufs

import "testing"
import "fmt"
import "strconv"
import "time"
import "ustr"
//import "runtime/pprof"

const (
	nlogblks   = 32
	ninodeblks = 1
	ndatablks  = 20
)

var STOP = false

//
// Util
//
func uniqdir(id int) string {
	return "d" + strconv.Itoa(id)
}

func doTestSimple(tfs *Ufs_t, d ustr.Ustr) string {
	ub := mkData(1, 11)
	e := tfs.MkFile(d.ExtendStr("f1"), ub)
	if e != 0 {
		return fmt.Sprintf("mkFile %v failed", "f1")
	}
	e = tfs.MkFile(d.ExtendStr("f2"), ub)
	if e != 0 {
		return fmt.Sprintf("mkFile %v failed", "f2")
	}
	e = tfs.Unlink(d.ExtendStr("f2"))
	if e != 0 {
		return fmt.Sprintf("unlink %v failed", "f2")
	}
	e = tfs.Unlink(d.ExtendStr("f1"))
	if e != 0 {
		return fmt.Sprintf("unlink %v failed", "f1")
	}
	return ""
}

func concurrent(t *testing.T) {
	n := 2
	dst := "tmp.img"
	MkDisk(dst, nil, nlogblks, ninodeblks*2, ndatablks*10)

	c := make(chan int)
	tfs := BootMemFS(dst)
	start := time.Now()
	stop := start
	for i := 0; i < n; i++ {
		go func(id int) {
			iter := 0
			d := ustr.Ustr(uniqdir(id))
			tfs.MkDir(d)
			for !STOP {
				doTestSimple(tfs, d)
				tfs.Sync()
				iter++
			}
			c <- iter
		}(i)
	}
	go func() {
		time.Sleep(2*time.Second)
		STOP = true
		stop = time.Now()
		fmt.Printf("Timer Thread Done\n")
		c <- 0
	}()
	s := 0
	for i := 0; i < n+1; i++ {
		s += <-c
		fmt.Printf("Got %d tests from %d threads\n", s, i)
	}
	fmt.Printf("Did %d tests in %v seconds\n", s, stop.Sub(start))
	fmt.Printf("Created/Wrote/Closed %f files/sec\n", float64(s*2)/float64(stop.Sub(start).Seconds()))
	ShutdownFS(tfs)
}

func TestFSConcurNotSame(t *testing.T) {
	/*if *flagCpuprofile != "" {
	    f, err := os.Create(*flagCpuprofile)
	    if err != nil {
	        log.Fatal(err)
	    }
	    pprof.StartCPUProfile(f)
	    defer pprof.StopCPUProfile()
	}*/
	concurrent(t)
}

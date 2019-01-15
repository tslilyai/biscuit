package ufs

import "testing"
import "fmt"
//import "io"
//import "os"
import "strconv"
//import "sync"
import "time"

//import "bpath"
//import "defs"
//import "fd"
//import "fs"
//import "mem"
import "ustr"

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

func doTestSimple(tfs *Ufs_t, d ustr.Ustr) {
	ub := mkData(1, 11)
	tfs.MkDir(d)
	tfs.MkFile(d.ExtendStr("f1"), ub)
	tfs.MkFile(d.ExtendStr("f2"), ub)
	//tfs.MkDir(d.ExtendStr("d0"))
	//tfs.MkDir(d.ExtendStr("d0/d1"))
	//tfs.Append(d.ExtendStr("f1"), ub)
	tfs.Unlink(d.ExtendStr("f2"))
	tfs.Unlink(d.ExtendStr("f1"))
	//tfs.Unlink(d.ExtendStr("d0/d1"))
	//tfs.Unlink(d.ExtendStr("d0"))
	tfs.Unlink(d)
}

func concurrent(t *testing.T) {
	n := 1
	dst := "tmp.img"
	MkDisk(dst, nil, nlogblks, ninodeblks*2, ndatablks*10)

	c := make(chan int)
	tfs := BootMemFS(dst)
	start := time.Now()
	stop := start
	for i := 0; i < n; i++ {
		go func(id int) {
			iter := 0
			for !STOP {
				d := ustr.Ustr(uniqdir(id+(iter*n)))
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
		fmt.Printf("Timer Thread Done")
		c <- 0
	}()
	s := 0
	for i := 0; i < n+1; i++ {
		s += <-c
		fmt.Printf("Got %d tests from %d threads", s, i)
	}
	fmt.Printf("Did %d tests in %v seconds", s, stop.Sub(start));
	ShutdownFS(tfs)
}

func TestFSConcurNotSame(t *testing.T) {
	concurrent(t)
}

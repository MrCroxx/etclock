package dao

import (
	"flag"
	"testing"
	"time"
)

func TestDao(t *testing.T) {

	if !flag.Parsed() {
		flag.Parse()
	}

	args := flag.Args()

	if len(args) == 0 {
		t.Fatalf("missing etcd host")
	}

	host := args[0]

	d, err := NewDao(Config{
		Endpoints: []string{
			host + ":2301",
			host + ":2302",
			host + ":2303",
		},
		Timeout: time.Second * 3,
	})
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	r, err := d.RLock("res1", "owner1")
	if err != nil {
		t.Errorf("%s", err)
	}
	t.Logf("%t\n", r)
}

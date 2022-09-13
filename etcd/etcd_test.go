package etcd

import (
	"testing"
	"time"
)

func TestInitEtcd(t *testing.T) {
	ca := "./ca.pem"
	key := "./etcd-key.pem"
	cert := "./etcd.pem"

	cli, err := InitEtcd(ca, key, cert,
		WithTimeOut(3*time.Second),
		WithEndpoints([]string{"https://10.0.2.114:2379", "10.0.2.115:2379"}))
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		wch := cli.Watch("/test")

		for item := range wch {
			for _, ev := range item.Events {
				t.Logf("Type:%s, key:%s, value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()

	cli.Put("/test", "11111")
	cli.Put("/test/2", "22222")
	result, _ := cli.Get("/test")
	t.Logf("%#v\n", result)
	cli.Delete("/test")
	result, _ = cli.Get("/test")
	t.Logf("%#v\n", result)

	time.Sleep(time.Second)

}

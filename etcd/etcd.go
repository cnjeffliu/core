/*
 * @Author: Jeffrey.Liu
 * @Date: 2021-09-06 16:08:13
 * @LastEditors: Jeffrey.Liu
 * @LastEditTime: 2021-12-15 15:32:16
 * @Description:
 */
package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"go.etcd.io/etcd/clientv3"
)

type EtcdClient struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
	tm     time.Duration // 操作超时时间
}

type etcdCliOption func(c *clientv3.Config)

const defaultTimeout time.Duration = 5 * time.Second

var (
	defaultEndpoints = []string{"127.0.0.1:2379"}
)

func WithEndpoints(endpoints []string) etcdCliOption {
	return func(c *clientv3.Config) {
		c.Endpoints = endpoints
	}
}

func WithTimeOut(timeout time.Duration) etcdCliOption {
	return func(c *clientv3.Config) {
		c.DialTimeout = timeout
	}
}

func InitEtcd(ca string, key string, cert string, opts ...etcdCliOption) (etcdCli *EtcdClient, err error) {
	fmt.Println("begin InitEtcd ...")
	certData, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	caData, err := ioutil.ReadFile(ca)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certData},
		RootCAs:      pool,
	}

	cfg := clientv3.Config{
		Endpoints:   defaultEndpoints,
		DialTimeout: defaultTimeout,
		TLS:         tlsConfig,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)

	etcdCli = &EtcdClient{
		client: client,
		kv:     kv,
		lease:  lease,
		tm:     time.Second * 10,
	}
	return etcdCli, nil
}

func (cli *EtcdClient) Get(key string) (result []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), cli.tm)
	defer cancel()
	resp, err := cli.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, kvPair := range resp.Kvs {
		result = append(result, string(kvPair.Value))
	}
	return result, nil
}

func (cli *EtcdClient) Put(key, value string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), cli.tm)
	defer cancel()
	_, err = cli.kv.Put(ctx, key, value, clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (cli *EtcdClient) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), cli.tm)
	defer cancel()
	_, err := cli.kv.Delete(ctx, key, clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (cli *EtcdClient) PutWithLease(key, value string, timeout int64) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), cli.tm)
	defer cancel()
	leaseResp, err := cli.lease.Grant(ctx, timeout)
	if err != nil {
		return
	}

	_, err = cli.kv.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return
	}
	return
}

func (cli *EtcdClient) Watch(key string) clientv3.WatchChan {
	ctx, cancel := context.WithTimeout(context.Background(), cli.tm)
	defer cancel()
	ch := cli.client.Watch(ctx, key, clientv3.WithPrefix())
	return ch
}

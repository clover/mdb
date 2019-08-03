// Copyright 2019 Clover Network, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.etcd.io/etcd/client"
	kerr "go.etcd.io/etcd/etcdserver/api/v2error"
)

func etcdConnect(host string, port int) client.KeysAPI {
	endpoint := fmt.Sprintf("http://%s:%d", host, port)
	cfg := client.Config{
		Endpoints:               []string{endpoint},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 2 * time.Second,
	}

	c, err := client.New(cfg)
	check(err)

	kapi := client.NewKeysAPI(c)
	return kapi
}

func isKeyNotFound(err error) bool {
	if err != nil {
		if e, ok := err.(client.Error); ok {
			if e.Code == kerr.EcodeKeyNotFound {
				return true
			}
		}
	}

	return false
}

func isKeyChanged(err error) bool {
	if err != nil {
		if e, ok := err.(client.Error); ok {
			if e.Code == kerr.EcodeTestFailed {
				return true
			}
		}
	}

	return false
}

func isKeyExists(err error) bool {
	if err != nil {
		if e, ok := err.(client.Error); ok {
			if e.Code == kerr.EcodeNodeExist {
				return true
			}
		}
	}

	return false
}

// EcodeNotFile
func isKeyNotFile(err error) bool {
	if err != nil {
		if e, ok := err.(client.Error); ok {
			if e.Code == kerr.EcodeNotFile {
				return true
			}
		}
	}

	return false
}

// Get all the "files" in a directory, rejecting subdirs explicitly
func getfilelist(k client.KeysAPI, path string) ([]string, error) {
	resp, err := k.Get(context.Background(), path, &client.GetOptions{Recursive: false})
	if err != nil {
		return []string{}, err
	}

	nodes := []string{}
	for _, node := range resp.Node.Nodes {
		if node.Dir {
			logwarn("Unexpected directory in filelist: %s\n", node.Key)
		} else {
			tag := node.Key
			nodes = append(nodes, strings.TrimPrefix(tag, path))
		}
	}

	return nodes, nil
}

// Get all the "files" in a directory, including directories (but not
// recursively). Doesn't distinguish between file and directory (is
// this a FIXME?)
func getdirlist(k client.KeysAPI, path string) ([]string, error) {
	resp, err := k.Get(context.Background(), path, &client.GetOptions{Recursive: false})
	if err != nil {
		return []string{}, err
	}

	nodes := []string{}
	for _, node := range resp.Node.Nodes {
		name := node.Key
		nodes = append(nodes, strings.TrimPrefix(name, path))
	}

	return nodes, nil
}

func checknode(k client.KeysAPI, path string) (bool, error) {
	_, err := k.Get(context.Background(), path, &client.GetOptions{Recursive: false})

	if err == nil {
		return true, nil
	}

	if isKeyNotFound(err) {
		return false, nil
	}

	return false, err
}

func getnode(k client.KeysAPI, path string) (string, error) {
	resp, err := k.Get(context.Background(), path, &client.GetOptions{Recursive: false})
	// check(err)

	if err != nil {
		return "", err
	}

	if resp.Node.Dir {
		logwarn("Unexpectedly retrieved directory when retrieving node: %s\n", path)
		return "", nil
	} else {
		value := resp.Node.Value
		return value, nil
	}
}

// Conditional set if onlyif  != "", otherwise always set
func setnode(k client.KeysAPI, path string, value string, onlyif string) error {
	logdebug("Setting %s to value %s\n", path, value)

	opts := client.SetOptions{
		PrevValue: onlyif,
		Dir:       false,
	}

	_, err := k.Set(context.Background(), path, value, &opts)
	return err
}

func incrnode(k client.KeysAPI, path string, increment int64) (int64, error) {
	for tries := 0; tries < 10; tries++ {
		cur, err := getnode(k, path)
		if err != nil {
			return -1, err
		}

		curint, err := strconv.ParseInt(cur, 10, 32)
		if err != nil {
			return -1, err
		}

		new := curint + increment
		err = setnode(k, path, strconv.FormatInt(new, 10), cur)

		if err == nil {
			return new, nil
		}

		if isKeyChanged(err) {
			continue
		}

		return -1, err
	}

	// If we fell out the bottom of the loop, we've retried too much
	return -1, errors.New("Couldn't atomically increment variable after 10 tries")
}

func deletenode(k client.KeysAPI, path string, recursive bool) error {
	logdebug("Deleting %s\n", path)
	delopts := client.DeleteOptions{
		Recursive: recursive,
		Dir:       recursive,
	}
	_, err := k.Delete(context.Background(), path, &delopts)
	return err
}

func createdir(k client.KeysAPI, path string) error {
	logdebug("Creating directory %s\n", path)
	_, err := k.Set(context.Background(), path, "", &client.SetOptions{Dir: true})
	return err
}

func watch(k client.KeysAPI, path string) (string, error) {
	log("Watching %s\n", path)
	w := k.Watcher(path, &client.WatcherOptions{Recursive: false})

	resp, err := w.Next(context.TODO())
	if err != nil {
		return "", err
	}

	// If the watch was on a directory, return no value, no error
	if resp.Node.Dir == true {
		return "", nil
	}

	return resp.Node.Value, nil
}

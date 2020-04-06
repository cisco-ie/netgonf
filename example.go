/**
 * Copyright (c) 2019-2020 Cisco Systems
 *
 * Author: Steven Barth <stbarth@cisco.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cisco-ie/netgonf/netconf"
	"golang.org/x/crypto/ssh"
)

type yangPush struct {
	netconf.Notification
	PushUpdate struct {
		Content struct {
			InnerXML []byte `xml:",innerxml"`
		} `xml:"datastore-contents-xml"`
	} `xml:"urn:ietf:params:xml:ns:yang:ietf-yang-push push-update"`
}

type establishSubscription struct {
	XMLName     xml.Name `xml:"urn:ietf:params:xml:ns:yang:ietf-event-notifications establish-subscription"`
	YangPush    string   `xml:"xmlns:yp,attr"`
	Stream      string   `xml:"stream"`
	XPathFilter string   `xml:"urn:ietf:params:xml:ns:yang:ietf-yang-push xpath-filter"`
	Period      int      `xml:"urn:ietf:params:xml:ns:yang:ietf-yang-push period"`
}

func main() {
	var address string
	var username string
	var password string
	var stream string
	var filter string
	var get string
	var period int
	var keyfile string

	flag.StringVar(&address, "address", "localhost", "Address")
	flag.StringVar(&username, "user", "cisco", "Username")
	flag.StringVar(&password, "pass", "cisco", "Password")
	flag.StringVar(&keyfile, "keyfile", "", "SSH-Key")
	flag.StringVar(&stream, "stream", "", "Stream")
	flag.StringVar(&filter, "filter", "/process-cpu-ios-xe-oper:cpu-usage/cpu-utilization/five-seconds", "Filter")
	flag.StringVar(&get, "get", "/interfaces-state/interface[name='TenGigabitEthernet1/0/1']/oper-status", "Get")
	flag.IntVar(&period, "period", 3, "Period")
	flag.Parse()

	var client netconf.Client
	var err error

	if len(keyfile) == 0 {
		client, err = netconf.DialSSHWithPassword(address, username, password, ssh.InsecureIgnoreHostKey())
	} else {
		var key []byte
		key, err = ioutil.ReadFile(keyfile)
		if err == nil {
			var signer ssh.Signer
			signer, err = ssh.ParsePrivateKey(key)
			if err == nil {
				client, err = netconf.DialSSHWithPublicKey(address, username, signer, ssh.InsecureIgnoreHostKey())
			}
		}
	}
	if err != nil {
		log.Println(err.Error())
		os.Exit(2)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		session, err := client.NewSession()
		if err != nil {
			client.Close()
			log.Println(err.Error())
			os.Exit(3)
		}

		if len(stream) > 0 {
			subscribe := &netconf.CreateSubscription{
				Stream: &stream,
			}
			err = session.CallSimple(subscribe)
		} else {
			establish := &establishSubscription{
				YangPush:    "urn:ietf:params:xml:ns:yang:ietf-yang-push",
				Stream:      "yp:yang-push",
				XPathFilter: filter,
				Period:      period * 100,
			}
			err = session.CallSimple(establish)
		}
		if err != nil {
			log.Println(err.Error())
		} else {
			push := &yangPush{}
			for {
				err = session.Receive(push)
				if err != nil {
					log.Println(err.Error())
				} else {
					fmt.Printf("%s: %s\n\n", push.EventTime, push.PushUpdate.Content.InnerXML)
				}
			}
		}

		session.Close()
	}()

	go func() {
		defer wg.Done()
		client, _ = netconf.DialSSHWithPassword(address, username, password, ssh.InsecureIgnoreHostKey())
		session, err := client.NewSession()
		if err != nil {
			client.Close()
			log.Println(err.Error())
			os.Exit(3)
		}

		request := netconf.Get{Filter: &netconf.Filter{Type: "xpath", Select: get}}
		response := netconf.RPCReplyData{}
		for len(get) > 0 {
			if err := session.Call(&request, &response); err != nil {
				log.Println(err)
			}

			fmt.Printf("%s\n\n", response.Data.InnerXML)

			select {
			case <-time.After(time.Duration(period) * time.Second):
			}
		}

		session.Close()
	}()

	wg.Wait()
	client.Close()
}

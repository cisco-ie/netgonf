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

package netconf

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"strings"
)

// ErrCapabilitiesExchange indicates a failed NETCONF hello-exchange due to incompatible versions or invalid session ID
var ErrCapabilitiesExchange = errors.New("Capabilities exchange failed")

// Client defines a transport-independent interface for NETCONF clients
type Client interface {
	io.Closer
	NewSession() (*Session, error)
}

// Session represents a session towards the server
type Session struct {
	SessionID    uint64
	Capabilities map[string]string

	transport   io.ReadWriteCloser
	newFramer   func(io.Writer) io.WriteCloser
	newUnframer func(io.Reader) io.ReadCloser
	messageID   int
}

func newSession(transport io.ReadWriteCloser) (*Session, error) {
	session := Session{
		transport:   transport,
		newFramer:   newFramerV10,
		newUnframer: newUnframerV10,
	}

	// Exchange capabilities
	hello := &struct {
		XMLName      xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 hello"`
		SessionID    uint64   `xml:"session-id,omitempty"`
		Capabilities []string `xml:"capabilities>capability"`
	}{Capabilities: []string{CapNetconf10, CapNetconf11}}

	// Send hello message
	framer := session.newFramer(transport)
	if err := xml.NewEncoder(framer).Encode(hello); err != nil {
		return nil, err
	}
	if err := framer.Close(); err != nil {
		return nil, err
	}

	// Retrieve and parse server's hello message
	reader := session.NewReader()
	hello.Capabilities = hello.Capabilities[:0]
	if err := xml.NewDecoder(reader).Decode(&hello); err != nil {
		return nil, err
	} else if err := reader.Close(); err != nil {
		return nil, err
	}

	// Check for non-0 session ID
	if session.SessionID = hello.SessionID; session.SessionID == 0 {
		return nil, ErrCapabilitiesExchange
	}

	// Parse capabilities
	session.Capabilities = make(map[string]string)
	for _, capability := range hello.Capabilities {
		cap := strings.SplitN(capability, "?", 2)
		if len(cap) > 1 {
			session.Capabilities[cap[0]] = cap[1]
		} else {
			session.Capabilities[cap[0]] = ""
		}
	}

	// Check for compatible version and switch framing method if necessary
	if _, compatible := session.Capabilities[CapNetconf11]; compatible {
		session.newFramer = newFramerV11
		session.newUnframer = newUnframerV11
	} else if _, compatible := session.Capabilities[CapNetconf10]; !compatible {
		return nil, ErrCapabilitiesExchange
	}

	return &session, nil
}

// Call a NETCONF RPC and retrieve its reply
func (s *Session) Call(request interface{}, response interface{}) error {
	var err error
	s.messageID++
	writer := s.NewWriter()
	element := xml.StartElement{
		Name: xml.Name{Local: "rpc", Space: "urn:ietf:params:xml:ns:netconf:base:1.0"},
		Attr: []xml.Attr{{Name: xml.Name{Local: "message-id"}, Value: strconv.Itoa(s.messageID)}},
	}
	rpc := &struct{ Operation interface{} }{Operation: request}
	if err = xml.NewEncoder(writer).EncodeElement(rpc, element); err == nil {
		err = writer.Close()
	}

	// Read until rpc-reply (skip spurious notifications etc.)
	for haveReply := false; !haveReply && err == nil && response != nil; {
		reader := s.NewReader()
		decoder := xml.NewDecoder(reader)

		for err == nil {
			var token xml.Token
			token, err = decoder.Token() // Read until the XML document root
			if element, ok := token.(xml.StartElement); ok {
				if element.Name.Local == "rpc-reply" {
					err = decoder.DecodeElement(response, &element)
					haveReply = true
				}
				break
			}
		}

		reader.Close()
	}
	return err
}

// NewReader creates a low-level reader for receiving the next NETCONF message
func (s *Session) NewReader() io.ReadCloser {
	return s.newUnframer(s.transport)
}

// NewWriter creates a low-level writing for sending the next NETCONF message
func (s *Session) NewWriter() io.WriteCloser {
	return s.newFramer(s.transport)
}

// Receive a message from the server, e.g. a notification
func (s *Session) Receive(response interface{}) error {
	reader := s.NewReader()
	err := xml.NewDecoder(reader).Decode(response)
	errReader := reader.Close()
	if err == nil {
		err = errReader
	}
	return err
}

// CallSimple calls a NETCONF RPC and returns the first rpc-error or nil if there was none
func (s *Session) CallSimple(request interface{}) error {
	reply := &RPCReply{}
	err := s.Call(request, reply)
	if err == nil && len(reply.RPCError) > 0 {
		err = &reply.RPCError[0]
	}
	return err
}

// Close the session gracefully
func (s *Session) Close() error {
	closeSession := &struct {
		XMLName xml.Name `xml:"close-session"`
	}{}
	err := s.CallSimple(closeSession)
	errTransport := s.transport.Close()
	if err == nil {
		err = errTransport
	}
	return err
}

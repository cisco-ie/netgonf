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
	"fmt"
	"strings"
	"time"
)

// RPCReply models a NETCONF <rpc-reply> element for inclusion in other structs
type RPCReply struct {
	XMLName  xml.Name   `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 rpc-reply"`
	RPCError []RPCError `xml:"rpc-error"`
}

// RPCReplyData models a NETCONF <rpc-reply> element with a data child
type RPCReplyData struct {
	RPCReply
	Data struct {
		InnerXML []byte `xml:",innerxml"`
	} `xml:"data"`
}

// RPCError models a NETCONF <rpc-error> element
type RPCError struct {
	XMLName       xml.Name `xml:"rpc-error"`
	ErrorType     string   `xml:"error-type"`
	ErrorTag      string   `xml:"error-tag"`
	ErrorSeverity string   `xml:"error-severity"`
	ErrorAppTag   string   `xml:"error-app-tag"`
	ErrorPath     string   `xml:"error-path"`
	ErrorMessage  string   `xml:"error-message"`
	ErrorInfo     struct {
		BadElement   string `xml:"bad-element"`
		BadAttribute string `xml:"bad-attribute"`
		BadNamespace string `xml:"bad-namespace"`
		SessionID    string `xml:"session-id"`
		InnerXML     []byte `xml:",innerxml"`
	} `xml:"error-info"`
}

func (e RPCError) Error() string {
	return fmt.Sprintf("NETCONF RPC error %s: %s", e.ErrorTag, e.ErrorMessage)
}

// Filter defines the basic structure of a filter for <get> and <get-config> operations
type Filter struct {
	Type    string `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 type,attr,omitempty"`
	Select  string `xml:"select,attr,omitempty"`
	Subtree string `xml:",innerxml"`
}

// Get defines the <get> operation for use with Session.Call
type Get struct {
	XMLName      xml.Name     `xml:"get"`
	Filter       *Filter      `xml:"filter,omitempty"`
	WithDefaults DefaultsMode `xml:"urn:ietf:params:xml:ns:yang:ietf-netconf-with-defaults with-defaults,omitempty"`
}

// GetConfig defines the <get-config> operation for use with Session.Call
type GetConfig struct {
	XMLName      xml.Name     `xml:"get-config"`
	Source       Datastore    `xml:"source"`
	Filter       *Filter      `xml:"filter,omitempty"`
	WithDefaults DefaultsMode `xml:"urn:ietf:params:xml:ns:yang:ietf-netconf-with-defaults with-defaults,omitempty"`
}

// EditConfig defines the <edit-config> operation for use with Session.CallProcedure
type EditConfig struct {
	XMLName          xml.Name          `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 edit-config"`
	Target           Datastore         `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 target"`
	DefaultOperation *DefaultOperation `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 default-operation,omitempty"`
	TestOption       *TestOption       `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 test-option,omitempty"`
	ErrorOption      *ErrorOption      `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 error-option,omitempty"`
	Config           struct {
		InnerXML []byte `xml:",innerxml"`
	} `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 config"`
}

// CopyConfig defines the <copy-config> operation for use with Session.CallProcedure
type CopyConfig struct {
	XMLName xml.Name  `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 copy-config"`
	Target  Datastore `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 target"`
	Source  Datastore `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 source"`
}

// DeleteConfig defines the <delete-config> operation for use with Session.CallProcedure
type DeleteConfig struct {
	XMLName xml.Name  `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 delete-config"`
	Target  Datastore `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 target"`
}

// Lock defines the <lock> operation for use with Session.CallProcedure
type Lock struct {
	XMLName xml.Name  `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 lock"`
	Target  Datastore `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 target"`
}

// Unlock defines the <unlock> operation for use with Session.CallProcedure
type Unlock struct {
	XMLName xml.Name  `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 lock"`
	Target  Datastore `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 target"`
}

// KillSession defines the <kill-session> operation for use with Session.CallProcedure
type KillSession struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 kill-session"`
	SessionID uint64   `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 session-id"`
}

// Commit defines the <commit> operation for use with Session.CallProcedure
type Commit struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 commit"`
	PersistID *string  `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 persist-id,omitempty"`
}

// CommitConfirmed defines the <commit> operation for confirmed commits for use with Session.CallProcedure
type CommitConfirmed struct {
	XMLName        xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 commit"`
	Confirmed      struct{} `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 confirmed"`
	ConfirmTimeout *uint    `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 confirm-timeout,omitempty"`
	Persist        *string  `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 persist,omitempty"`
}

// CancelCommit defines the <cancel-commit> operation for use with Session.CallProcedure
type CancelCommit struct {
	XMLName   xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 cancel-commit"`
	PersistID *string  `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 persist-id,omitempty"`
}

// DiscardChanges defines the <discard-changes> operation for use with Session.CallProcedure
type DiscardChanges struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 discard-changes"`
}

// Validate defines the <validate> operation for use with Session.CallProcedure
type Validate struct {
	XMLName xml.Name  `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 validate"`
	Source  Datastore `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 source"`
}

// ValidateConfig defines the <validate> operation on an explicit config for use with Session.CallProcedure
type ValidateConfig struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 validate"`
	Config  struct {
		InnerXML []byte `xml:",innerxml"`
	} `xml:"source>config"`
}

// GetSchema defines the <get-schema> operation for use with Session.Call
type GetSchema struct {
	XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring get-schema"`
	Identifier string   `xml:"urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring identifier"`
	Version    *string  `xml:"urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring version,omitempty"`
	Format     *string  `xml:"urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring format,omitempty"`
}

// CreateSubscription defines the <create-subscription> operation for use with Session.CallProcedure
type CreateSubscription struct {
	XMLName   xml.Name   `xml:"urn:ietf:params:xml:ns:netconf:notification:1.0 create-subscription"`
	Stream    *string    `xml:"urn:ietf:params:xml:ns:netconf:notification:1.0 stream,omitempty"`
	Filter    *Filter    `xml:"urn:ietf:params:xml:ns:netconf:notification:1.0 filter,omitempty"`
	StartTime *time.Time `xml:"urn:ietf:params:xml:ns:netconf:notification:1.0 startTime,omitempty"`
	StopTime  *time.Time `xml:"urn:ietf:params:xml:ns:netconf:notification:1.0 stoptTime,omitempty"`
}

// Notification describes the basic struct for NETCONF notifications (RFC 5277)
type Notification struct {
	XMLName   xml.Name  `xml:"urn:ietf:params:xml:ns:netconf:notification:1.0 notification"`
	EventTime time.Time `xml:"eventTime"`
}

// Action defines a Yang 1.1 action
type Action struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:yang:1 action"`
	InnerXML []byte   `xml:",innerxml"`
}

// TailfAction defines a tailf:action extension
type TailfAction struct {
	XMLName xml.Name `xml:"http://tail-f.com/ns/netconf/actions/1.0 action"`
	Data    struct {
		InnerXML []byte `xml:",innerxml"`
	} `xml:"data"`
}

// DefaultOperation specifies default behavior for edit-config operation
type DefaultOperation string

// TestOption specifies test-behavior for edit-config operation
type TestOption string

// ErrorOption specifies error-handling for edit-config operation
type ErrorOption string

// DefaultsMode specifies how to handle default values specified in YANG
type DefaultsMode string

// NCSCommitParameter defines custom commit parameters for NCS
type NCSCommitParameter string

// Datastore on NETCONF agent
type Datastore string

// MarshalXML datastore into XML depending if it is a URL (contains a ':') or not
func (d Datastore) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var element interface{}
	datastore := string(d)
	if strings.ContainsRune(datastore, ':') {
		element = &struct {
			URL string `xml:"urn:ietf:params:xml:ns:netconf:base:1.0 url"`
		}{URL: datastore}
	} else {
		element = &struct {
			Datastore struct {
				XMLName xml.Name
			}
		}{Datastore: struct{ XMLName xml.Name }{XMLName: xml.Name{Local: datastore}}}
	}
	e.EncodeElement(element, start)
	return nil
}

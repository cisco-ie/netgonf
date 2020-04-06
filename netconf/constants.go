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

// List of standardized constants related to the NETCONF protocol and its extensions
const (
	NsNetconf             = "urn:ietf:params:xml:ns:netconf:base:1.0"
	NsNetconfWithDefaults = "urn:ietf:params:xml:ns:yang:ietf-netconf-with-defaults"
	NsNetconfNotification = "urn:ietf:params:xml:ns:netconf:notification:1.0"
	NsNetmodNotification  = "urn:ietf:params:xml:ns:netmod:notification"
	NsNetconfMonitoring   = "urn:ietf:params:xml:ns:yang:ietf-netconf-monitoring"
	NsTailfActions        = "http://tail-f.com/ns/netconf/actions/1.0"

	CapNetconf10       = "urn:ietf:params:netconf:base:1.0"
	CapNetconf11       = "urn:ietf:params:netconf:base:1.1"
	CapConfirmedCommit = "urn:ietf:params:netconf:capability:confirmed-commit:1.1"
	CapValidate        = "urn:ietf:params:netconf:capability:validate:1.1"
	CapWithDefaults    = "urn:ietf:params:netconf:capability:with-defaults:1.0"
	CapNotifiction     = "urn:ietf:params:netconf:capability:notification:1.0"
	CapInterleave      = "urn:ietf:params:netconf:capability:interleave:1.0"
	CapStartup         = "urn:ietf:params:netconf:capability:startup:1.0"
	CapWritableRunning = "urn:ietf:params:netconf:capability:writable-running:1.0"
	CapCandidate       = "urn:ietf:params:netconf:capability:candidate:1.0"
	CapRollbackOnError = "urn:ietf:params:netconf:capability:rollback-on-error:1.0"
	CapURL             = "urn:ietf:params:netconf:capability:url:1.0"
	CapXPath           = "urn:ietf:params:netconf:capability:xpath:1.0"
	CapMonitoring      = NsNetconfMonitoring
	CapTailfActions    = NsTailfActions

	Running   Datastore = "running"
	Candidate Datastore = "candidate"
	Startup   Datastore = "startup"
	Intended  Datastore = "intended"

	OpMerge   DefaultOperation = "merge"
	OpReplace DefaultOperation = "replace"
	OpNone    DefaultOperation = "none"

	TestThenSet TestOption = "test-then-set"
	TestOnlySet TestOption = "set"
	TestOnly    TestOption = "test-only"

	StopOnError     ErrorOption = "stop-on-error"
	ContinueOnError ErrorOption = "continue-on-error"
	RollbackOnError ErrorOption = "rollback-on-error"

	ReportAll       DefaultsMode = "report-all"
	ReportAllTagged DefaultsMode = "report-all-tagged"
	Trim            DefaultsMode = "trim"
	Explicit        DefaultsMode = "explicit"
)

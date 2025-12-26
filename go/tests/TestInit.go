// © 2025 Sharon Aicler (saichler@gmail.com)
//
// Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package tests provides test infrastructure and utilities for L8 Pollaris.
// This file contains the test setup and teardown functions, as well as
// the test topology configuration.
package tests

import (
	"github.com/saichler/l8bus/go/overlay/protocol"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	. "github.com/saichler/l8test/go/infra/t_resources"
	. "github.com/saichler/l8test/go/infra/t_topology"
	. "github.com/saichler/l8types/go/ifs"
	"github.com/saichler/probler/go/prob/common"
)

// topo is the global test topology used by all tests.
// It provides a simulated distributed environment with multiple VNics.
var topo *TestTopology

// init configures the test environment at package load time.
// It sets the log level to Trace for detailed test output and
// initializes the TargetLinks with a mock implementation.
func init() {
	Log.SetLogLevel(Trace_Level)
	targets.Links = &common.Links{}
}

// setup initializes the test topology before tests run.
func setup() {
	setupTopology()
}

// tear shuts down the test topology after tests complete.
func tear() {
	shutdownTopology()
}

// reset clears handler state between tests and logs test completion.
func reset(name string) {
	Log.Info("*** ", name, " end ***")
	topo.ResetHandlers()
}

// setupTopology creates a new test topology with 4 nodes across 3 port ranges.
// Message logging is enabled for debugging test failures.
func setupTopology() {
	protocol.MessageLog = true
	topo = NewTestTopology(4, []int{20000, 30000, 40000}, Info_Level)
}

// shutdownTopology gracefully shuts down all nodes in the test topology.
func shutdownTopology() {
	topo.Shutdown()
}

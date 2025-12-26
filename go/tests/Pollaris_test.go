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

// Package tests contains integration tests for the L8 Pollaris service.
// These tests verify the functionality of polling configuration management
// including creation, retrieval by name, and retrieval by group.
package tests

import (
	"testing"

	"github.com/saichler/l8collector/go/collector/common"
	"github.com/saichler/l8parser/go/parser/boot"
	"github.com/saichler/l8pollaris/go/pollaris"
	"github.com/saichler/l8types/go/ifs"
)

// TestMain is the test entry point that sets up and tears down the test topology.
// It initializes the distributed test environment before running tests
// and ensures proper cleanup afterward.
func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

// TestPollaris verifies the core Pollaris service functionality:
// 1. Registers and activates the PollarisService
// 2. Creates a polling configuration using boot.CreateBoot01()
// 3. Posts the configuration to the PollarisCenter
// 4. Retrieves the configuration by name
// 5. Retrieves the configuration by group (BOOT_STAGE_01)
func TestPollaris(t *testing.T) {
	vnic := topo.VnicByVnetNum(2, 2)
	vnic.Resources().Registry().Register(pollaris.PollarisService{})
	sla := ifs.NewServiceLevelAgreement(&pollaris.PollarisService{}, pollaris.ServiceName, 0, true, nil)
	vnic.Resources().Services().Activate(sla, vnic)
	p := pollaris.Pollaris(vnic.Resources())
	pollrs := boot.CreateBoot01()
	err := p.Post(pollrs, false)
	if err != nil {
		vnic.Resources().Logger().Fail(t, err.Error())
		return
	}

	byName := p.PollarisByName(pollrs.Name)
	if byName == nil {
		vnic.Resources().Logger().Fail(t, "No such pollaris")
		return
	}

	byGroup, err := pollaris.PollarisByGroup(vnic.Resources(), common.BOOT_STAGE_01,
		"", "", "", "", "", "")
	if byGroup == nil {
		vnic.Resources().Logger().Fail(t, "No such group")
		return
	}
}

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

package tests

import (
	"fmt"
	"github.com/saichler/l8pollaris/go/pollaris/targets"
	"github.com/saichler/l8pollaris/go/types/l8tpollaris"
	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/probler/go/prob/common"
	"github.com/saichler/probler/go/prob/common/creates"
	"strconv"
	"testing"
	"time"
)

// TestTargetService verifies the Target service functionality:
// 1. Activates the Targets service with PostgreSQL backend
// 2. Creates and posts a single target device
// 3. Updates the target state to UP using Patch
// 4. Creates and posts a batch of 100 target devices
// 5. Waits for collector distribution to complete
func TestTargetService(t *testing.T) {
	nic := topo.VnicByVnetNum(2, 2)
	targets.Activate("postgres", "probler", nic)
	tr, _ := targets.Targets(nic)
	device := creates.CreateDevice("10.10.10.10", common.NetworkDevice_Links_ID, "sim")
	resp := tr.Post(object.New(nil, device), nic)
	if resp.Error() != nil {
		fmt.Println(resp.Error().Error())
		nic.Resources().Logger().Fail(t, resp.Error().Error())
		return
	}

	device = &l8tpollaris.L8PTarget{TargetId: device.TargetId}
	device.State = l8tpollaris.L8PTargetState_Up
	resp = tr.Patch(object.New(nil, device), nic)
	if resp.Error() != nil {
		fmt.Println(resp.Error().Error())
		nic.Resources().Logger().Fail(t, resp.Error().Error())
		return
	}

	deviceList := &l8tpollaris.L8PTargetList{List: make([]*l8tpollaris.L8PTarget, 0)}

	ip := 1
	sub := 40
	for i := 1; i <= 100; i++ {
		device = creates.CreateDevice("60.50."+strconv.Itoa(sub)+"."+strconv.Itoa(ip), common.NetworkDevice_Links_ID, "sim")
		deviceList.List = append(deviceList.List, device)
		ip++
		if ip > 254 {
			sub++
			ip = 1
		}
	}

	resp = tr.Post(object.New(nil, deviceList), nic)
	if resp.Error() != nil {
		fmt.Println(resp.Error().Error())
		nic.Resources().Logger().Fail(t, resp.Error().Error())
		return
	}
	fmt.Println("Sleeping")
	time.Sleep(time.Second * 15)
}

/*
 * Copyright The Dragonfly Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dfgettask

import (
	"context"
	"testing"

	"github.com/dragonflyoss/Dragonfly/apis/types"
	"github.com/dragonflyoss/Dragonfly/pkg/errortypes"

	"github.com/go-check/check"
	"github.com/prometheus/client_golang/prometheus"
	prom_testutil "github.com/prometheus/client_golang/prometheus/testutil"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

func init() {
	check.Suite(&DfgetTaskMgrTestSuite{})
}

type DfgetTaskMgrTestSuite struct {
}

// SetUpTest does common setup in the beginning of each test.
func (s *DfgetTaskMgrTestSuite) SetUpTest(c *check.C) {
	// In every test, we should reset Prometheus default registry, otherwise
	// it will panic because of duplicate metricsutils.
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
}

func (s *DfgetTaskMgrTestSuite) TestDfgetTaskAdd(c *check.C) {
	manager, _ := NewManager()
	dfgetNum := manager.metrics.dfgetTasks
	dfgetDaemonNum := manager.metrics.dfgetTasksDaemon

	var testCases = []struct {
		dfgetTask *types.DfGetTask
		Expect    *types.DfGetTask
	}{
		{
			dfgetTask: &types.DfGetTask{
				CID:        "foo",
				CallSystem: "foo",
				Dfdaemon:   true,
				Path:       "/peer/file/taskFileName",
				PieceSize:  4 * 1024 * 1024,
				TaskID:     "test1",
				PeerID:     "peer1",
			},
			Expect: &types.DfGetTask{
				CID:        "foo",
				CallSystem: "foo",
				Dfdaemon:   true,
				Path:       "/peer/file/taskFileName",
				PieceSize:  4 * 1024 * 1024,
				TaskID:     "test1",
				PeerID:     "peer1",
				Status:     types.DfGetTaskStatusWAITING,
			},
		},
		{
			dfgetTask: &types.DfGetTask{
				CID:        "bar",
				CallSystem: "bar",
				Dfdaemon:   true,
				Path:       "/peer/file/taskFileName",
				PieceSize:  4 * 1024 * 1024,
				TaskID:     "test2",
				PeerID:     "peer2",
			},
			Expect: &types.DfGetTask{
				CID:        "bar",
				CallSystem: "bar",
				Dfdaemon:   true,
				Path:       "/peer/file/taskFileName",
				PieceSize:  4 * 1024 * 1024,
				TaskID:     "test2",
				PeerID:     "peer2",
				Status:     types.DfGetTaskStatusWAITING,
			},
		},
	}

	for _, tc := range testCases {
		err := manager.Add(context.Background(), tc.dfgetTask)
		c.Check(err, check.IsNil)
		if tc.dfgetTask.Dfdaemon {
			c.Assert(1, check.Equals,
				int(prom_testutil.ToFloat64(
					dfgetDaemonNum.WithLabelValues(tc.dfgetTask.TaskID, tc.dfgetTask.CallSystem))))
		} else {
			c.Assert(1, check.Equals,
				int(prom_testutil.ToFloat64(
					dfgetNum.WithLabelValues(tc.dfgetTask.TaskID, tc.dfgetTask.CallSystem))))
		}

		dt, err := manager.Get(context.Background(), tc.dfgetTask.CID, tc.dfgetTask.TaskID)
		c.Check(err, check.IsNil)
		c.Check(dt, check.DeepEquals, tc.Expect)
	}
}

func (s *DfgetTaskMgrTestSuite) TestDfgetTaskUpdate(c *check.C) {
	manager, _ := NewManager()
	var testCases = []struct {
		dfgetTask  *types.DfGetTask
		taskStatus string
		Expect     *types.DfGetTask
	}{
		{
			dfgetTask: &types.DfGetTask{
				CID:        "foo",
				CallSystem: "foo",
				Dfdaemon:   true,
				Path:       "/peer/file/taskFileName",
				PieceSize:  4 * 1024 * 1024,
				TaskID:     "test1",
				PeerID:     "peer1",
			},
			taskStatus: types.DfGetTaskStatusFAILED,
			Expect: &types.DfGetTask{
				CID:        "foo",
				CallSystem: "foo",
				Dfdaemon:   true,
				Path:       "/peer/file/taskFileName",
				PieceSize:  4 * 1024 * 1024,
				TaskID:     "test1",
				PeerID:     "peer1",
				Status:     types.DfGetTaskStatusFAILED,
			},
		},
		{
			dfgetTask: &types.DfGetTask{
				CID:        "bar",
				CallSystem: "bar",
				Dfdaemon:   true,
				Path:       "/peer/file/taskFileName",
				PieceSize:  4 * 1024 * 1024,
				TaskID:     "test2",
				PeerID:     "peer2",
			},
			taskStatus: types.DfGetTaskStatusSUCCESS,
			Expect: &types.DfGetTask{
				CID:        "bar",
				CallSystem: "bar",
				Dfdaemon:   true,
				Path:       "/peer/file/taskFileName",
				PieceSize:  4 * 1024 * 1024,
				TaskID:     "test2",
				PeerID:     "peer2",
				Status:     types.DfGetTaskStatusSUCCESS,
			},
		},
	}

	for _, tc := range testCases {
		err := manager.Add(context.Background(), tc.dfgetTask)
		c.Check(err, check.IsNil)

		err = manager.UpdateStatus(context.Background(), tc.dfgetTask.CID, tc.dfgetTask.TaskID, tc.taskStatus)
		c.Check(err, check.IsNil)

		dt, err := manager.Get(context.Background(), tc.dfgetTask.CID, tc.dfgetTask.TaskID)
		c.Check(dt, check.DeepEquals, tc.Expect)
	}
}

func (s *DfgetTaskMgrTestSuite) TestDfgetTaskDelete(c *check.C) {
	manager, _ := NewManager()
	dfgetNum := manager.metrics.dfgetTasks
	dfgetDaemonNum := manager.metrics.dfgetTasksDaemon

	var testCases = []struct {
		dfgetTask *types.DfGetTask
	}{
		{
			dfgetTask: &types.DfGetTask{
				CID:        "foo",
				CallSystem: "foo",
				Dfdaemon:   true,
				Path:       "/peer/file/taskFileName",
				PieceSize:  4 * 1024 * 1024,
				TaskID:     "test1",
				PeerID:     "peer1",
			},
		},
		{
			dfgetTask: &types.DfGetTask{
				CID:        "bar",
				CallSystem: "bar",
				Dfdaemon:   true,
				Path:       "/peer/file/taskFileName",
				PieceSize:  4 * 1024 * 1024,
				TaskID:     "test2",
				PeerID:     "peer2",
			},
		},
	}

	for _, tc := range testCases {
		err := manager.Add(context.Background(), tc.dfgetTask)
		c.Check(err, check.IsNil)

		err = manager.Delete(context.Background(), tc.dfgetTask.CID, tc.dfgetTask.TaskID)
		c.Check(err, check.IsNil)
		if tc.dfgetTask.Dfdaemon {
			c.Assert(0, check.Equals,
				int(prom_testutil.ToFloat64(
					dfgetDaemonNum.WithLabelValues(tc.dfgetTask.TaskID, tc.dfgetTask.CallSystem))))
		} else {
			c.Assert(0, check.Equals,
				int(prom_testutil.ToFloat64(
					dfgetNum.WithLabelValues(tc.dfgetTask.TaskID, tc.dfgetTask.CallSystem))))
		}

		_, err = manager.Get(context.Background(), tc.dfgetTask.CID, tc.dfgetTask.TaskID)
		c.Check(errors.IsDataNotFound(err), check.Equals, true)
	}
>>>>>>> 40bbeb8... add some supernode metrics
}

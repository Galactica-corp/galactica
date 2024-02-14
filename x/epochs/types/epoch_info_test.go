// Copyright 2024 Galactica Network
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

package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type EpochInfoTestSuite struct {
	suite.Suite
}

func TestEpochInfoTestSuite(t *testing.T) {
	suite.Run(t, new(EpochInfoTestSuite))
}

func (suite *EpochInfoTestSuite) TestStartEndEpoch() {
	startTime := time.Now()
	duration := time.Hour * 24
	ei := EpochInfo{StartTime: startTime, Duration: duration}

	ei.StartInitialEpoch()
	suite.Require().True(ei.EpochCountingStarted)
	suite.Require().Equal(int64(1), ei.CurrentEpoch)
	suite.Require().Equal(startTime, ei.CurrentEpochStartTime)

	ei.EndEpoch()
	suite.Require().Equal(int64(2), ei.CurrentEpoch)
	suite.Require().Equal(startTime.Add(duration), ei.CurrentEpochStartTime)
}

func (suite *EpochInfoTestSuite) TestValidateEpochInfo() {
	testCases := []struct {
		name       string
		ei         EpochInfo
		expectPass bool
	}{
		{
			"invalid - blank identifier",
			EpochInfo{
				"  ",
				time.Now(),
				time.Hour * 24,
				1,
				time.Now(),
				true,
				1,
			},
			false,
		},
		{
			"invalid - epoch duration zero",
			EpochInfo{
				WeekEpochID,
				time.Now(),
				time.Hour * 0,
				1,
				time.Now(),
				true,
				1,
			},
			false,
		},
		{
			"invalid - negative current epoch",
			EpochInfo{
				WeekEpochID,
				time.Now(),
				time.Hour * 24,
				-1,
				time.Now(),
				true,
				1,
			},
			false,
		},
		{
			"invalid - negative epoch start height",
			EpochInfo{
				WeekEpochID,
				time.Now(),
				time.Hour * 24,
				1,
				time.Now(),
				true,
				-1,
			},
			false,
		},
		{
			"pass",
			EpochInfo{
				WeekEpochID,
				time.Now(),
				time.Hour * 24,
				1,
				time.Now(),
				true,
				1,
			},
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.ei.Validate()

		if tc.expectPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

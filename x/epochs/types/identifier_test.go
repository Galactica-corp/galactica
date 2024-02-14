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

	"github.com/stretchr/testify/suite"
)

type EpochIdentifierTestSuite struct {
	suite.Suite
}

func TestEpochIdentifierTestSuite(t *testing.T) {
	suite.Run(t, new(EpochIdentifierTestSuite))
}

func (suite *EpochIdentifierTestSuite) TestValidateEpochIdentifierInterface() {
	testCases := []struct {
		name       string
		id         interface{}
		expectPass bool
	}{
		{
			"invalid - blank identifier",
			"",
			false,
		},
		{
			"invalid - blank identifier with spaces",
			"   ",
			false,
		},
		{
			"invalid - non-string",
			3,
			false,
		},
		{
			"pass",
			WeekEpochID,
			true,
		},
	}

	for _, tc := range testCases {
		err := ValidateEpochIdentifierInterface(tc.id)

		if tc.expectPass {
			suite.Require().NoError(err, tc.name)
		} else {
			suite.Require().Error(err, tc.name)
		}
	}
}

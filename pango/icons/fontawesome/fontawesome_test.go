// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fontawesome

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/soumya92/barista/pango"
	"github.com/soumya92/barista/testing/cron"
	"github.com/soumya92/barista/testing/githubfs"
	pangoTesting "github.com/soumya92/barista/testing/pango"
)

func TestInvalid(t *testing.T) {
	fs = afero.NewMemMapFs()
	require.Error(t, Load("/src/no-such-directory"))

	afero.WriteFile(fs, "/src/fa-error-1/web-fonts-with-css/scss/_variables.scss", []byte(
		`
$blah: "red";
$fa-var-xy: xy;
$fa-var-foobar: \61;
`,
	), 0644)
	require.Error(t, Load("/src/fa-error-1"))
}

func TestValid(t *testing.T) {
	fs = afero.NewMemMapFs()
	afero.WriteFile(fs, "/src/fa/web-fonts-with-css/scss/_variables.scss", []byte(
		`
$fa-var-some-icon: \61;
$fa-var-other-icon: \62;
`,
	), 0644)
	require.NoError(t, Load("/src/fa"))
	pangoTesting.AssertText(t, "a", pango.Icon("fa-some-icon").String())
	pangoTesting.AssertText(t, "b", pango.Icon("fa-other-icon").String())
}

// TestLive tests that current master branch of the icon font works with
// this package. This test only runs when CI runs tests in 'cron' mode,
// which provides timely notifications of incompatible changes while
// keeping default tests hermetic.
func TestLive(t *testing.T) {
	fs = githubfs.New()
	cron.Test(t, func() error {
		if err := Load("/FortAwesome/Font-Awesome/master"); err != nil {
			return err
		}
		// At least one of these icons should be loaded.
		testIcons := pango.New(
			pango.Icon("fa-arrow-circle-right"),
			pango.Icon("fa-cloud"),
			pango.Icon("fa-music"),
			pango.Icon("fa-tags"),
		)
		require.NotEmpty(t, testIcons.String(), "No expected icons were loaded")
		return nil
	})
}

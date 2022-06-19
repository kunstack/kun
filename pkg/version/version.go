/*
 * Copyright 2021 Aapeli.Smith<aapeli.nian@gmail.com>.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package version

import (
	"runtime"
	"time"
)

var (
	// Semver holds the current version of kun.
	Semver = "dev"
	// BuildDate holds the build date of kun.
	BuildDate = "I don't remember exactly"
	// StartDate holds the start date of kun.
	StartDate = time.Now()
	// GitCommit The commit ID of the current commit.
	GitCommit = "I don't remember exactly"
)

// Version describes compile time information.
type Version struct {
	// Semver is the current semver.
	Semver string `json:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `json:"git_commit,omitempty"`
	// BuildDate holds the build date of this component.
	BuildDate string `json:"build_date,omitempty"`
	// StartDate holds the start date of this component.
	StartDate time.Time `json:"startDate,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"go_version,omitempty"`
}

// Get returns build info
func Get() *Version {
	return &Version{
		Semver:    Semver,
		GitCommit: GitCommit,
		StartDate: StartDate,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
	}
}

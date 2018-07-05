// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package linux

import (
	"bytes"
	"strconv"

	"github.com/elastic/go-sysinfo/types"
)

type SeccompMode uint8

const (
	SeccompModeDisabled SeccompMode = iota
	SeccompModeStrict
	SeccompModeFilter
)

func (m SeccompMode) String() string {
	switch m {
	case SeccompModeDisabled:
		return "disabled"
	case SeccompModeStrict:
		return "strict"
	case SeccompModeFilter:
		return "filter"
	default:
		return strconv.Itoa(int(m))
	}
}

func readSeccompFields(content []byte) (*types.SeccompInfo, error) {
	var seccomp types.SeccompInfo

	err := parseKeyValue(bytes.NewReader(content), ":", func(key, value string) error {
		switch key {
		case "Seccomp":
			mode, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return err
			}
			seccomp.Mode = SeccompMode(mode).String()
		case "NoNewPrivs":
			noNewPrivs, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			seccomp.NoNewPrivs = &noNewPrivs
		}
		return nil
	})

	return &seccomp, err
}

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
	"io"
	"time"

	"github.com/elastic/go-sysinfo/providers/linux/meminfo"
	"github.com/elastic/go-sysinfo/types"
)

func parseHostMemoryInfo(r io.Reader) (*types.HostMemoryInfo, error) {
	var memInfo meminfo.Info
	metrics := make(map[string]uint64)
	if err := meminfo.Parse(r, &memInfo, metrics); err != nil {
		return nil, err
	}
	return &types.HostMemoryInfo{
		Timestamp:    time.Now().UTC(),
		Total:        memInfo.Total,
		Used:         memInfo.Total - memInfo.Free,
		Available:    memInfo.Available,
		Free:         memInfo.Free,
		VirtualTotal: memInfo.VirtualTotal,
		VirtualUsed:  memInfo.VirtualTotal - memInfo.VirtualFree,
		VirtualFree:  memInfo.VirtualFree,
		Metrics:      metrics,
	}, nil
}

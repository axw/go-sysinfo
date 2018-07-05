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

package meminfo

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Info holds selected fields parsed from /proc/meminfo.
type Info struct {
	Total        uint64
	Free         uint64
	Available    uint64
	VirtualTotal uint64
	VirtualFree  uint64
}

// Parse parses /proc/meminfo into memInfo, and all other fields into other if non-nil.
func Parse(r io.Reader, memInfo *Info, other map[string]uint64) error {
	var buffers, cached uint64
	var hasAvailable bool

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i == -1 {
			continue
		}
		key := line[:i]
		value := strings.TrimSpace(line[i+1:])

		var field *uint64
		switch key {
		case "MemTotal":
			field = &memInfo.Total
		case "MemAvailable":
			field = &memInfo.Available
			hasAvailable = true
		case "MemFree":
			field = &memInfo.Free
		case "SwapTotal":
			field = &memInfo.VirtualTotal
		case "SwapFree":
			field = &memInfo.VirtualFree
		case "Buffers":
			field = &buffers
		case "Cached":
			field = &cached
		default:
			if other == nil {
				continue
			}
		}
		num, err := parseBytesOrNumber(value)
		if err != nil {
			return errors.Wrapf(err, "failed to parse %v value of %v", key, value)
		}
		if field != nil {
			*field = num
		} else {
			other[key] = num
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if other != nil {
		other["Buffers"] = buffers
		other["Cached"] = cached
	}
	if !hasAvailable {
		// MemAvailable was added in kernel 3.14.
		//
		// Linux uses this for the calculation (but we are using a simpler calculation):
		// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=34e431b0ae398fc54ea69ff85ec700722c9da773
		memInfo.Available = memInfo.Free + buffers + cached
	}
	return nil
}

func parseBytesOrNumber(data string) (uint64, error) {
	parts := strings.Fields(data)
	if len(parts) == 0 {
		return 0, errors.New("empty value")
	}

	num, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse value")
	}

	var multiplier uint64 = 1
	if len(parts) >= 2 {
		switch parts[1] {
		case "kB":
			multiplier = 1024
		default:
			return 0, errors.Errorf("unhandled unit %v", parts[1])
		}
	}

	return num * multiplier, nil
}

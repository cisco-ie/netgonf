/**
 * Copyright (c) 2019-2020 Cisco Systems
 *
 * Author: Steven Barth <stbarth@cisco.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package netconf

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

// ErrFraming describes a NETCONF protocol error due to invalid message framing
var ErrFraming = errors.New("NETCONF message framing error")

// NETCONF 1.0 message delimiter sequence
var eom = []byte{']', ']', '>', ']', ']', '>'}

type framerV10 struct {
	writer io.Writer
}

func newFramerV10(writer io.Writer) io.WriteCloser {
	return &framerV10{writer: writer}
}

func (c *framerV10) Write(p []byte) (int, error) {
	return c.writer.Write(p)
}

func (c *framerV10) Close() error {
	_, err := c.writer.Write(eom)
	return err
}

type unframerV10 struct {
	reader io.Reader
	buffer []byte
	len    int
	err    error
}

func newUnframerV10(reader io.Reader) io.ReadCloser {
	return &unframerV10{reader: reader, buffer: make([]byte, len(eom))}
}

func (c *unframerV10) Read(p []byte) (int, error) {
	if c.err != nil {
		return 0, c.err
	}

	for c.len < len(c.buffer) {
		n, err := c.reader.Read(c.buffer[c.len:])
		if err != nil {
			if err == io.EOF {
				err = ErrFraming
			}
			c.err = err
			return 0, c.err
		}
		c.len += n
	}

	var i int
	for i = 0; i < len(c.buffer); i++ {
		if bytes.Equal(c.buffer[i:], eom[:len(eom)-i]) {
			break
		}
	}

	if i == 0 {
		c.err = io.EOF
		return 0, c.err
	}

	len := copy(p, c.buffer[:i])
	c.len = copy(c.buffer, c.buffer[len:])

	return len, nil
}

func (c *unframerV10) Close() error {
	dummy := make([]byte, 16)
	for {
		_, err := c.Read(dummy)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

type framerV11 struct {
	writer io.Writer
}

func newFramerV11(writer io.Writer) io.WriteCloser {
	return &framerV11{writer: writer}
}

func (c *framerV11) Write(p []byte) (int, error) {
	if len(p) > 0 {
		_, err := c.writer.Write([]byte("\n#" + strconv.Itoa(len(p)) + "\n"))
		if err != nil {
			return 0, err
		}
	}
	return c.writer.Write(p)
}

func (c *framerV11) Close() error {
	_, err := c.writer.Write([]byte{'\n', '#', '#', '\n'})
	return err
}

type unframerV11 struct {
	reader io.Reader
	len    int
	err    error
}

func newUnframerV11(reader io.Reader) io.ReadCloser {
	return &unframerV11{reader: reader}
}

func (c *unframerV11) Read(p []byte) (int, error) {
	if c.err != nil || len(p) == 0 {
		return 0, c.err
	} else if c.len == 0 {
		_, err := c.reader.Read(p[:1])
		if err == nil && p[0] == '\n' {
			_, err = c.reader.Read(p[:1])
			if (err == nil && p[0] != '#') || err == io.EOF {
				err = ErrFraming
			}
		} else if err == nil || err == io.EOF {
			err = ErrFraming
		}

		for err == nil {
			_, err = c.reader.Read(p[:1])
			if err == nil {
				if c.len == 0 && p[0] == '#' {
					_, err = c.reader.Read(p[:1])
					if err == nil && p[0] == '\n' {
						err = io.EOF
					} else if err == nil || err == io.EOF {
						err = ErrFraming
					}
				} else if p[0] >= '0' && p[0] <= '9' {
					c.len = c.len*10 + int(p[0]-'0')
				} else if p[0] == '\n' && c.len > 0 {
					break
				} else {
					err = ErrFraming
				}
			} else if err == io.EOF {
				err = ErrFraming
			}
		}

		if err != nil {
			c.err = err
			return 0, c.err
		}
	}

	if c.len < len(p) {
		p = p[:c.len]
	}

	n, err := c.reader.Read(p)
	if err != nil {
		if err == io.EOF {
			err = ErrFraming
		}
		c.err = err
		return 0, c.err
	}

	c.len -= n
	return n, nil
}

func (c *unframerV11) Close() error {
	dummy := make([]byte, 16)
	for {
		_, err := c.Read(dummy)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

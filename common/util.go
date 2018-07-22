/*
 * Copyright 2018 Anoop Vijayan Maniankara
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package common

import (
	"bufio"
	"io"
	"os"

	pb "github.com/maniankara/gcpd/gcp"
)

// SendToStream The implementor implements how to send to stream
// E.g.
// var stream pb.Gcp_CopyFromServer // or
// var stream pb.Gcp_CopyToClient
//
// fn := func (chunk *pb.CopyFile) error {
// 		return stream.Send(chunk)
// }
type SendToStream func(*pb.CopyFile) error

// GetFromStream The implementor implements how its read from stream
// e.g.
// var stream pb.Gcp_CopyFromClient // or
// var stream pb.Gcp_CopyToServer
//
// fn := func() (int, error) {
// 		return stream.Recv()
// }
type GetFromStream func() (*pb.CopyFile, error)

// WriteToStream the contents of the file (given in path)
// fn() defined whats needed to be done with the chunk
func WriteToStream(path string, fn SendToStream) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	for {
		buf := make([]byte, 1024)
		count, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		chunk := &pb.CopyFile{
			Chunktype: &pb.CopyFile_FileData{
				FileData: buf[:count],
			},
		}
		if err := fn(chunk); err != nil {
			return err
		}
		buf = nil // empty the buffer
	}
	return nil

}

// ReadFromStream the contents and write to path
// fn() defines what needed to be done with the chunk
func ReadFromStream(path string, fn GetFromStream) error {

	fh, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		return err
	}

	for {
		c, err := fn()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		buf := c.GetFileData()
		if _, err := fh.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

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

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	util "github.com/maniankara/gcpd/common"
	pb "github.com/maniankara/gcpd/gcp"
	grpc "google.golang.org/grpc"
)

type gcpServer struct{}

func main() {
	var port = flag.String("p", "10001", "Port to listen to")
	flag.Parse()

	conn, err := net.Listen("tcp", fmt.Sprintf(":"+*port))
	if err != nil {
		log.Fatal("Unable to serve on port: ", err)
	}
	log.Println("Serving on port: ", *port)

	grpcServer := grpc.NewServer()
	pb.RegisterGcpServer(grpcServer, &gcpServer{})
	grpcServer.Serve(conn)
}

// Copy to server
func (s *gcpServer) CopyTo(stream pb.Gcp_CopyToServer) error {

	var fh *os.File
	var path string

	for {
		c, err := stream.Recv()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		switch c.Chunktype.(type) {
		case *pb.CopyFile_FilePath:
			path = c.GetFilePath()
		case *pb.CopyFile_FileData:
			if fh == nil {
				fh, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0700)
				if err != nil || fh == nil {
					return err
				}
			}
			buf := c.GetFileData()
			if _, err := fh.Write(buf); err != nil {
				return err
			}
		}
	}
	defer fh.Close()
	return stream.SendAndClose(&pb.TransferPath{Path: path})
}

// Copy from server
func (s *gcpServer) CopyFrom(path *pb.TransferPath, stream pb.Gcp_CopyFromServer) error {

	fn := func(chunk *pb.CopyFile) error {
		return stream.Send(chunk)
	}

	err := util.WriteToStream(path.GetPath(), fn)
	if err != nil {
		return err
	}
	return nil
}

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
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	util "github.com/maniankara/gcpd/common"
	pb "github.com/maniankara/gcpd/gcp"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func usage() {
	progName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Syntax:\n")
	fmt.Fprintf(os.Stderr, "./%s [-p port] host:filepath path # Copy from remote host\n", progName)
	fmt.Fprintf(os.Stderr, "./%s [-p port] path host:filepath # Copy to remote host\n", progName)
	fmt.Fprintf(os.Stderr, "E.g.\n")
	fmt.Fprintf(os.Stderr, "./%s dhcp-101:/var/opt/y.iso /var/tmp/y.iso\n", progName)
	fmt.Fprintf(os.Stderr, "./%s /var/tmp/x.iso 35.128.27.105:/var/opt/x.iso\n", progName)
	fmt.Fprintf(os.Stderr, "./%s -p 10010 /var/tmp/x.iso 35.128.27.105:/var/opt/x.iso\n", progName)
}

func main() {

	// Handle command line args
	var serverPort = flag.String("p", "10001", "Port to use")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		os.Exit(-1)
	}

	params, copyFrom, err := parseCli()
	if err != nil {
		log.Println(err)
		usage()
		os.Exit(-1)
	}

	// dial grpc
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(params[0]+":"+*serverPort, opts...)
	if err != nil {
		log.Fatal("Error while connecting: ", err)
	}
	defer conn.Close()

	// create client
	client := pb.NewGcpClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// invoke remote method
	if copyFrom {
		copyFromServer(ctx, client, params[1], params[2])
	} else {
		copyToServer(ctx, client, params[1], params[2])
	}
}

func parseCli() ([]string, bool, error) {

	var (
		params   []string // hostname, remotepath, localpath
		copyFrom bool
		err      error
	)

	arg1 := flag.Arg(0)
	arg2 := flag.Arg(1)
	if strings.Contains(arg1, ":") {
		// from remote server
		params = strings.Split(arg1, ":")
		params = append(params, arg2)
		copyFrom = true
	} else if strings.Contains(arg2, ":") {
		// to remote server
		params = strings.Split(arg2, ":")
		params = append(params, arg1)
		copyFrom = false
	} else {
		// wrong syntax
		err = errors.New("None of the params define a host")
	}
	return params, copyFrom, err
}

// Copies file as a stream from client to server
func copyToServer(ctx context.Context, client pb.GcpClient, rPath string, lPath string) {

	stream, err := client.CopyTo(ctx)
	if err != nil {
		log.Fatal("Failed to call CopyTo rpc", err)
	}

	// Send the file path as first chunk
	pChunk := &pb.CopyFile{
		Chunktype: &pb.CopyFile_FilePath{FilePath: rPath},
	}
	stream.Send(pChunk)

	// Send the file contents to stream
	fn := func(chunk *pb.CopyFile) error {
		return stream.Send(chunk)
	}

	err = util.WriteToStream(lPath, fn)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Unsucessful closing of stream: ", err)
	}

}

// Copies file as a stream from server to client
func copyFromServer(ctx context.Context, client pb.GcpClient, rPath string, lPath string) {

	tPath := &pb.TransferPath{Path: rPath}
	stream, err := client.CopyFrom(ctx, tPath)
	if err != nil {
		log.Fatal("Unable to connect to stream", err)
	}

	fn := func() (*pb.CopyFile, error) {
		return stream.Recv()
	}

	err = util.ReadFromStream(lPath, fn)
	if err != nil {
		log.Fatal(err)
	}

}

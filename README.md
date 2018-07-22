# GCP/GCPD - A gRPC implementation of scp
A simple secure file copy implementation in gRPC.

## How it works
![GCPD](/gcpd.png)


## Reference/Summary of 4 types gRPC implementations

| Proto                                            | Client side                         | Server side                                  |
|:------------------------------------------------:|:-----------------------------------:|:--------------------------------------------:|
| rpc FUNC(TYPE1) returns (TYPE2) {};              | obvious                             | obvious                                      | 
| rpc FUNC(stream TYPE1) returns (TYPE2) {};       | stream.Recv() &&                    | stream.Send(TYPE1) &&                        |
|                                                  |        stream.SendAndClose(TYPE2)   |      stream.CloseAndRecv()                   |
| rpc FUNC(TYPE1) returns (stream TYPE2) {};       | stream.Send(TYPE2)                  | stream.Recv() && EOF                         |
| rpc FUNC(stream TYPE1) returns (stream TYPE2) {};| stream.Recv() && EOF                | stream.Recv() && EOF &&                      |
|                                                  |        stream.Send(TYPE2)           |      stream.Send(TYPE1) && stream.CloseSend()| 


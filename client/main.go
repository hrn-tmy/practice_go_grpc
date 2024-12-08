package main

import (
	"context"
	"fmt"
	"go-grpc/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	// "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	certFile := "/Users/tomoya/Library/Application Support/mkcert/rootCA.pem"
	creds, _ := credentials.NewClientTLSFromFile(certFile, "")
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
	// 「grpc.dial()」と「grpc.withInsecure()」は非推奨になったのでコメントアウト
	// conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	// callListFiles(client)
	callDownload(client)
	// callUpload(client)
	// calluploadAndNotifyProgress(client)
}

func callListFiles(client pb.FileServiceClient) {
	md := metadata.New(map[string]string{"authorization": "Bearer test-token"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	res, err := client.ListFiles(ctx, &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res.GetFilenames())
}

func callDownload(client pb.FileServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	req := &pb.DownloadRequest{Filemane: "name.txt"}
	stream, err := client.Download(ctx, req)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			resErr, ok := status.FromError(err)
			if ok {
				if resErr.Code() == codes.NotFound {
					log.Fatalf("error code: %v, Error Message: %v", resErr.Code(), resErr.Message())
				} else if resErr.Code() == codes.DeadlineExceeded {
					log.Fatalln("deadline exceeded")
				} else {
					log.Fatalln("unknown grpc error")
				}
			} else {
				log.Fatalln(err)
			}
		}

		log.Printf("Response from download(bytes): %v", res.GetData())
		log.Printf("Response from download(string): %v", string(res.GetData()))
	}
}

func callUpload(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/Users/tomoya/Downloads/CreatedApp/go-grpc/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		req := &pb.UploadReqest{Data: buf[:n]}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalln(sendErr)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("received data size: %v", res.GetSize())
}

func calluploadAndNotifyProgress(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/Users/tomoya/Downloads/CreatedApp/go-grpc/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.UploadAndNotifyProgress(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// リクエスト側のゴルーチン
	buf := make([]byte, 5)
	go func() {
		for {
			n, err := file.Read(buf)
			if n == 0 || err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}
			req := &pb.UploadAndNotifyProgressRequest{Data: buf[:n]}
			sendErr := stream.Send(req)
			if sendErr != nil {
				log.Fatalln(sendErr)
			}
			time.Sleep(1 * time.Second)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// レスポンス側のゴルーチン
	ch := make(chan struct{})
	go func ()  {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("received message: %v", res.GetMsg())
		}
		close(ch)
	}()
	<-ch
}

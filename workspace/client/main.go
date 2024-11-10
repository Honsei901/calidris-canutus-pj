package main

import (
	"context"
	"fmt"
	"grpc/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	callListFiles(client)
	callDownload(client)
	callUpload(client)
	callUploadAndNotifyProgress(client)
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
	req := &pb.DownloadRequest{
		Filename: "language.txt",
	}

	stream, err := client.Download(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln(err)
		}

		log.Println(string(res.GetData()))
	}
}

func callUpload(client pb.FileServiceClient) {
	filename := "name.txt"
	path := "D:/code-study/repository/github/for-public/training/web/calidris-canutus-pj/workspace/storage/" + filename

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

		req := &pb.UploadRequest{
			Data: buf[:n],
		}

		if sendErr := stream.Send(req); sendErr != nil {
			log.Fatalln(sendErr)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Received data size: %v", res.GetSize())
}

func callUploadAndNotifyProgress(client pb.FileServiceClient) {
	filename := "language.txt"
	path := "D:/code-study/repository/github/for-public/training/web/calidris-canutus-pj/workspace/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	stream, err := client.UploadAndNotifyProgress(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// Request
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

			req := &pb.UploadAndNotifyProgressRequest{
				Data: buf[:n],
			}

			if sendErr := stream.Send(req); sendErr != nil {
				log.Fatalln(sendErr)
			}
			time.Sleep(1 * time.Second)
		}

		if err := stream.CloseSend(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Response
	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("Received message: %v", res.GetMsg())
		}
		close(ch)
	}()
	<-ch
}

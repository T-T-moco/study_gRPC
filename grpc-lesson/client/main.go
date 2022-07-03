package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:30000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	// callListFiles(client) // 普通の関数を呼び出す
	// callDownload(client) // メイン関数からサーバーストリーミングRPCを呼び出す
	CallUpload(client) // メイン関数からクライアントストリーミングRPCを呼び出す
}

func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFilenames())
}

/*****************************************************
* サーバーストリーミングRPC（１リクエストに複数レスポンス）
******************************************************/
// １つのリクエストを渡した後複数のレスポンスが帰ってくる。
// 1. クライアントからファイルをアップロード
// 2. バックエンドからファイルの中身を5バイトずつクライアントに返す
func callDownload(client pb.FileServiceClient) {
	req := &pb.DownloadRequest{Filename: "name.txt"}
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

		log.Printf("Response from Download(bytes): %v", res.GetData())
		log.Printf("Response from Download(string): %v", string(res.GetData()))
	}
}

/*****************************************************
* クライアントストリーミングRPC（複数リクエストに１レスポンス）
******************************************************/
// 複数のリクエストを渡し、渡し終わってからレスポンスを受け取る
// 1. 5バイトずつリクエストを送り、バックエンドで都度表示している。
// 2. 全てリクエストで送り終わってから、バックエンドから総サイズ数をクライアントに返す
func CallUpload(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/Users/takemasatomoko/Desktop/source/practice/gRPC/grpc-lesson/strage/" + filename // クライアントのストレージであると想定

	file, err := os.Open(path) // ファイルを開く
	if err != nil {            // ファイルを開けなかったら
		log.Fatalln(err)
	}
	defer file.Close() // ファイルをクローズする

	stream, err := client.Upload(context.Background()) // ストリームが取得できる
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 5) // データ格納用のバッファを用意する
	for {
		n, err := file.Read(buf)     // ファイルの内容をバッファに格納し
		if n == 0 || err == io.EOF { // この場合はループを抜ける
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		req := &pb.UploadRequest{Data: buf[:n]} // 読み込んだ内容をリクエストに詰める
		sendErr := stream.Send(req)             // ストリームからリクエストを送信することができる
		if sendErr != nil {
			log.Fatalln(sendErr)
		}

		time.Sleep(1 * time.Second) // 1秒間スリープを入れる
	}

	res, err := stream.CloseAndRecv() // リクエストの終了をサーバーに通知し、レスポンスを受け取ることができる
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("received data size: %v", res.GetSize())
}

/*****************************************************
* 双方向ストリーミングRPC（複数リクエストに複数レスポンス）
******************************************************/
func CallUploadAndNotifyProgress(client pb.FileServiceClient) {
	// 処理の半分はCallUpload関数と同じなため、コピペ
	filename := "sports.txt"
	path := "/Users/takemasatomoko/Desktop/source/practice/gRPC/grpc-lesson/strage/" + filename // クライアントのストレージであると想定

	file, err := os.Open(path) // ファイルを開く
	if err != nil {            // ファイルを開けなかったら
		log.Fatalln(err)
	}
	defer file.Close() // ファイルをクローズする

	stream, err := client.UploadAndNotifyProgress(context.Background()) // ストリームを取得
	if err != nil {                                                     // エラーハンドリング
		log.Fatalln(err)
	}

	// 双方向ストリーミングRPCでは複数のリクエストを送信しつつ、複数のレスポンスも受け取る必要がある。
	// これを実現するためにボルーチン？を使用した平行処理を行う

	// request
	buf := make([]byte, 5) // データ格納用のバッファを用意
	go func() {
		for {
			n, err := file.Read(buf) // ファイルの内容をバッファに格納する
			if n == 0 || err == io.EOF {
				break // データの読み込み終了時にループを抜ける
			}
			if err != nil { // その他のエラーではerrを表示
				log.Fatalln(err)
			}

			req := &pb.UploadAndNotifyProgressRequest{Data: buf[:n]} // リクエストのメッセージを作成し
			sendErr := stream.Send(req)                              // ストリームを返してリクエストを送信する
			if sendErr != nil {                                      // エラーハンドリング
				log.Fatalln(err)
			}
			time.Sleep(1 * time.Second) // スリープ処理
		}

		err := stream.CloseSend() // リクエストの終了を通知
		// レスポンスの処理は別のボルーチン？で実装するのでCallUpload関数で使用したCloseAndRecvメソッドではなく、CloseSendメソッドを使用すること。
		if err != nil {
			log.Fatalln(err)
		}
	}()
}

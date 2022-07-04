package main

import (
	"bytes"
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

/*****************************************************
* Unary RPC（１リクエストに１レスポンス）
******************************************************/
type server struct {
	pb.UnimplementedFileServiceServer // エンベットする
}

/*
proto/file.protoの中で定義したrpc ListFilesの関数を記述する
func (*server) ListFiles() まで書いたら、pb/file_grpc.pb.goで[type FileServiceServer interface]のシグネチャを確認してくる（引数など何を入れたらいいか全て書いてある）
>> ListFiles(context.Context, *ListFilesRequest) (*ListFilesResponse, error)
※ １つめの塊が引数、２つめの塊が戻り値
*/
func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	fmt.Println("ListFiles was invoked") // メソッドがコールされたことをわかるようにする

	// strageに作ったtxtの中身を読み込む
	dir := "/Users/takemasatomoko/Desktop/source/practice/gRPC/grpc-lesson/strage"

	// 引数にdirを渡し、strage配下のパスを取得する
	paths, err := ioutil.ReadDir(dir)
	if err != nil { // エラーであればnilとerrを返す
		return nil, err
	}

	var filenames []string // ストリングのスライスでファイル名一覧を格納する変数を定義
	for _, path := range paths {
		if !path.IsDir() { // ；パスがファイルである場合
			filenames = append(filenames, path.Name()) // スライスにファイル名を配列に入れ込む
		}
	}

	res := &pb.ListFilesResponse{ // 戻り値のメッセージを作成
		Filenames: filenames,
	}
	return res, nil
}

/*****************************************************
* サーバーストリーミングRPC（１リクエストに複数レスポンス）
******************************************************/
// サーバーストリーミングRPC（１リクエストに複数レスポンス）
func (*server) Download(req *pb.DownloadRequest, stream pb.FileService_DownloadServer) error {
	fmt.Println("Download was invoked")

	filename := req.GetFilename()
	path := "/Users/takemasatomoko/Desktop/source/practice/gRPC/grpc-lesson/strage/" + filename

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		res := &pb.DownloadResponse{Data: buf[:n]}
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

/*****************************************************
* クライアントストリーミングRPC（複数リクエストに１レスポンス）
******************************************************/
// pb/file.pb.goのtype FileServiceServer interface の中にある Uploadメソットを参照。引数にFileService_UploadServerが必要 / 戻り値はerrorであることがわかる
func (*server) Upload(stream pb.FileService_UploadServer) error { // <- 戻り値はerror
	fmt.Println("Upload was invoked")

	var buf bytes.Buffer // クライアントからアップロードされたバッファを格納するための変数を用意
	for {
		req, err := stream.Recv() // クライアントからストリーム経由で複数のリクエストを受け取る
		if err == io.EOF {        // クライアントからの終了信号が到達した場合
			res := &pb.UploadResponse{Size: int32(buf.Len())} // int32にキャストしたバッファのサイズをメッセージに含め
			return stream.SendAndClose(res)                   // 引数にresを渡し、サーバーからのレスポンスを返す
		}
		if err != nil { // その他のエラーの場合
			return err
		}

		data := req.GetData() // リクエストからのデータを変数に格納し、その内容を出力する
		log.Printf("received data(bytes): %v", data)
		log.Printf("received data(string): %v", string(data))
		buf.Write(data) // データをバッファに書き込む
	}
}

/*****************************************************
* 双方向ストリーミングRPC（複数リクエストに複数レスポンス）
******************************************************/
func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress was invoked")

	size := 0 // 受信したサイズのデータを格納する変数を用意

	for {
		req, err := stream.Recv() // クライアントからストリーム経由で複数のリクエストを受け取れるようにする
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		data := req.GetData() //リクエストからデータを取り出し
		log.Printf("received data: %v", data)
		size += len(data) // クライアントから受信したデータのサイズをsize変数に足し合わせる

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("received %vbytes", size), // ここまでに受信したデータのサイズをメッセージとしてレスポンスに入れる
		}
		err = stream.Send(res) // レスポンスを返却
		if err != nil {
			return err // エラーがあればエラーをレスポンス
		}
	}
}

/*****************************************************
* main関数
******************************************************/
func main() {
	lis, err := net.Listen("tcp", "localhost: 30000") // 第一引数に通信プロトコル、第二引数にアドレスを指定（アドレスにはポート番号まで記載）
	if err != nil {                                   // エラーハンドリング
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer() // gRPCのサーバー構造体を取得する
	// NewServerメソッドはgrpcライブラリで定義されていて、サーバーのポインタ型を返す
	pb.RegisterFileServiceServer(s, &server{}) // pb/file_grpc.pb.goに定義されており、第一引数にgrpcサーバーを、第二引数にFileServiceServerインターフェイスを実装した構造体を渡すことで、grpcサーバーに構造体の内容を登録する
	// --> つまりgrpcサーバーが、ListFilesなどのメソッドを提供できるようになる

	fmt.Println("server is running...")
	if err := s.Serve(lis); err != nil { // 指定したリッスンポートでサーバーを起動できる
		log.Fatalf("Failed to serve: %v", err)
	}
}

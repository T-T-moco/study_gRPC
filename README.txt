/********************************
* RPCとは
*********************************/
＊RPCとは：Remote procedure Call
・Remote = 遠隔地（リモート）サーバーの
・Procedure = 手続き（メソッド）を
・Call = 呼び出す（実行する）
・ネットワーク上の他の端末と通信するための仕組み
・REST APIのようにパスやメソッドを指定する必要はなく、メソッド名と引数を指定する
・gRPC以外にもJSON-RPCなどがあるが、今はgRPCがデファクトスタンダード

【gRPCの特徴】
・データフォーマットにProtocol Buffersを使用
　・バイナリにシリアライズすることで送信データ量が減り高速な通信を実現
　・型付けされたデータ転送が可能
・IDL（Protocol Buffers）からサーバー側・クライアント側に必要なソースコードを生成
・通信にはHTTP/2を使用。これにより高速化が可能
・特定の言語やプラットフォームに依存しない

【gRPCが適したケース】
・Microservice間の通信
　・複数の言語やプラットフォームで構成される可能性がある
　・バックエンド間であればgRPCの恩恵が多く得られる
・モバイルユーザーが利用するサービス
　・通信量が削減できるため、通信容量制限にかかりにくい
・速度が求められる場合

【gRPCの開発の流れ】
1. protoファイルの生成
　messageやserviceを定義していく
2. protoファイルをコンパイルしてサーバー・クライアントの雛形コードを作成
　コンパイル時のオプションによって、何の言語にコンパイルするのか？を決める
3. 雛形コードを使用してサーバー・クライアントを実装

※gRPCの「g」はGoogleのgではない
・バージョンごとに意味が異なる

【HTTP/1について】
＊HTTP/1.1の課題
　・リクエストの多重化
　　・1リクエストに対して1レスポンスという制約があり、大量のリソースで構成されているページを表示するには大きなネックになる
　・プロトコルオーバーヘッド
　　・Cookieやトークンなどを毎回リクエストヘッダに付与してリクエストするため、オーバーヘッドが大きくなる

＊HTTP/2の特徴
　・ストリームという概念を導入
　　・1つのTCP接続を用いて、複数のリクエスト/レスポンスのやり取りが可能
　　・TCP接続を減らすことができるので、サーバーの負荷軽減が可能
　・ヘッダーの圧縮
　　・ヘッダーをHPACKという圧縮方式で圧縮し、さらにキャッシュを行うことで差分のみを送受信することで効率化
　・サーバープッシュ機能
　　・クライアントからのリクエスト無しにサーバーからデータを送信できる
　　・事前に必要と思われるリソースを送信しておくことで、ラウンドトリップの回数を削減しリソース読み込みまでの時間を短縮

Demo
・http://www.http2demo.io


/********************************
* messageについて
*********************************/
https://developers.google.com/protocol-buffers/docs/proto3

＊コンパイル後にimportインポートエラーが起きていた時：
　　command: go mod tidy


-IPATH, --proto_path={PATH}
protoファイルのimport文のパスを特定する

・proto3とproto2には互換性がないため、基本的にはproto3を指定
・syntaxの記述を忘れるとデフォルトでproto2が使用されるので注意


// Scalar型
// 最初は大文字、単語１つくらいのシンプルな名前
message Employee {
  int32 id = 1;
  string name = 2;
  string email = 3;
  repeated string phone_number = 5; // repeated fiels: 配列のように複数の要素を含めることができる
  map<string, Company.Project> Project = 6; // map: key, valueのような連想配列的な。keyがString 
  // mapにはrepeatedをつけることはできない
  // keyにはstring, int, boolのいずれか
  oneof profile {
    string text = 7;
    Video video = 8;
  }
  date.Date birthday = 9;
  // ワンオブは、複数の型のどれか１つを値として持つフィールドを定義する時に使う
  // repeatedにすることはできない
}
// 列挙型
enum Occupation {
  OCCUPATION_UNKNOWN = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}
message Company {
  message Project {}
}
message Video {}



【デフォルト値】
・定義したmessageでデータをやり取りする際に、定義したフィールドがセットされていない場合、そのフィールドのデフォルト値が設定される
・デフォルト値は型ごとに決められている
  string: ''
  bytes: 空のbyte
  bool: false
  整数型: 浮動小数点: 0
  列挙型: タグ番号0の型
  repeated: 空のリスト


【messageとは】
・複数のフォールドを持つことができる型定義
  ・それぞれのフィールドはスカラ型もしくはコンポジット型
・各言語のコードとしてコンパイルした場合、構造体やクラスとして変換される
・１つのprotoファイルに複数のmessage型を定義することも可能

ex)
message Person {
  ※ フィールドの型  フィールド名 = タグ番号（順番）
  string name = 1;
  int32 id = 2;
  string email = 3;
}


[ Scalar型について ]
＊不動小数点（扱えるサイズが違う）
  ・double
  ・float
＊数値
  ・int32 / uint32 / sint32 などがある
  u = アンサインド 正の整数を使用する時
  s = サインド　負の整数を使用する時
  ※パフォーマンスを気にしない時は普通のint型でおk
＊fixed32 / 64
  とても大きな数値を使う時に使用。でも普通にint使えば大体おk
＊bool
  ・真偽値
＊string
  ・文字列
＊byte
  ・バイト型

[ Tag について ]
・Protocol Buffersではフィールドはフィールド名ではなく、タグ番号によって識別される
・そのため重複は許されず、一意である必要がある
・タグの最小値は1, 最大値は536,870,911
・19000 ~ 19999 はProtocol Buffersの予約番号のため使用不可
・1 ~ 15番までは1byteで表すことができるため、よく使うフィールドには1 ~ 15番を割り当てることでパフォーマンスが向上する
・タグは連番にする必要はなし
  ・なのであまり使わないフィールドはあえて16番以降を割り当てることも可能
・タグ番号を予約するなど、安全にProtocol Buffersを使用する方法も用意されている


【列挙型について】
enumキーワードを使用
・型は不要
・タグ番号0から開始する必要あり
・タグ番号0は慣例的にUNKNOWNとするのが一般的

ex)
enum Occupation {
  OCCUPATION_UNKNOWN = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}


【repeated fiels について】
enumキーワードを使用
・型は不要
・タグ番号0から開始する必要あり
・タグ番号0は慣例的にUNKNOWNとするのが一般的

ex)
enum Occupation {
  OCCUPATION_UNKNOWN = 0;
  ENGINEER = 1;
  DESIGNER = 2;
  MANAGER = 3;
}


/********************************
* Serviceについて
*********************************/
・RPC（メソッド）の実装単位
　・サービス内に定義するメソッドがエンドポイントになる
　・1サービス内に複数のメソッドを定義できる
・サービスを利用するにはサービス名、メソッド名、引数、戻り値を定義する必要がある
・コンパイルしてgoファイルに変換すると、インターフェイスとなる
　・アプリケーション側でこのインターフェイスを実装する

gRPCでは通信方式によってそのメソッドの定義方法も変わる

【gRPCの通信方法】
＊４種類の通信方式
1. Unary RPC
2. Server Streaming RPC
3. Client Streaming RPC
4. Bidrectional Streaming RPC


＊Unary RPC ================================================
・最もシンプルな通信方式で、1リクエスト1レスポンスの方式
・通常の関数コールのように扱うことができる
・用途：サーバー間通信のAPIなど

ex)
message SayHelloRequest {}
message SayHelloResponse {}
Service Greeter {
  rpc SayHello(SayHelloRequest) returns (SayHelloResponse);
  rpc メソッド名(リクエスト※引数) returns※戻り値 (レスポンス);
}


＊Server Streaming RPC ======================================
・1リクエストに対して複数レスポンスの方式
・クライアントからはいくつレスポンスがあるか分からないので、サーバーから送信完了の信号が送信されるまでストリームのメッセージを読み続ける
・用途：サーバーからのプッシュ通知など

ex)
message SayHelloRequest {}
message SayHelloResponse {}
Service Greeter {
  rpc SayHello(SayHelloRequest) returns (stream SayHelloResponse);
  rpc メソッド名(リクエスト※引数) returns※戻り値 (stream レスポンス);
}


＊Client Streaming RPC =======================================
・複数リクエストに対して1レスポンスの方式
・サーバーはクライアントからリクエスト完了の信号が送信されるまでストリームからメッセージを読み続け、レスポンスを返さない
  ・クライアントからの終了信号を待ち、終了信号が来たらレスポンスを返す
・用途：大きなファイルのアップロードなど

ex)
message SayHelloRequest {}
message SayHelloResponse {}
Service Greeter {
  rpc stream SayHello(SayHelloRequest) returns (SayHelloResponse);
  rpc stream メソッド名(リクエスト※引数) returns※戻り値 (レスポンス);
}


＊Bidrirectional Streaming RPC ===============================
・複数リクエスト・複数1レスポンスの方式
・クライアントとサーバーのストリームが独立しており、リクエストとレスポンスはどのような順序でもよい
・用途：チャットやオンライン対戦など

ex)
message SayHelloRequest {}
message SayHelloResponse {}
Service Greeter {
  rpc stream SayHello(SayHelloRequest) returns (stream SayHelloResponse);
  rpc stream メソッド名(リクエスト※引数) returns※戻り値 (stream レスポンス);
}




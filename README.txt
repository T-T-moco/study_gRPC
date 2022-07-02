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

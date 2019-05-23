# 課題

下記の要件を満たすキャンペーンのバナーを制御するクラス及びテストコードを実装してください。
また、注意した点や本問題における不安点などをREADME.mdの形で記載してください。

# 要件
* 以下から得意な言語を使用してください(ただ入社後はGo, PHPがメインの使用言語となります)
  * PHP, Go, Java
* 本試験では設計力と選択した言語での実装力を確認します
  * この機能のみの実装で構いませんが、一般的なWebアプリケーションに組み込み運用することを想定して実装してください
* この課題は**REST API**の実装を必要としません
  * すなわちControllerやViewなどの実装は不要です
* データレイヤの実装も不要です
  * もしデータレイヤへのインターフェースを考慮に入れる場合はこれらをStubし、Mockもしくはテストデータを返すような構造を考慮してください  
* 標準関数以外に、オープンソースライセンスのライブラリなどを利用できます
  * ただし、Webアプリケーションフレームワークの利用は評価が難しくなるためご遠慮ください
* テストコードも必要です

# 仕様
* バナーの表示期間は、キャンペーンによって異なりますので、バナー毎に設定できるようにしてください
* バナーの表示可否は、以下のように判断します
  1. 表示期間中は、バナーを表示できる
      * 表示期間は、年月日時分秒の精度を必要とします
      * タイムゾーンも考慮してください
  2. バナーの表示開始日時前でも、許可した IPアドレス `( 10.0.0.1, 10.0.0.2 ）` からのアクセスの場合はバナーを表示できる
  3. `i, ii` 以外のときは、バナーを表示できない
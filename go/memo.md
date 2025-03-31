# 勉強メモ

勉強するものリスト {
    ・ProgateでGit, Go, SQLあたり終わらせる
    ・Linuxのしくみ読み進める
    ・
}

参考になりそうなurl {
    https://zenn.dev/o_ga/books/dc6c7b055b65a3/viewer/chapter1
    https://zenn.dev/isawa/articles/069cff08d64904
    https://qiita.com/y_murotani/items/1c1a93fba22035b45d6a##GinWebFramework
}


## Go

フォーマット指定子 {
    %s 文字列
    %d 整数
    %f 浮動小数点(デフォは6桁)
    %.2f 小数点以下2個まで表示
    %x 16進数
    %X 16進数 大文字
    %t 真偽値 
    %v 変数のデフォルトフォーマット
    %% 変数の型を表示
}


## SQL

基本文法 {
    SELECT *
    が基本

    SELECT DISTINCT(カラム)
    で重複する要素を消せる

    SELECT SUM(カラム)
    で合計値が取得できる

    SELECT AVG(カラム)
    で平均値が取得できる


    FROM カテゴリー

    WHERE カラム LIKE "%名前の一部%";
    で名前の一部に指定した文字列が含まれるものを探し出せる
    (LIKE : LIKE演算子, % : ワイルドカード)

    WHERE は
    AND
    OR を使うと複数指定できる

    NOTを入れる位置は
    WHERE NOT (条件文)
    WHERE カラム IS NOT NULL

    ORDER BY カラム ASC/DESC
    昇順/降順で並び替え

    LIMIT 数字
    で表示件数を制限
}
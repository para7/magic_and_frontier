# 魔法と探索データパック

ジョブデータパック v2.x.x

開発中…

- そのままサーバーに投入できるよう、開発ツールとか色々入ってます。そのうち導入方法をドキュメント化します

- まだ色々と設置中なので、恐らく今後リンク先等を変更していきます

- なんかいい感じの名称案ください

# 現在のあそびかた

## 魔法について

明らかに変なやつがスポーンします。そいつが魔法書を落とします。

魔法は 魔法の杖と魔法書を両手それぞれに持って右クリックで発動します。

## 杖について

まだ杖をつくる方法がないので、コマンドで生成してください。

```
/summon item ~ ~ ~ {Item:{id:"minecraft:carrot_on_a_stick",Count:1b,tag:{display:{Name:'{"text":"魔法杖","color":"gold","italic":false}'},HideFlags:3,RepairCost:99,wandID:1,Enchantments:[{id:"minecraft:mending",lvl:1s}]}}}
```

## アナライズについて

アナライズの魔法は、インベントリ欄のもっとも左上の欄に解析したい魔法書を置いてから詠唱してください。

## 魔法書が直接ほしい

maf:tool/ 以下に魔法書の入手関数があります

```
/function maf:tool/029001
```

などで実行してください。

指定するのは詠唱 ID(castid)です。  
魔法 ID はデータパック外の \_tool/out/sample/maf.json にあります

# 導入方法

## 現行開発版

まるごと zip でダウンロードしてから、解凍します。

致命的なバグ等含む場合があります。

## リリース版

本画面 右側の Releases からダウンロードできます。

現行開発版と比べて安定

<https://github.com/para7/magic_and_frontier/releases/>

# マイグレーション

過去バージョンから更新して正常に動かない場合、次のコマンドを実行すれば動作するかもしれません（α 版につき保証はいたしかねます）

```
/function maf:reinstall
```

# 各種連絡・報告フォーム

バグ報告等は以下にお願いします

【一般向け】  
https://docs.google.com/forms/d/e/1FAIpQLSdutG5Q5O34SY3zoA6wShnRi0LcSfw72-UXy7nagEcP9JHbbQ/viewform?usp=sf_link

【github わかる人向け】  
ここの issue に直接立てて頂いて構いません。よろしくお願いします。

# 利用条件

- 利用報告等は任意です。  
  報告してくれると開発モチベーションが上がって機能が増えたりするかもしれません。

- このページへリンクを掲載して頂ければ、再配布や改造版の配布も許可します。  
  <https://github.com/para7/magic_and_frontier>

# 前提データパック

- p7BaseSystem : 自作のユーティリティパック

ほか赤石愛さんの [oh my dat](https://github.com/Ai-Akaishi/OhMyDat) や [ScoreToHealth](https://github.com/Ai-Akaishi/ScoreToHealth) を導入検討中

---

---

---

## 予定

- 2023 年内にマルチ検証開始したい　 → done. オープンベータもそのうちやりたい
- discord サーバーとか作る？　 → プライベートサーバーは作った
- 公式 wiki みたいなチュートリアルサイト作りたい
- 整ってきたら開発メンバーの募集・追加してみたり？（テキスト入力周りが大変なのと、グラフィック作成してくれる人いたらなぁ）

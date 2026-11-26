# DB定義・管理

- DDLやDMLを管理します。
- 番号の帯域で以下のルールとします。
  - `00 ~ 09`: 全環境で適用するDDLやDML（マスタデータ）
  - `10 ~ 99`: ローカル用のテストデータ　* 本番環境適用外

## 既存DBへの変更適用

DBイメージ内の初期化SQLは、新規DBの作成時にのみ実行されます。
既存DBの `report_finding` に検索用索引を追加する場合は、`middleware/db` を作業ディレクトリとして対象DBに接続し、次のSQLを一度適用してください。

```sql
SOURCE migrations/20261126_add_report_finding_project_date_index.sql;
```

`(report_date, project_id, data_source, score)` の既存の一意索引は保持します。
追加索引は `ALGORITHM=INPLACE, LOCK=NONE` を指定しますが、開始・終了時のメタデータロック待ちや、構築中のI/O・追加領域は発生し得ます。

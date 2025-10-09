-- init-scripts/01-init.sql

-- 1. アプリケーション用ユーザー作成
CREATE USER app_user WITH PASSWORD 'app_password';

-- 2. データベースへの接続権限のみ（最小限）
GRANT CONNECT ON DATABASE app_db TO app_user;

-- 3. app スキーマを作成（app_user が所有者）
CREATE SCHEMA IF NOT EXISTS app AUTHORIZATION app_user;

-- 4. app スキーマの使用権限
GRANT USAGE ON SCHEMA app TO app_user;

-- 5. app スキーマ内の既存テーブルへの権限
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA app TO app_user;

-- 6. app スキーマ内のシーケンス（自動採番）への権限
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA app TO app_user;

-- 7. 今後作成されるテーブルにも同じ権限を自動付与
ALTER DEFAULT PRIVILEGES FOR USER app_user IN SCHEMA app 
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_user;

-- 8. 今後作成されるシーケンスにも権限を自動付与
ALTER DEFAULT PRIVILEGES FOR USER app_user IN SCHEMA app 
GRANT USAGE, SELECT ON SEQUENCES TO app_user;

-- 9. public スキーマへのアクセスを明示的に拒否（セキュリティ強化）
REVOKE ALL ON SCHEMA public FROM app_user;
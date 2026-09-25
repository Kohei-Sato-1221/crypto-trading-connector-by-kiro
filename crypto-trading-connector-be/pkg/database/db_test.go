package database

import "testing"

func TestEnsureBinaryParameters(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		want string
	}{
		{
			name: "クエリパラメータ有りのURLに追加する",
			dsn:  "postgresql://user:pass@host.pooler.supabase.com:6543/postgres?sslmode=require",
			want: "postgresql://user:pass@host.pooler.supabase.com:6543/postgres?sslmode=require&binary_parameters=yes",
		},
		{
			name: "クエリパラメータ無しのURLに追加する",
			dsn:  "postgres://user:pass@host:5432/postgres",
			want: "postgres://user:pass@host:5432/postgres?binary_parameters=yes",
		},
		{
			name: "既に指定済みなら変更しない",
			dsn:  "postgresql://user:pass@host:6543/postgres?binary_parameters=no",
			want: "postgresql://user:pass@host:6543/postgres?binary_parameters=no",
		},
		{
			name: "key=value形式のDSNは変更しない",
			dsn:  "host=localhost port=5432 user=postgres sslmode=disable",
			want: "host=localhost port=5432 user=postgres sslmode=disable",
		},
		{
			name: "MySQLのDSNは変更しない",
			dsn:  "user:pass@tcp(localhost:3306)/db?parseTime=true",
			want: "user:pass@tcp(localhost:3306)/db?parseTime=true",
		},
		{
			name: "パスワードの特殊文字を壊さない",
			dsn:  "postgresql://user:p%40ss%26word@host:6543/postgres?sslmode=require",
			want: "postgresql://user:p%40ss%26word@host:6543/postgres?sslmode=require&binary_parameters=yes",
		},
		{
			name: "空文字列はそのまま",
			dsn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ensureBinaryParameters(tt.dsn); got != tt.want {
				t.Errorf("ensureBinaryParameters()\n got = %q\nwant = %q", got, tt.want)
			}
		})
	}
}

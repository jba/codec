module github.com/jba/codec/internal/benchmarks

go 1.15

replace github.com/jba/codec => ../..

require (
	cloud.google.com/go/storage v1.10.0
	github.com/GoogleCloudPlatform/cloudsql-proxy v1.18.0
	github.com/google/licensecheck v0.0.0-20200805042302-c54f297c3b57
	github.com/jackc/pgx/v4 v4.10.0
	github.com/jba/codec v0.0.0-00010101000000-000000000000
	github.com/ugorji/go/codec v1.2.2
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324
)

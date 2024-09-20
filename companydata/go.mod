module github.com/wopta/goworkspace/companydata

go 1.21

replace github.com/wopta/goworkspace/companydata => ./

require (
	cloud.google.com/go v0.115.1
	cloud.google.com/go/bigquery v1.62.0
	cloud.google.com/go/firestore v1.16.0
	github.com/GoogleCloudPlatform/functions-framework-go v1.9.0
	github.com/dustin/go-humanize v1.0.1
	github.com/go-gota/gota v0.12.0
	github.com/google/uuid v1.6.0
	github.com/oliveagle/jsonpath v0.0.0-20180606110733-2e52cf6e6852
	github.com/wopta/goworkspace/inclusive v0.0.4
	github.com/wopta/goworkspace/lib v1.0.128
	github.com/wopta/goworkspace/mail v1.0.120
	github.com/wopta/goworkspace/models v1.2.11
	github.com/wopta/goworkspace/network v1.0.73
	github.com/wopta/goworkspace/product v1.1.7
	github.com/wopta/goworkspace/user v1.0.105
	github.com/xuri/excelize/v2 v2.8.1
	google.golang.org/api v0.195.0
)

require (
	cloud.google.com/go/auth v0.9.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.4 // indirect
	cloud.google.com/go/compute/metadata v0.5.0 // indirect
	cloud.google.com/go/iam v1.1.13 // indirect
	cloud.google.com/go/longrunning v0.5.12 // indirect
	cloud.google.com/go/pubsub v1.42.0 // indirect
	cloud.google.com/go/storage v1.43.0 // indirect
	firebase.google.com/go v3.13.0+incompatible // indirect
	firebase.google.com/go/v4 v4.10.0 // indirect
	github.com/MicahParks/keyfunc v1.5.1 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20220527190237-ee62e23da966 // indirect
	github.com/apache/arrow/go/v15 v15.0.2 // indirect
	github.com/bmatcuk/doublestar v1.3.2 // indirect
	github.com/cloudevents/sdk-go/v2 v2.15.2 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-chi/chi/v5 v5.0.12 // indirect
	github.com/go-chi/cors v1.2.1 // indirect
	github.com/go-jose/go-jose/v4 v4.0.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.13.0 // indirect
	github.com/hyperjumptech/grule-rule-engine v1.12.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kevinburke/ssh_config v0.0.0-20190725054713-01f96b0aa0cd // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pkg/sftp v1.13.5 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.3 // indirect
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/src-d/gcfg v1.4.0 // indirect
	github.com/wopta/goworkspace/wiseproxy v1.0.3 // indirect
	github.com/xanzy/ssh-agent v0.2.1 // indirect
	github.com/xuri/efp v0.0.0-20231025114914-d1ff6096ae53 // indirect
	github.com/xuri/nfp v0.0.0-20230919160717-d98342af3f05 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/oauth2 v0.22.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/time v0.6.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	gonum.org/v1/gonum v0.12.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/appengine/v2 v2.0.2 // indirect
	google.golang.org/genproto v0.0.0-20240823204242-4ba0660f739c // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240823204242-4ba0660f739c // indirect
	google.golang.org/grpc v1.65.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.2 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)

module authentication_service

go 1.18

require (
	github.com/docker/go-connections v0.4.0
	github.com/gin-contrib/sessions v0.0.5
	github.com/gin-gonic/gin v1.8.1
	github.com/go-playground/validator/v10 v10.11.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/ilyakaznacheev/cleanenv v1.3.0
	github.com/jaswdr/faker v1.16.0
	github.com/joho/godotenv v1.4.0
	github.com/markbates/goth v1.75.2
	github.com/nats-io/nats.go v1.23.0
	github.com/pquerna/otp v1.4.0
	github.com/rs/zerolog v1.28.0
	github.com/sethvargo/go-password v0.2.0
	github.com/stretchr/testify v1.8.1
	github.com/testcontainers/testcontainers-go v0.17.0
	go.mongodb.org/mongo-driver v1.11.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.37.0
	go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo v0.37.0
	go.opentelemetry.io/otel v1.11.2
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.11.2
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.11.2
	go.opentelemetry.io/otel/sdk v1.11.2
	go.uber.org/dig v1.15.0
	golang.org/x/crypto v0.5.0
	google.golang.org/grpc v1.51.0
)

require (
	cloud.google.com/go v0.67.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc // indirect
	github.com/bos-hieu/mongostore v0.0.2 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/containerd/containerd v1.6.12 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.20+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/gorilla/securecookie v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.15 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.0.0-20221128092401-c43b287e0e0f // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nats-io/nats-server/v2 v2.9.14 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/opencontainers/runc v1.1.3 // indirect
	github.com/pelletier/go-toml/v2 v2.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.11.2 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220617124728-180714bec0ad // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/docker/docker => github.com/docker/docker v20.10.3-0.20221013203545-33ab36d6b304+incompatible // 22.06 branch

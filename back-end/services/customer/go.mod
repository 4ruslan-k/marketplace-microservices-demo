module customer

go 1.20

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.14.0
	github.com/ilyakaznacheev/cleanenv v1.3.0
	github.com/joho/godotenv v1.4.0
	github.com/nats-io/nats.go v1.28.0
	github.com/rs/zerolog v1.30.0
	github.com/uptrace/bun v1.1.14
	github.com/urfave/cli/v2 v2.25.1
	gopkg.in/yaml.v2 v2.4.0
	shared v0.0.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.4.1 // indirect
	github.com/nats-io/nkeys v0.4.4 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	github.com/uptrace/bun/dialect/pgdialect v1.1.14 // indirect
	github.com/uptrace/bun/driver/pgdriver v1.1.14 // indirect
	github.com/uptrace/bun/extra/bundebug v1.1.14 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mellium.im/sasl v0.3.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/docker/docker => github.com/docker/docker v20.10.3-0.20221013203545-33ab36d6b304+incompatible // 22.06 branch

replace shared v0.0.0 => ../../pkg

module service2

go 1.22.4

replace services-challenge/proto => ../proto

require (
	github.com/jinzhu/configor v1.2.2
	google.golang.org/grpc v1.65.0
	services-challenge/proto v0.0.0-00010101000000-000000000000
)

require (
	github.com/BurntSushi/toml v1.2.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

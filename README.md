# gframe
go framework

protoc-gen-go-gframe

安装：
  go install github.com/brunowang/gframe/cmd/protoc-gen-go-gframe@latest
  
使用：
  protoc -I $include_proto_dir --go-gframe_out=paths=source_relative,pbGoDir=$pbgen_go_dir:./ $proto_path

dao-gen

安装：
  go install github.com/brunowang/gframe/cmd/dao-gen@latest

使用：
  dao-gen -f create_table.sql

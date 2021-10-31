# protoc-gen-jsonnet
Generate Jsonnet definition for JSON representation of protobuf object

## Generated message docs

Assuming protobuf

```protobuf
enum EnumValues {
    OK = 0
}

message Test {
    string v = 1;
    EnumValues e = 2;
}
```

Construct this message in Jsonnet by

```jsonnet
local test = import "test.jsonnet";

test.Test.new(v="test", e=test.EnumValues.OK)
// or
test.Test.new()
    + test.Test.withV("test")
    + test.Test.withE(test.EnumValues.OK)
```

## TODO

- [ ] WKT

## License

[Licensed under 3 clause BSD license](LICENSE) with output exception, same as Protobuf itself.

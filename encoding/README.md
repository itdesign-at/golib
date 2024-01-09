package encoding

| Function                                       | Comment                            |
|------------------------------------------------|------------------------------------|
| func MarshalFile(file string, in any) error    | encode yaml, json or gob into file |
| func UnmarshalFile(file string, out any) error | decode yaml, json or gob from file |

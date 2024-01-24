package converter

| Function                                                    | Comment                                               |
|-------------------------------------------------------------|-------------------------------------------------------|
| func Denormalize(in string) string                          | Denormalize does the opposite of Normalize            |
| func Normalize(in string) string                            | Normalize a string by removing all special characters |
| func GetRandomString(length int) string                     | returns a random string with length given             |
| func GetMobileNumber(input string) (string, error)          | extract the mobile number without blanks              |
| func PrepareMobileNumber(input interface{}) (string, error) | PrepareMobileNumber returns 0043664....               |
| func Fnv1aHash(input string) string                         | returns a 32 bit FNV-1a hash                          |
| func B64Dec(b64s string) (string, error)                    | decode string in base64 format                        |
| func B64MustDec(b64s string) string                         | decode string in base64 format, panics on error       |
| func B64Enc(s string) string                                | encode string to base64 format                        |
| func Dec(format string, src string, dst interface{}) error  | decode string in specific format                      |
| func Enc(format string, src interface{}) (string, error)    | encode string to specific format                      |
| func MustDec(format string, src string, dst interface{})    | decode string to specific format, panics on error     |
| func MustEnc(format string, src interface{}) string         | encode string to specific format, panics on error     |
| func JsonDec(src string, dst interface{}) error             | decode json in specific format                        |
| func JsonEnc(src interface{}) (string, error)               | encode json to specific format                        |
| func JsonMustDec(src string, dst interface{})               | decode json to specific format, panics on error       |
| func JsonMustEnc(src interface{}) string                    | encode json to specific format, panics on error       |

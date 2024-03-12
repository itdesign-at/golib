# golib

first fetch all your tags and display all of them

    git fetch --tags
    git tag -l
    ... output with list of tags ...
    
you have to prefix the tag with the folder name, e.g.: commandLine/v1.0.0

    git tag -a commandLine/v1.0.1 -m "Release 1.0.1"
    git tag -a converter/v1.0.3 -m "Release 1.0.3"
    git tag -a crypt/v1.0.0 -m "Release 1.0.0"
    git tag -a encoding/v1.0.1 -m "Release 1.0.1"
    git tag -a keyvalue/v1.0.2 -m "Release 1.0.2"
    git tag -a macro/v1.0.0 -m "Release 1.0.0"
    git tag -a structuredLogging/v1.0.3 -m "Release 1.0.3"

after you create a new tag for a specific package you also have to create a new tag for the whole library

    git tag -a v1.0.2 -m "Release 1.0.2"
    
    git push --tags

See [go module documentation](https://go.dev/doc/modules/managing-source) for more information.

# package commandLine

| Function                                               | Comment                                                             |
|--------------------------------------------------------|---------------------------------------------------------------------|
| func CheckRootUserAndExit()                            | CheckRootUserAndExit terminates the program when you are not "root" |
| func Parse(args []string) keyvalue.Record              | Reads os.Args into the map returned                                 |
| func PrintVersion(args keyvalue.Record, detailed bool) | Prints version info                                                 |

# package converter

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

# package encoding

| Function                                       | Comment                            |
|------------------------------------------------|------------------------------------|
| func MarshalFile(file string, in any) error    | encode yaml, json or gob into file |
| func UnmarshalFile(file string, out any) error | decode yaml, json or gob from file |

# package keyvalue

| Function                                                     | Comment                         |
|--------------------------------------------------------------|---------------------------------|
| func NewRecord() Record                                      | does make for a new record      |
| func (r Record) Exists(key string) bool                      | checks the existence of the key |
| func (r Record) Value(key string) interface{}                | returns the value               |
| func (r Record) Bool(key string, convert ...bool) bool       | returns a bool value            |
| func (r Record) Float64(key string, convert ...bool) float64 | returns a float value           |
| func (r Record) Int(key string, convert ...bool) int         | returns an int value            |
| func (r Record) Int64(key string, convert ...bool) int64     | returns an int64 value          |
| func (r Record) String(key string, convert ...bool) string   | returns a string value          |
| func (r Record) Copy() Record                                | copy a record                   |
| func (r Record) GetSortedKeys() []string                     | returns keys sorted             |

# package structuredLogging

| Function                                                            | Comment                        |
|---------------------------------------------------------------------|--------------------------------|
| func New(dsn ...string) *SlogLogger                                 | New creates an slog writer     |
| func (sl *SlogLogger) Parameter(params keyvalue.Record) *SlogLogger | add params                     |
| func (sl *SlogLogger) InitJsonHandler() *slog.JSONHandler           | Get JSONHandler with writer(s) |
| func (sl *SlogLogger) Init() *SlogLogger                            | Init the logger with handler   |
| func (sl *SlogLogger) Write(p []byte) (int, error)                  | Write to all log writers       |
| func GenerateLogfileName(layout string) string                      | Get a logfilename              |                 

# package crypt

| Function                                                                  | Comment                                                     |
|---------------------------------------------------------------------------|-------------------------------------------------------------|
| func NewBcrypt() Bcrypt                                                   | creates an bcrypt handler                                   |
| func (b *bcryptProperties) Encrypt(plainText string) (string, error)      | encrypt plaintext                                           |
| func (b *bcryptProperties) HashedText(hashedText string)                  | set a hashed text string                                    |
| func (b *bcryptProperties) PlainText(plainText string)                    | set a plain text string                                     |
| func (b *bcryptProperties) Cost(cost int)                                 | sets bcrypt costs                                           |
| func (b *bcryptProperties) Compare() bool                                 | checks the hash against the plaintext                       |                 
| func GenerateEd25519KeyFiles(dir string, filename string) (string, error) | generates ed25519 key files (id_ed25519 and id_ed25519.pub) |                 
| func NewEncryptedString(plainTextValue string) EncryptedString            | creates an Encrypted string                                 |                 
| func NewDecryptedString(encryptedValue string) string                     | decrypt an encrypted string                                 |                 
| func (v EncryptedString) Value() string                                   | returns the decrypted plainTextValue                        |                 
| func (v EncryptedString) String() string                                  | returns the decrypted plainTextValue                        |                 
| func (v EncryptedString) MarshalText() ([]byte, error)                    |                                                             |                 
| func (v *EncryptedString) UnmarshalText(text []byte) error                |                                                             |                 
| func (v EncryptedString) MarshalBinary() ([]byte, error)                  |                                                             |                 
| func NewSymmetricEncryption() *SymCrypt                                   | creates an SymCrypt handler                                 |                 
| func (s *SymCrypt) SetKey(key string) *SymCrypt                           | set an AES key                                              |                 
| func (s *SymCrypt) SetPlainText(plainText string) *SymCrypt               | set plain text to encrypt.                                  |                 
| func (s *SymCrypt) GetCypherBase64() string                               | returns the encrypted data stream as base64 encoded         |                 
| func (s *SymCrypt) SetCypherBase64(base64String string) *SymCrypt         | set cipher text as base64 string.                           |                 
| func (s *SymCrypt) GetPlainText() (string, error)                         | returns the plaintext                                       |                 

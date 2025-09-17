# 🐧 HRMS Penguin

## TO DO:

- [x] Check attendance record
  - [x] Basic
  - [x] With leave application records
- [x] Save Password safely in config
- [x] Vacation Application
- [ ] Clean up old struct with `UnmarshalJSON`
- [x] Remove logging pw :(

## Maybe DO:
- [ ] Proper docs
- [ ] GitHub actions to build binaries
  - [ ] Tag versioning
  - [ ] Release notes
- [ ] TUI with Charm

## Go
- Unmarshal JSON
  - with tags
  - custom `UnmarshalJSON`
- time package
  - parse
  - format
  - compare
- net/http
  - cookies jar
  - make requests
- cobra/cli
  - sub commands
  - flags
- test
  - t.Errorf
  - reflect.DeepEqual
- slices
  - DeleteFunc
  - SortFunc
- Go's package
  - import internal package
- Debugger
  - dlv
  - launch.json
# lazydb

## TODO

### Basic Functionality

- [ ] add new row
- [ ] custom SQL editor on new page (by pressing 3)
- [ ] support for other DB drivers (refactor to use interfaces)
- [ ] show error modals for failed queries and other user-facing errors
- [ ] table results pagination
- [x] update timestamp fields to NOW()
- [x] row value viewer (broken atm)
- [x] view for columns and indexes
- [x] home page - read config file for list of databases
- [x] home page - connect to a database
- [x] select/switch databases
- [x] delete row

### QOL Improvements

- [ ] tabs similar to Sequel Ace (keep sessions for each table opened)
- [ ] allow refreshing of tables with hidden columns
- [ ] dynamic hiding/showing of columns in a new modal - similar to Lazygit staging/unstaging files where pressing space toggles, and pressing A toggles all
- [ ] improve UI colors - similar to lazygit
- [ ] query history
- [ ] keyboard shorcuts
  - [ ] ctrl+hjkl to move panels (in addtn to tab)
  - [ ] 0 and $ goes to start/end of row
  - [ ] ctrl+p command palette type to go to table
  - [x] ctrl+f from anywhere goes to table filter
  - [x] press keybind on a cell (W), automatically write a WHERE condition for column
- [ ] queries history (press ctrl+n or ctrl+p on WHERE filter scrolls through history)
- [ ] saved queries
- [ ] help menu by pressing ?
- [ ] improve SQL editor autocomplete for DB columns
- [x] easy way to show/hide columns
- [x] yank cell value
- [x] sort columns by highlighting header name and pressing a keybind

# plan

- 3 layer design

p.1
- modules
- http handler 
- unit testing,unit httptest
- - assert, require 
- - go test [-count=1, -cover]
- - vscode run test, show cover
- basic main
- - zerolog
- chi router + middleware
- context 


p.2
- database/sql
- database integration test
- - debugging
- gracefull shutdown
- create config 
- add logger
- - info levels
- - init in main
- - inject into handlers
- http server setup
- http client setup
- server to server 
- domain error handler
- - custom errors tricks with interface{}

p.3
- compose API 
- sync pkg
- errgroup

p.4
- linter
- makefile
- - build, test, tools

p.X
- grpc


Requested:

- integration test for handler 
- show channels



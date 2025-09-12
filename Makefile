LEXERTEST=./test/lexer_test/lexer_test.go
##ASTTEST=./test/ast_test/ast_test.go
PARSERTEST=./test/parser_test/parser_test.go
MAINFILE=./src/main.go
OUTPUTFILE=./bin/main 
GXX=go
test:
	@echo "start lexer test"
	@ $(GXX) test $(LEXERTEST)
	@echo "finish lexer test"
	@echo "start parser test"
	@ $(GXX) test $(PARSERTEST)
	@echo "finish parser test"
build:
	@echo "start to build ......"
	@ $(GXX) build -o $(OUTPUTFILE) $(MAINFILE) 
	@echo "finish the building ........"
run: build
	@echo "start to run the repl ........"
	@./bin/main 

.PHONY: test

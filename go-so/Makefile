
clean:
	rm *.h *.so

shared-so:
	/opt/dev/go/bin/go build -o a.so -buildmode=c-shared main.go

run-test: shared-so
	python3 test.py

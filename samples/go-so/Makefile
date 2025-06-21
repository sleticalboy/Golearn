
clean:
	rm *.h *.so

shared-so:
	go build -o libfib.so -buildmode=c-shared main.go

run-pytest: shared-so
	python3 test.py

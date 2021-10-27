.PHONY: all run
all: run

run:
	ls *.go | grep -v '_test.go' | xargs go run

enque:
	aws sqs send-message \
		--endpoint-url http://localhost:9324 \
		--queue-url http://localhost:9324/queue/osint-website \
		--message-body '{"osint_id":1002, "osint_data_source_id":1003,"rel_osint_data_source_id":1002, "data_source":"osint:website", "project_id":1001, "resource_name":"https://console.security-hub.jp", "resource_type":"website", "detect_word":""}'

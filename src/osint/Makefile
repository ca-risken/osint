.PHONY: all
all: help

.PHONY: help
help:
	@echo "Usage: make <sub-command>"
	@echo "\n---------------- sub-command list ----------------"
	@cat Makefile | grep -e "^.PHONY:" | grep -v "all" | cut -f2 -d' '

.PHONY: list-service
list-service:
	grpcurl -plaintext localhost:18081 list osint.osint.OsintService

.PHONY: list-osint
list-osint:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		localhost:18081 osint.osint.OsintService.ListOsint

.PHONY: get-osint
get-osint:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "osint_id":1001}' \
		localhost:18081 osint.osint.OsintService.GetOsint

.PHONY: put-osint
put-osint:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001 ,"osint":{"resource_type":"Domain","resource_name":"cyberagent.co.jp","project_id":1001}}' \
		localhost:18081 osint.osint.OsintService.PutOsint

.PHONY: delete-osint
delete-osint:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "osint_id":1002}' \
		localhost:18081 osint.osint.OsintService.DeleteOsint

.PHONY: list-osint_datasource
list-osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		localhost:18081 osint.osint.OsintService.ListOsintDataSource

.PHONY: get-osint_datasource
get-osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "osint_data_source_id":1001}' \
		localhost:18081 osint.osint.OsintService.GetOsintDataSource

.PHONY: put-osint_datasource
put-osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001 ,"osint_data_source":{"name":"test_ds","description":"for_test","max_score":10}}' \
		localhost:18081 osint.osint.OsintService.PutOsintDataSource

.PHONY: delete-osint_datasource
delete-osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "osint_data_source_id":1002}' \
		localhost:18081 osint.osint.OsintService.DeleteOsintDataSource

.PHONY: list-rel_osint_datasource
list-rel_osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		localhost:18081 osint.osint.OsintService.ListRelOsintDataSource

.PHONY: get-rel_osint_datasource
get-rel_osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "rel_osint_data_source_id":1001}' \
		localhost:18081 osint.osint.OsintService.GetRelOsintDataSource

.PHONY: put-rel_osint_datasource
put-rel_osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "rel_osint_data_source": {"osint_id":1001, "osint_data_source_id":1001, "project_id":1001, "status":"CONFIGURED"}}' \
		localhost:18081 osint.osint.OsintService.PutRelOsintDataSource

.PHONY: delete-rel_osint_datasource
delete-rel_osint_datasource:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "rel_osint_data_source_id":1001}' \
		localhost:18081 osint.osint.OsintService.DeleteRelOsintDataSource

.PHONY: list-rel_osint_detect_word
list-rel_osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		localhost:18081 osint.osint.OsintService.ListRelOsintDetectWord

.PHONY: get-rel_osint_detect_word
get-rel_osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "rel_osint_detect_word_id":1001}' \
		localhost:18081 osint.osint.OsintService.GetRelOsintDetectWord

.PHONY: put-rel_osint_detect_word
put-rel_osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "rel_osint_detect_word": {"osint_detect_word_id":1001, "rel_osint_data_source_id":1001, "project_id":1001}}' \
		localhost:18081 osint.osint.OsintService.PutRelOsintDetectWord

.PHONY: delete-rel_osint_detect_word
delete-rel_osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "rel_osint_detect_word_id":1001}' \
		localhost:18081 osint.osint.OsintService.DeleteRelOsintDetectWord

.PHONY: list-osint_detect_word
list-osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001}' \
		localhost:18081 osint.osint.OsintService.ListOsintDetectWord

.PHONY: get-osint_detect_word
get-osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "osint_detect_word_id":1001}' \
		localhost:18081 osint.osint.OsintService.GetOsintDetectWord

.PHONY: put-osint_detect_word
put-osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "osint_detect_word": {"word":"hoge", "project_id":1001}}' \
		localhost:18081 osint.osint.OsintService.PutOsintDetectWord

.PHONY: delete-osint_detect_word
delete-osint_detect_word:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "osint_detect_word_id":1001}' \
		localhost:18081 osint.osint.OsintService.DeleteOsintDetectWord

.PHONY: invoke-scan
invoke-scan:
	grpcurl \
		-plaintext \
		-d '{"project_id":1001, "rel_osint_data_source_id":1001}' \
		localhost:18081 osint.osint.OsintService.InvokeScan

.PHONY: invoke-scan_all
invoke-scan_all:
	grpcurl \
		-plaintext \
		-d '{}' \
		localhost:18081 osint.osint.OsintService.InvokeScanAll

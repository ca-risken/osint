# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [harvester/services.proto](#harvester/services.proto)
    - [InvokeScanRequest](#osint.harvester.InvokeScanRequest)
  
    - [ResourceService](#osint.harvester.ResourceService)
  
- [osint/entities.proto](#osint/entities.proto)
    - [Osint](#osint.osint.Osint)
    - [OsintDataSource](#osint.osint.OsintDataSource)
    - [OsintDataSourceForUpsert](#osint.osint.OsintDataSourceForUpsert)
    - [OsintDetectWord](#osint.osint.OsintDetectWord)
    - [OsintDetectWordForUpsert](#osint.osint.OsintDetectWordForUpsert)
    - [OsintForUpsert](#osint.osint.OsintForUpsert)
    - [RelOsintDataSource](#osint.osint.RelOsintDataSource)
    - [RelOsintDataSourceForUpsert](#osint.osint.RelOsintDataSourceForUpsert)
  
    - [Status](#osint.osint.Status)
  
- [osint/services.proto](#osint/services.proto)
    - [DeleteOsintDataSourceRequest](#osint.osint.DeleteOsintDataSourceRequest)
    - [DeleteOsintDetectWordRequest](#osint.osint.DeleteOsintDetectWordRequest)
    - [DeleteOsintRequest](#osint.osint.DeleteOsintRequest)
    - [DeleteRelOsintDataSourceRequest](#osint.osint.DeleteRelOsintDataSourceRequest)
    - [GetOsintDataSourceRequest](#osint.osint.GetOsintDataSourceRequest)
    - [GetOsintDataSourceResponse](#osint.osint.GetOsintDataSourceResponse)
    - [GetOsintDetectWordRequest](#osint.osint.GetOsintDetectWordRequest)
    - [GetOsintDetectWordResponse](#osint.osint.GetOsintDetectWordResponse)
    - [GetOsintRequest](#osint.osint.GetOsintRequest)
    - [GetOsintResponse](#osint.osint.GetOsintResponse)
    - [GetRelOsintDataSourceRequest](#osint.osint.GetRelOsintDataSourceRequest)
    - [GetRelOsintDataSourceResponse](#osint.osint.GetRelOsintDataSourceResponse)
    - [InvokeScanRequest](#osint.osint.InvokeScanRequest)
    - [InvokeScanResponse](#osint.osint.InvokeScanResponse)
    - [ListOsintDataSourceRequest](#osint.osint.ListOsintDataSourceRequest)
    - [ListOsintDataSourceResponse](#osint.osint.ListOsintDataSourceResponse)
    - [ListOsintDetectWordRequest](#osint.osint.ListOsintDetectWordRequest)
    - [ListOsintDetectWordResponse](#osint.osint.ListOsintDetectWordResponse)
    - [ListOsintRequest](#osint.osint.ListOsintRequest)
    - [ListOsintResponse](#osint.osint.ListOsintResponse)
    - [ListRelOsintDataSourceRequest](#osint.osint.ListRelOsintDataSourceRequest)
    - [ListRelOsintDataSourceResponse](#osint.osint.ListRelOsintDataSourceResponse)
    - [PutOsintDataSourceRequest](#osint.osint.PutOsintDataSourceRequest)
    - [PutOsintDataSourceResponse](#osint.osint.PutOsintDataSourceResponse)
    - [PutOsintDetectWordRequest](#osint.osint.PutOsintDetectWordRequest)
    - [PutOsintDetectWordResponse](#osint.osint.PutOsintDetectWordResponse)
    - [PutOsintRequest](#osint.osint.PutOsintRequest)
    - [PutOsintResponse](#osint.osint.PutOsintResponse)
    - [PutRelOsintDataSourceRequest](#osint.osint.PutRelOsintDataSourceRequest)
    - [PutRelOsintDataSourceResponse](#osint.osint.PutRelOsintDataSourceResponse)
  
    - [OsintService](#osint.osint.OsintService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="harvester/services.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## harvester/services.proto



<a name="osint.harvester.InvokeScanRequest"></a>

### InvokeScanRequest
Invoke Scan


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| resource_name | [string](#string) |  |  |
| resource_type | [string](#string) |  |  |





 

 

 


<a name="osint.harvester.ResourceService"></a>

### ResourceService
Resource
rpc ListOsint(ListOsintRequest) returns (ListOsintResponse) {}

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| InvokeScan | [InvokeScanRequest](#osint.harvester.InvokeScanRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | InvokeScan |

 



<a name="osint/entities.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osint/entities.proto



<a name="osint.osint.Osint"></a>

### Osint



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| resource_type | [string](#string) |  |  |
| resource_name | [string](#string) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="osint.osint.OsintDataSource"></a>

### OsintDataSource



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_data_source_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| max_score | [float](#float) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="osint.osint.OsintDataSourceForUpsert"></a>

### OsintDataSourceForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_data_source_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| max_score | [float](#float) |  |  |






<a name="osint.osint.OsintDetectWord"></a>

### OsintDetectWord



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_detect_word_id | [uint32](#uint32) |  |  |
| rel_osint_data_source_id | [uint32](#uint32) |  |  |
| word | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="osint.osint.OsintDetectWordForUpsert"></a>

### OsintDetectWordForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_detect_word_id | [uint32](#uint32) |  |  |
| rel_osint_data_source_id | [uint32](#uint32) |  |  |
| word | [string](#string) |  |  |
| project_id | [uint32](#uint32) |  |  |






<a name="osint.osint.OsintForUpsert"></a>

### OsintForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| resource_type | [string](#string) |  |  |
| resource_name | [string](#string) |  |  |






<a name="osint.osint.RelOsintDataSource"></a>

### RelOsintDataSource



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_osint_data_source_id | [uint32](#uint32) |  |  |
| osint_data_source_id | [uint32](#uint32) |  |  |
| osint_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| status | [Status](#osint.osint.Status) |  |  |
| status_detail | [string](#string) |  |  |
| scan_at | [int64](#int64) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="osint.osint.RelOsintDataSourceForUpsert"></a>

### RelOsintDataSourceForUpsert



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_osint_data_source_id | [uint32](#uint32) |  |  |
| osint_data_source_id | [uint32](#uint32) |  |  |
| osint_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |
| status | [Status](#osint.osint.Status) |  |  |
| status_detail | [string](#string) |  |  |
| scan_at | [int64](#int64) |  |  |





 


<a name="osint.osint.Status"></a>

### Status
Status

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| OK | 1 |  |
| CONFIGURED | 2 |  |
| NOT_CONFIGURED | 3 |  |
| ERROR | 4 |  |


 

 

 



<a name="osint/services.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## osint/services.proto



<a name="osint.osint.DeleteOsintDataSourceRequest"></a>

### DeleteOsintDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| osint_data_source_id | [uint32](#uint32) |  |  |






<a name="osint.osint.DeleteOsintDetectWordRequest"></a>

### DeleteOsintDetectWordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| osint_detect_word_id | [uint32](#uint32) |  |  |






<a name="osint.osint.DeleteOsintRequest"></a>

### DeleteOsintRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| osint_id | [uint32](#uint32) |  |  |






<a name="osint.osint.DeleteRelOsintDataSourceRequest"></a>

### DeleteRelOsintDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| rel_osint_data_source_id | [uint32](#uint32) |  |  |






<a name="osint.osint.GetOsintDataSourceRequest"></a>

### GetOsintDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| osint_data_source_id | [uint32](#uint32) |  |  |






<a name="osint.osint.GetOsintDataSourceResponse"></a>

### GetOsintDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_data_source | [OsintDataSource](#osint.osint.OsintDataSource) |  |  |






<a name="osint.osint.GetOsintDetectWordRequest"></a>

### GetOsintDetectWordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| osint_detect_word_id | [uint32](#uint32) |  |  |






<a name="osint.osint.GetOsintDetectWordResponse"></a>

### GetOsintDetectWordResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_detect_word | [OsintDetectWord](#osint.osint.OsintDetectWord) |  |  |






<a name="osint.osint.GetOsintRequest"></a>

### GetOsintRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| osint_id | [uint32](#uint32) |  |  |






<a name="osint.osint.GetOsintResponse"></a>

### GetOsintResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint | [Osint](#osint.osint.Osint) |  |  |






<a name="osint.osint.GetRelOsintDataSourceRequest"></a>

### GetRelOsintDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_osint_data_source_id | [uint32](#uint32) |  |  |
| project_id | [uint32](#uint32) |  |  |






<a name="osint.osint.GetRelOsintDataSourceResponse"></a>

### GetRelOsintDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_osint_data_source | [RelOsintDataSource](#osint.osint.RelOsintDataSource) |  |  |






<a name="osint.osint.InvokeScanRequest"></a>

### InvokeScanRequest
Invoke Scan


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| rel_osint_data_source_id | [uint32](#uint32) |  |  |






<a name="osint.osint.InvokeScanResponse"></a>

### InvokeScanResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [string](#string) |  |  |






<a name="osint.osint.ListOsintDataSourceRequest"></a>

### ListOsintDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| name | [string](#string) |  |  |






<a name="osint.osint.ListOsintDataSourceResponse"></a>

### ListOsintDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_data_source | [OsintDataSource](#osint.osint.OsintDataSource) | repeated |  |






<a name="osint.osint.ListOsintDetectWordRequest"></a>

### ListOsintDetectWordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| rel_osint_data_source_id | [uint32](#uint32) |  |  |






<a name="osint.osint.ListOsintDetectWordResponse"></a>

### ListOsintDetectWordResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_detect_word | [OsintDetectWord](#osint.osint.OsintDetectWord) | repeated |  |






<a name="osint.osint.ListOsintRequest"></a>

### ListOsintRequest
Osint Service


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |






<a name="osint.osint.ListOsintResponse"></a>

### ListOsintResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint | [Osint](#osint.osint.Osint) | repeated |  |






<a name="osint.osint.ListRelOsintDataSourceRequest"></a>

### ListRelOsintDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| osint_id | [uint32](#uint32) |  |  |
| osint_data_source_id | [uint32](#uint32) |  |  |






<a name="osint.osint.ListRelOsintDataSourceResponse"></a>

### ListRelOsintDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_osint_data_source | [RelOsintDataSource](#osint.osint.RelOsintDataSource) | repeated |  |






<a name="osint.osint.PutOsintDataSourceRequest"></a>

### PutOsintDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| osint_data_source | [OsintDataSourceForUpsert](#osint.osint.OsintDataSourceForUpsert) |  |  |






<a name="osint.osint.PutOsintDataSourceResponse"></a>

### PutOsintDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_data_source | [OsintDataSource](#osint.osint.OsintDataSource) |  |  |






<a name="osint.osint.PutOsintDetectWordRequest"></a>

### PutOsintDetectWordRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| osint_detect_word | [OsintDetectWordForUpsert](#osint.osint.OsintDetectWordForUpsert) |  |  |






<a name="osint.osint.PutOsintDetectWordResponse"></a>

### PutOsintDetectWordResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint_detect_word | [OsintDetectWord](#osint.osint.OsintDetectWord) |  |  |






<a name="osint.osint.PutOsintRequest"></a>

### PutOsintRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| osint | [OsintForUpsert](#osint.osint.OsintForUpsert) |  |  |






<a name="osint.osint.PutOsintResponse"></a>

### PutOsintResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| osint | [Osint](#osint.osint.Osint) |  |  |






<a name="osint.osint.PutRelOsintDataSourceRequest"></a>

### PutRelOsintDataSourceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| project_id | [uint32](#uint32) |  |  |
| rel_osint_data_source | [RelOsintDataSourceForUpsert](#osint.osint.RelOsintDataSourceForUpsert) |  |  |






<a name="osint.osint.PutRelOsintDataSourceResponse"></a>

### PutRelOsintDataSourceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rel_osint_data_source | [RelOsintDataSource](#osint.osint.RelOsintDataSource) |  |  |





 

 

 


<a name="osint.osint.OsintService"></a>

### OsintService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListOsint | [ListOsintRequest](#osint.osint.ListOsintRequest) | [ListOsintResponse](#osint.osint.ListOsintResponse) | Osint |
| GetOsint | [GetOsintRequest](#osint.osint.GetOsintRequest) | [GetOsintResponse](#osint.osint.GetOsintResponse) |  |
| PutOsint | [PutOsintRequest](#osint.osint.PutOsintRequest) | [PutOsintResponse](#osint.osint.PutOsintResponse) |  |
| DeleteOsint | [DeleteOsintRequest](#osint.osint.DeleteOsintRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListRelOsintDataSource | [ListRelOsintDataSourceRequest](#osint.osint.ListRelOsintDataSourceRequest) | [ListRelOsintDataSourceResponse](#osint.osint.ListRelOsintDataSourceResponse) | RelOsintDataSource |
| GetRelOsintDataSource | [GetRelOsintDataSourceRequest](#osint.osint.GetRelOsintDataSourceRequest) | [GetRelOsintDataSourceResponse](#osint.osint.GetRelOsintDataSourceResponse) |  |
| PutRelOsintDataSource | [PutRelOsintDataSourceRequest](#osint.osint.PutRelOsintDataSourceRequest) | [PutRelOsintDataSourceResponse](#osint.osint.PutRelOsintDataSourceResponse) |  |
| DeleteRelOsintDataSource | [DeleteRelOsintDataSourceRequest](#osint.osint.DeleteRelOsintDataSourceRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListOsintDataSource | [ListOsintDataSourceRequest](#osint.osint.ListOsintDataSourceRequest) | [ListOsintDataSourceResponse](#osint.osint.ListOsintDataSourceResponse) | OsintDataSource |
| GetOsintDataSource | [GetOsintDataSourceRequest](#osint.osint.GetOsintDataSourceRequest) | [GetOsintDataSourceResponse](#osint.osint.GetOsintDataSourceResponse) |  |
| PutOsintDataSource | [PutOsintDataSourceRequest](#osint.osint.PutOsintDataSourceRequest) | [PutOsintDataSourceResponse](#osint.osint.PutOsintDataSourceResponse) |  |
| DeleteOsintDataSource | [DeleteOsintDataSourceRequest](#osint.osint.DeleteOsintDataSourceRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| ListOsintDetectWord | [ListOsintDetectWordRequest](#osint.osint.ListOsintDetectWordRequest) | [ListOsintDetectWordResponse](#osint.osint.ListOsintDetectWordResponse) | OsintDetectWord |
| GetOsintDetectWord | [GetOsintDetectWordRequest](#osint.osint.GetOsintDetectWordRequest) | [GetOsintDetectWordResponse](#osint.osint.GetOsintDetectWordResponse) |  |
| PutOsintDetectWord | [PutOsintDetectWordRequest](#osint.osint.PutOsintDetectWordRequest) | [PutOsintDetectWordResponse](#osint.osint.PutOsintDetectWordResponse) |  |
| DeleteOsintDetectWord | [DeleteOsintDetectWordRequest](#osint.osint.DeleteOsintDetectWordRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |
| InvokeScan | [InvokeScanRequest](#osint.osint.InvokeScanRequest) | [InvokeScanResponse](#osint.osint.InvokeScanResponse) | Invoke |
| InvokeScanAll | [.google.protobuf.Empty](#google.protobuf.Empty) | [.google.protobuf.Empty](#google.protobuf.Empty) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |


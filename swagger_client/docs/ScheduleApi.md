# swagger_client.ScheduleApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**admin_schedule_get**](ScheduleApi.md#admin_schedule_get) | **GET** /admin/schedule | 
[**admin_schedule_id_delete**](ScheduleApi.md#admin_schedule_id_delete) | **DELETE** /admin/schedule/{id} | 
[**admin_schedule_id_get**](ScheduleApi.md#admin_schedule_id_get) | **GET** /admin/schedule/{id} | 
[**admin_schedule_post**](ScheduleApi.md#admin_schedule_post) | **POST** /admin/schedule | 

# **admin_schedule_get**
> Segments admin_schedule_get(start=start, stop=stop)



Get schedule cut

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.ScheduleApi(swagger_client.ApiClient(configuration))
start = 56 # int | start cut (optional)
stop = 56 # int | stop cut (optional)

try:
    api_response = api_instance.admin_schedule_get(start=start, stop=stop)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling ScheduleApi->admin_schedule_get: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **start** | **int**| start cut | [optional] 
 **stop** | **int**| stop cut | [optional] 

### Return type

[**Segments**](Segments.md)

### Authorization

[editorAuth](../README.md#editorAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **admin_schedule_id_delete**
> admin_schedule_id_delete(id)



Delete segment by id

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.ScheduleApi(swagger_client.ApiClient(configuration))
id = 789 # int | ID

try:
    api_instance.admin_schedule_id_delete(id)
except ApiException as e:
    print("Exception when calling ScheduleApi->admin_schedule_id_delete: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| ID | 

### Return type

void (empty response body)

### Authorization

[editorAuth](../README.md#editorAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **admin_schedule_id_get**
> Segment admin_schedule_id_get(id)



Get segment information by id

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.ScheduleApi(swagger_client.ApiClient(configuration))
id = 789 # int | ID

try:
    api_response = api_instance.admin_schedule_id_get(id)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling ScheduleApi->admin_schedule_id_get: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| ID | 

### Return type

[**Segment**](Segment.md)

### Authorization

[editorAuth](../README.md#editorAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **admin_schedule_post**
> InlineResponse2001 admin_schedule_post(body=body)



Create new segment

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.ScheduleApi(swagger_client.ApiClient(configuration))
body = swagger_client.SegmentRegister() # SegmentRegister | ... (optional)

try:
    api_response = api_instance.admin_schedule_post(body=body)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling ScheduleApi->admin_schedule_post: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**SegmentRegister**](SegmentRegister.md)| ... | [optional] 

### Return type

[**InlineResponse2001**](InlineResponse2001.md)

### Authorization

[editorAuth](../README.md#editorAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


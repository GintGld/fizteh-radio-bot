# swagger_client.RadioApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**id_file_get**](RadioApi.md#id_file_get) | **GET** /{id}/{file} | 
[**man_mpd_get**](RadioApi.md#man_mpd_get) | **GET** /man.mpd | 

# **id_file_get**
> str id_file_get(id, file)



Load source files for radio. Player automatically loads them.

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = swagger_client.RadioApi()
id = 789 # int | ID
file = 'file_example' # str | .m4s file

try:
    api_response = api_instance.id_file_get(id, file)
    pprint(api_response)
except ApiException as e:
    print("Exception when calling RadioApi->id_file_get: %s\n" % e)
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| ID | 
 **file** | **str**| .m4s file | 

### Return type

**str**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/octet-stream

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **man_mpd_get**
> str man_mpd_get()



Dash manifest for streaming. Entrypoint to start listening DASH streaming.

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint

# create an instance of the API class
api_instance = swagger_client.RadioApi()

try:
    api_response = api_instance.man_mpd_get()
    pprint(api_response)
except ApiException as e:
    print("Exception when calling RadioApi->man_mpd_get: %s\n" % e)
```

### Parameters
This endpoint does not need any parameter.

### Return type

**str**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/xml

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


# swagger_client.RootRadioApi

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**radio_start_get**](RootRadioApi.md#radio_start_get) | **GET** /radio/start | 
[**radio_stop_get**](RootRadioApi.md#radio_stop_get) | **GET** /radio/stop | 

# **radio_start_get**
> radio_start_get()



Start radio

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.RootRadioApi(swagger_client.ApiClient(configuration))

try:
    api_instance.radio_start_get()
except ApiException as e:
    print("Exception when calling RootRadioApi->radio_start_get: %s\n" % e)
```

### Parameters
This endpoint does not need any parameter.

### Return type

void (empty response body)

### Authorization

[rootAuth](../README.md#rootAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **radio_stop_get**
> radio_stop_get()



Stop radio

### Example
```python
from __future__ import print_function
import time
import swagger_client
from swagger_client.rest import ApiException
from pprint import pprint


# create an instance of the API class
api_instance = swagger_client.RootRadioApi(swagger_client.ApiClient(configuration))

try:
    api_instance.radio_stop_get()
except ApiException as e:
    print("Exception when calling RootRadioApi->radio_stop_get: %s\n" % e)
```

### Parameters
This endpoint does not need any parameter.

### Return type

void (empty response body)

### Authorization

[rootAuth](../README.md#rootAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)


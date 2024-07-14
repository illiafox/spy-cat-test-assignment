---
name: Bug report
about: Create a report to help us improve
title: "[BUG]"
labels: bug
assignees: ''

---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Application version: `v0.1.0`
2. (Optional) Postman collection version: `v1.0.0`
3. `curl` version of request Request. Example:
	```bash
	curl --location 'http://127.0.0.1:8080/cats/' \
	--header 'Content-Type: application/json' \
	--data '{
	    "name": "Alex",
	    "breed": "Abyssinian",
	    "experience": 10,
	    "salary": 5000
	}'
	```

**Expected behavior**
Response example
```
HTTP/1.1 201 Created
Content-Type: application/json
```
```json
{
    "ok": true,
    "id": 9
}
```

**Actual behavior**
Bad response example
```
HTTP/1.1 500 Internal Server Error
Content-Type: application/json
Content-Length: 97
```
```json
{
	"ok":  false,
	"code":  "INTERNAL_ERROR",
	"message":  "internal error, please try again later"
}
```

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Additional context**
Add any other context about the problem here.

Feature: API Testing with Godog
    In order to be a great programmer
    As a software engineer
    I need to be able to use Godog

  Scenario: Get User by ID
    Given the API is running
    When I receive "<request_method>" request to "<request_path>" with "<request_body>" body
    Then the response http status code should be "<response_http_status_code>"
    And the response body should be "<response_body>"

    Examples:
      | request_method | request_path | request_body | response_http_status_code | response_body              |
      | GET            | /users/1     | {}           |                       200 | {"id":1,"name":"Jhon Doe"} |

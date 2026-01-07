# API DOC

## Users

#### /users

- Method: GET
- Description: Retrieve a list of all users.
- Response:
  - 200 OK: Returns a list of users.

#### /create_user

- Method: POST
- Description: Create a new user.
- Request Body:
  - email: string (required)
  - password: string (required)
- Response:
  - 201 Created: User created successfully.
  - 400 Bad Request: Invalid input data.
  - 500 Internal Server Error: Server error.

#### /verify_user

- Method: GET
- Description: Verify user and create account.
- Request Body:
  - verification_code: string (required)
- Response:
  - 200 OK: User verified successfully.
  - 400 Bad Request: Invalid verification code.
  - 500 Internal Server Error: Server error.

# API DOC

## Server

#### /status

- Method: GET
- Description: Check the status of the server.
- Response:
  - 200 OK: Server is running.

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

## Events

#### /get_all_current_events

- Method: POST
- Description: Retrieve all current events for a user.
- Request Body:
  - user_id: string (required)
- Response:
  - 200 OK: Returns a list of current events.
  - 400 Bad Request: Invalid user ID.
  - 500 Internal Server Error: Server error.

#### /create_event

- Method: POST
- Description: Create a new event.
- Request Body:
  - start: timestamp (required)
  - end: timestamp (required)
  - event_name: string (required)
- Response:
  - 200 OK: Event created successfully.
  - 400 Bad Request: Invalid input data.
  - 500 Internal Server Error: Server error.

#### /delete_event

- Method: POST
- Description: Delete an existing event.
- Request Body:
  - event_id: string (required)
- Response:
  - 200 OK: Event deleted successfully.
  - 400 Bad Request: Invalid event ID.
  - 500 Internal Server Error: Server error.

#### /change_availability

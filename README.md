# fliqt-assignment

### API Documentation

* `POST -> /auth/signup`
    * Description: Use to create a new user
    * Request body: 
        ```
        {
            "username": "testuser", 
            "password": "testpassword"
        }
        ```
    * Response:
        - 201 Created: 
            ```
            {
                "message": "user created successfully"
            }
            ```
        - 400 Bad Request

* `POST -> /auth/login`
    * Description: Use to login and get JWT
    * Request body: 
        ```
        {
            "username": "testuser", 
            "password": "testpassword"
        }
        ```
    * Response:
        - 201 Created: 
            ```
            {
                "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzkyODQzOTMsInN1YiI6MX0.isXxf2Vshcur6VcSLAAikCKFimcdpXf3c1XLA7VnZ9U"
            }
            ```
        - 400 Bad Request

* `GET -> /user/<id>`
    * Description: Use to get user profile
    * Authentication inside Header: 
        ```
        Authorization: Bearer <JWT Token>
        ```
    * Request Params:
	    * id (required): user id
    
    * Response:
        - 200 OK: 
        ```
        {
            "id": 1,
            "username": "testuser",
            "created_at": "2025-02-11T00:44:53.306+08:00",
            "is_active": true,
            "roles": [
                "admin",
                "staff"
            ]
        }
        ```
        - 400 Bad Request
        - 401 Unauthorized

* `PUT -> /user/<id>/roles`
    * Description: Use to update specific user's role
    * Authentication inside Header:
        * only `admin` user can do this operation
        ```
        Authorization: Bearer <JWT Token>
        ```
    * Request Params:
	    * id (required): user id
    * Request Body:
        * only accept `staff` and `admin` role 
	    ```
        {
            "roles": ["staff", "admin"]
        }
        ```
    
    * Response:
        - 200 OK: 
        ```
        {
            "message": "user roles updated successfully"
        }
        ```
        - 400 Bad Request
        - 401 Unauthorized
        - 500 Internal Server Error


* `POST -> /clock/in`
    * Description: Use to clock in
    * Authentication inside Header: 
        ```
        Authorization: Bearer <JWT Token>
        ```
    * Response:
        - 201 Created: 
        ```
        {
            "message": "user clocked in for today successfully"
        }
        ```
        - 400 Bad Request
        - 401 Unauthorized
        - 500 Internal Server Error

* `POST -> /clock/out`
    * Description: Use to clock out
    * Authentication inside Header: 
        ```
        Authorization: Bearer <JWT Token>
        ```
    * Response:
        - 201 Created: 
        ```
        {
            "message": "user clocked out for today successfully"
        }
        ```
        - 400 Bad Request
        - 401 Unauthorized
        - 500 Internal Server Error

* `GET -> /clock`
    * Description: Use to get all user's punch cards
    * Authentication inside Header: 
        ```
        Authorization: Bearer <JWT Token>
        ```
    * Response:
        - 200 OK: 
        ```
        [
            {
                "clock_in": "2025-02-07T12:00:00+08:00",
                "clock_out": "2025-02-07T20:00:00+08:00",
                "created_at": "2025-02-07T12:00:00+08:00"
            },
            {
                "clock_in": "2025-02-06T12:00:00+08:00",
                "clock_out": "2025-02-06T20:00:00+08:00",
                "created_at": "2025-02-06T12:00:00+08:00"
            }
        ]
        ```
        - 401 Unauthorized
        - 500 Internal Server Error
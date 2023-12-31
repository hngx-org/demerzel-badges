>**Link to Postman API Documentation:** https://documenter.getpostman.com/view/4194134/2s9YJgU1Ez

# Team Demerzel Badge API Documentation

## Table of contents

* [Introduction](#introduction)
   
* [API Usage and Features](#api-usage-and-features)
   * [How to Call the API](#how-to-call-the-api)
   * [Authenticating to the API](#authenticating-to-the-api)
   <!-- * [Security Definitions](#security-definitions) -->

* [Request and Response Formats](#request-and-response-format)
   * [Request](#request)
   * [Reseponse](#response)
   * [Error Formats](#error-formats)

* [API Endpoints](#api-endpoints)
   * [API Health](#api-health)
   * [Badges](#badges)
  
## API Usage and Features
The Badges API is a service which the Zuri Portfolio calls to assign Badges to users.
These badges are given to users after passing a skill assessment test. Admins have
the ability to create badges to be assigned to users of the Zuri Portfolio.

### How to Call the API

The API can be accessed via HTTP requests. It exposes endpoints for CRUD
operations on Badges for Admins, Assigning Badges to Users, and also retrieval 
of a user's badge for a particular skill.
<!-- todo: change this to the current host url. -->
* **Current (Active) Host**:   
* **API Base Path**: `/api`

## Request and Response Format
### Request
<!-- * What Authentication token should be sent along side the request -->

### Response
The Api Response body follow the JSend format, whcih have a `status`, `data` or `error` and `message` key, the status falls under either `success` or `error` respectfully.
The data and error field is a JSON object, the error objects contains Form input
validation errors.  
Body:  
```Json
{
   "status": "success",
   "message": "User Badge Retrieved Successfully",
   "data": {
      "id": 123,
      "assessment_id": 321,
      "user_id": "a2218d8f-4cdb-4114-a847-4cf8fcbd2e54",
      "badge_id": 324,
      "created_at": "2023-09-20T18:28:42.523+01:00",
      "updated_at": "2023-09-20T18:28:42.523+01:00"
   }
}
```
Read more at:  [The JSend Specification](https://github.com/omniti-labs/jsend)

### Error Formats
1. Status Code: 404  
   Example:
   ```Json
   {
      "status": "error",
      "message": "User Badge not Found",
      "errors": {}
   }
   ```
2. Status Code: 400  
   Example:
   ```Json
   {
      "status": "error"
      "errors": {},
      "message": "Unable to parse payload: invalid character '/' looking for beginning of object key string",
   }
   ```
3. Status Code: 422  
   Example:
   ```Json
   {
      "errors": {
         "min_score": "min_score should be at least 0"
      },
      "message": "Invalid input",
      "status": "error"
   }
   ```
4. Status Code: 500  
   Example:
   ```Json
   {
      "status": "error",
      "message": "Unable to create badge",
      "errors": {
         ...
      }
   }
   ```

## API Endpoints
### API Health
* **GET /health**
   * **Sample Request URL**: `{host}/health `
   * **Response**:  
   Status Code: 200  
   Body:
      ```Json
      {
         "data": null,
         "message": "Team Demerzel Events API",
         "status": "success"
      }
      ```

### Badges
* **POST api/badges**
   * **Summary**: Create a Badge
   * **Description**: Create a badge for user after assessment by admin.
   * **Sample Request URL**: `{host}/api/badges`
   * **Parameters**:  
      Body:
      ```Json
      {
         "skill_id": 321,
         "name": "Intermediate",
         "min_score": 51,
         "max_score": 80
      }
      ```
   * **Response**:   
   Status Code: 201  
   Body:
      ```Json
      {
         "status": "success",
         "message": "Badge Created Successfully",
         "data": {
            "id":123,
            "skill_id": 321,
            "name": "Intermediate",
            "min_score": "51",
            "max_score": "80"
         }
      }
      ```
* **POST /api/user/badges**
   * **Summary**: Assign Badge to a user
   * **Description**: After a user has passes an assessment, assign a badge to the user, provide
   the required fields in the request body.
   * **Sample Request URL**: `{host}/api/user/badges`
   * **Parameters**:
      Body:
      ```Json
      {
         "user_id": "a2218d8f-4cdb-4114-a847-4cf8fcbd2e54",
         "badge_id": 123,
         "assessment_id": 321,
         "skill_id": 432
      }
      ```
   * **Response**:  
      Status Code: 201    
      Body:
      ```Json
      {
         "status": "success",
         "message": "Badge Assigned Successfully",
         "data": {
            "id": 123,
            "skill_id": 432,
            "badge_id": 123,
            "assessment_id": 321,
            "created_at": "2023-09-20T18:28:42.523+01:00",
            "updated_at": "2023-09-20T18:28:42.523+01:00"
         }
      }
      ```

* **GET /api/user/badges/{userId}/skill/{skillId}**
   * **Summary**: Retrive Badge of a user for a particular skill
   * **Sample Request URL**: `{host}/api/user/badges/a2218d8f-4cdb-4114-a847-4cf8fcbd/skill/123
   * **Response**:  
      Status Code: 200  
      Body:
      ```Json
      {
         "status": "success",
         "message": "User Badge Retrieved Successfully",
         "data": {
            "id": 123,
            "skill_id": 432,
            "badge_id": 123,
            "assessment_id": 321,
            "created_at": "2023-09-20T18:28:42.523+01:00",
            "updated_at": "2023-09-20T18:28:42.523+01:00"
         }
      }
      ```
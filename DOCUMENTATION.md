Document API endpoints, request/response structures, and usage in API documentation.

# Team Demerzel Badge API Documentation

## Table of contents

* [Introduction](#introduction)
   
* [API Usage and Features](#api-usage-and-features)
   * [How to Call the API](#how-to-call-the-api)
   * [Authenticating to the API](#authenticating-to-the-api)
   <!-- * [Security Definitions](#security-definitions) -->

* [API Endpoints](#api-endpoints)
   * [API Health](#api-health)
   * [Authentication](#authentication)
   * [Groups](#groups)
   * [Users](#users)
   * [Events](#events)
   * [Comments]()
   * [Images](#images)
  
* [Request and Response Formats](#request-and-response-format)

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
* What Authentication token should be sent along side the request

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
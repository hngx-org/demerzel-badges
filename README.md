# Team Demerzel

## About Project
The Badges API is a service which the Zuri Portfolio calls to assign Badges to users.
These badges are given to users after passing a skill assessment test. Admins have
the ability to create badges to be assigned to users of the Zuri Portfolio.

## Stack and Technologies used
- [Golang](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
## Prerequisites 
Golang v >= 1.19  
PostgreSQL v >= 15.0
## Building and running locally
1. Install Go and PostgreSQL
2. Clone the repo at https://github.com/hngx-org/demerzel-badges
3. `cd` into the project directory
4. Create a .env file and fill PORT and PostgreSQL details using the .env.example format.  
`cp .env.example .env`
```bash
PORT=3001

POSTGRES_USERNAME=
POSTGRES_PASSWORD=
POSTGRES_HOST=
POSTGRES_DBNAME=
POSTGRES_PORT=
```
5. Run `go get`
6. Run `go run main.go`
6. Access the API endpoints from localhost with the port specified in step 4.
7. Read the Documentation.md to check for available endpoints.
    
## As a maintainer

### Fork repo to personal github account
Your Repo's found at https://github.com/hngx-org/demerzel-badges
So to work with Forks you basically:
1. Fork your Team Repo to your personal Github account
2. Pull the code back to your local Machine
3. Checkout to your assignment branch
4. Do your thing
5. Push back to your Personal Github Repo. That'll be your 'ORIGIN' remote (not the 'UPSTREAM' remote)
6. You head over to Github and Create a Pull request to the Main Repository's branch
* Remember that a Pull Request can contain multiple commits.
That's basically it. Here's a [video](https://youtu.be/nT8KGYVurIU) to help further.
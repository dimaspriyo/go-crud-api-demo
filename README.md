# Go CRUD API Demo
Simple CRUD with 2 routes group 
- v1 is just plain route without authentication
- v2 using JWT authentication

### v1 
1. "/v1" - GET
2. "/v1/insert" - POST
3. "/v1/update:id" - PUT
4. "/v1/delete:id" - DELETE


### v2 
1. "/v2/token" - GET (Get token first)
2. "/v2" - GET
3. "/v2/insert" - POST
4. "/v2/update/:id" - PUT
5. "/v2/delete/:id" - DELETE
> Use Authorization: Bearer  **your token** 


## 

Insert Data

    {
    "name" : "..."
    "sex" : "..."
    "country" : "..."
    }

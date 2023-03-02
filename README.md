## GO API with Echo Framework
*Create an API using GO and the Echo Framework.*

To Start Server from root:
`go run echo/main.go`


Requirements 
1. GO v1.18
2. Echo Framework
3. Ginko/Gomega Testing
4. Squirrel for DB access

Notes:

The tutorial I used to acclimate to Go used the storm DB framework and used BSON for ID generation. 
Once a Postgres DB is created, the switch to the squirrel framework would be trivial.

![image](https://user-images.githubusercontent.com/6537603/222530430-6d05278a-2076-43e9-914b-a64119a9a55d.png)

I wrote the APIs and tests by hand in handlers/ and provided contingency for user_name uniqueness. All tests pass:

> ##### user_test.go
> ![image](https://user-images.githubusercontent.com/6537603/222533838-e6f904de-bb07-4b53-9454-b8b72fd399c4.png)
> ##### userHandlers_test.go
> ![image](https://user-images.githubusercontent.com/6537603/222534295-c3482d86-a6d8-4ea0-ab87-b2a15845cdb1.png)

I also provided caching and a custom writer for intercepting requests. I wrote tests for these:

> ##### writer_test.go
> ![image](https://user-images.githubusercontent.com/6537603/222534746-9da43049-c559-4535-a918-90946b0eb629.png)

The ported all requests and handlers to Echo. 
I still need to add Ginko/Gomega tests to Echo implementation. 

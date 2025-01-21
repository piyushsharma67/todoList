# todolist Interview Programming Task

- You have full internet access and may research the solution as you see fit.
- You should not consult another individual about your solution.

After the task there will be a discussion on the design and implementation of the system.


You have been provided with a simple REST API service which implements a todo list (for a single user).
This lets a user keep a simple list of things they need to do, e.g. wash car, fix bike, etc.

-   The service is designed to back a website (the fictional website is not included).
-   The website should be kept simple, with logic in the backend.
-	There are endpoints for Creating, Removing, Updating and Deleting an item in this list.
-	There is an endpoint for retrieving the entire list.
-	The service is backed by an SQLite database.
-	There are unit tests to ensure the service functions correctly.

It is desired to allow the user to order the items in this list. 
-   The order of the items should be persisted.
-   No two Items in the list should have the same order.
-   It should be possible to change the order of the items in the list (e.g. the fictional website allows dragging + dropping or moving the items).

You should: -
-   Expand the current service to accomodate ordering items.
-	Design REST APIs to allow the user to change the order in the list.
-	Desing the DB structure need to accommodate the ordering system.
-	Implement the ordering system.
-	Add unit tests to ensure the behaviour is correct.

My Solution: -

- created a client page where he could add task so you can run it by clicking over it for a interface
- i have used air in order to hot reload when i make changes to the files
- i have added priority as a column which is by default incremented 
- for now user could make changes to the priority, lower the number higher the priority and list fetched is sorted on basis of priority
- in case priority clashes then the list is sorted on the basis of updated_at column
- since i have added 3 column i have accordingly made changes to the test cases for the api endpoints and db queries.

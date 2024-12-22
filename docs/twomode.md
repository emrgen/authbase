# Two modes of execution

The authbase service originated from the need of a simple authentication system.

The main driving factor was price.

Initially I started using firebase, supabase.
Having the authentication service in different place than the main service have some problems.

1. Speed
2. Dependency
3. Future lock in

Though very easy to start with it felt someday I have to pay heavy price.

The initial goal was to create a multi tenant auth service like simplified keycloak.

The design is simple

## Backend

1. Create one database.
2. Create a master org and super user at start
3. Master org member(with write permission) can create more orgs with a admin user (also called org member)

## Frontend

1. Organization member logs in the authbase ui
2. Get the user view
3. Add new user
4. ...

The above design serves well when all the users are in one database.

But what if different projects wants to use different database but the same authbase service.

**NOTE**: Authbase backend is stateless and all the data is kept in the database.

The solution is to use project specific database. It works as follows

1. On start no master or is created.
2. A external user send request to create a new org in a project.
3. The user token is verified by an external service. (most probably another authbase service running in singlestore mode)
4. Once authenticated the external users write access in the project is ensured by another permission check api call. (calls external permission management system with <projectID#write@userID>)
5. If the user has permission. a new org is created.
6. If this is the first org in the database, it is marked as master.

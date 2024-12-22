# Stories

## Create org at start up

- [x] done

- **Modes**
  - In singlestore mode the service should create the master org and the super admin on startup
  - In multistore mode the service should start but no master org exists

## Create token for a org

- [x] done

- **Requirement**
  - To create the organization offline token the user must be a member of the master org or the target org.
- **Result**
  - The token will have same permission as the creator.

## Add a new user to a organization

- **Requirement**
  - To add a new user to the organization the caller must be a member of the master org or the target org.
- **Result**
  - A new user is created with provided password.
  - A email verification link is mailed if the mailer is configured

### User registration in an organization

- **Requirement**
  - User calls the register api to create a new account with username/password
  - Need to provide the organization id in the payload
- **Result**
  - A new user is created within the organization

### User login

- **Requirement**
  - User need to login using username/password
- **Result**
  - Users logs in with valid credentials
  - A auth+refresh token pair is returned

### User logout

- **Requirement**
- **Result**

### User reset password

- **Requirement**
- **Result**

### User forgot password

- **Requirement**
- **Result**

### User verify email

- **Requirement**
- **Result**

### User login via oauth

- **Requirement**
  - User login via oauth provider like google
- **Result**
  - User is logged in
  - Access tokens are returned
  - A HttpOnly cookie is returned with the response

### Add oauth provider for a organization

- **Requirement**
  - A user with permissions adds a new oauth provider
  - A organization should have at most one config per provider.
- **Result**

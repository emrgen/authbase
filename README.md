# authbase is a simple authentication service.

**NOTE**: This is a side project. Not to be meant for production.

## API Usage

```go
package main

import (
    "context"
    "log"

    "github.com/emrgen/authbase"
    "github.com/emrgen/authbase/apis/v1"
)

func main() {
    // Create a new client.
    client, err := authbase.NewClient("localhost:4000")
    if err != nil {
        log.Fatal(err)
    }

    // Create a new project. need token to create an organization.
    project, err := client.CreateProject(context.Background(), &v1.CreateProjectRequest{
        Name: "org",
        Email:    "admin@gmail.com",
        Password: "password",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Create a new user. (the user is crated within the organization)
    user, err := client.CreateAccount(context.Background(), &v1.CreateAccessKeyRequest{
        Email:    "user@mail.com",
        Passowrd: "password",
		PoolID: "pool-id",
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## Goal

1. The goal of this library is to provide a simple way to authenticate users in a web application.
2. Support multitenancy.
3. Support multiple authentication methods.
4. Support multiple storage backends.
5. Minimal permission system. Users within an org can have different permissions once they are members of the org.
6. Support for user registration, password reset, and email verification etc.

## Installation

```bash
# Install the buf dependencies.
make buf-deps
# Generate the protobuf files.
make protoc
# Install the go dependencies.
make deps
# Run the server.
make air
```

## CLI Usage

```bash
# Create master organization. (The first organization is the master organization)
# if the password is provided, the email verification will not be required strictly but need to verify later.
# without the password, if the email is not verified in 10 minutes, the organization will be deleted.
# NOTE: need to set the email config before creating the organization.
authbase org create --name=master --user=admin --email=email [--password=password] [--verify=true]

# in case there is no org no need to set the org flag
authbase config code set --medium=email --value=smtp://user:password@localhost:587

# Create a user token
authbase token create ---p=project-id-id --user=admin --password=admin

# Set context
authbase context set --token=token ---p=project-id-id

############################################
# Organization commands
############################################

# Create an organization
authbase org create --name=org

# Get an organization
authbase org get --name=org

# List all organizations
authbase org list --limit=10 --offset=0

# Delete an organization
authbase org delete --name=org

# Add database to an organization 
# with separate db set, on migrate the organization will be migrated to the database
authbase org db set ---p=project-id --db-url=postgres://user:password@localhost:5432/db

# Add migration db to an organization
authbase org migration db set ---p=project-id --db-url=postgres://user:password@localhost:5432/db

# Migrate an organization
authbase org migrate ---p=project-id

############################################
# Provider commands
############################################

# Add a oauth provider
# adding a prover automatically enables the strategy=oauth2
authbase provider add --name=org --provider=google --client-id=client-id --client-secret=client-secret

# Enable login strategy
authbase strategy enable ---p=project-id --strategy=password
authbase strategy enable ---p=project-id --strategy=oath2

# Disable login strategy
authbase strategy disable ---p=project-id --strategy=password
authbase strategy disable ---p=project-id --strategy=oath2

# Get a provider
authbase provider get ---p=project-id --provider=google

# List all providers
authbase provider list ---p=project-id

# Remove a oauth provider
authbase provider remove ---p=project-id --provider=google

############################################
# Config commands
############################################

# Set a email config
authbase config code set ---p=project-id --medium=email --value=smtp://user:password@localhost:587

# Get a email config
authbase config code get ---p=project-id --medium=email

# Remove a email config
authbase config code remove ---p=project-id --medium=email

# Set otp config
authbase config opt set ---p=project-id --key=phone --value=twilio://account-sid:auth-token@localhost:8080

############################################
# Member commands
############################################

# Create a member
authbase member create --username=user --email=example@mail.com --verify=true

# Get a member
authbase member get --username=user

# List all members
authbase member list --limit=10 --offset=0

# Delete a member
authbase member delete --username=user

# Add permission to a member
authbase user permission add --token=token --username=user --permission=read

# Remove permission from a member
authbase user permission remove --token=token --username=user --permission=read

# Verify a member
authbase member verify --username=user --code=123456

# Reset a member password
authbase member password reset --username=user --password=123456

# Send a member password reset email
authbase member password reset email --username=user

############################################
# User commands
############################################

# Create a user
authbase user create --username=user --email=example@mail.com --verify=true

# Get a user
authbase user get --username=user

# List all users
authbase user list --limit=10 --offset=0

# Delete a user
authbase user delete --username=user

# Verify a user
authbase user verify --username=user --code=123456

# Send a user password reset email
authbase user password reset email --username=user

# Reset a user password
authbase user password reset --username=user --password=123456
```

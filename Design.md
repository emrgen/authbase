# Design

## Use case 1

### One db to for orgs

#### Master org
In multitenant applications, we need to create a master organization that will have the ability to create other organizations.
The master org members will have the ability to create other organizations and manage them.

#### Other org
The other organizations will be created by the master org members with write permission. 
The organization will have one admin after creation. The org admin can not be removed from the organization.
The org admin will have the ability to create other users and manage the organization.
The admin can also create other admins by promoting the users to admin role.

## Use case 2

### Separate db for each org, no master org

NOTE: This use case is experimental and not implemented yet.

#### Org

In this use case, there will be no master organization. Each organization will have its own database.
The permission are not saved in the database, the permission check will be delegated to the authbac service.

The permission check to work the permission are updated when the organization, op0
00

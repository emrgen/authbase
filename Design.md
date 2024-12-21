# Design

## Use case 1

### One db to all orgs

#### Master org
In multitenant applications, we need to create a master organization that will have the ability to create other organizations.
The master org members will have the ability to create other organizations and manage them.

#### Other org
The other organizations will be created by the master org members. The organization will have one admin during creating. 
The admin will have the ability to create other users and manage the organization.
The admin can also create other admins by promoting the user to admin.

## Use case 2

### Separate db for each org, no master org

#### Org

In this use case, there will be no master organization. Each organization will have its own database.
The permission are not saved in the database, the permission check will be delegated to the authbac service.

The permission check to work the permission are updated when the organization, op0
00

## Use case 3

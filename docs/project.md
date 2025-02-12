# Everything falls under a project.

## Introduction

A project is a collection of user pools. A user pool is a collection of users. A user can be part of multiple user pools with different account. 
A user pool can have multiple groups. A group has some permissions. The user of a group inherits these permissions.

## Project

A project is a collection of user pools. A project is created with a name and a description. 
The project has a master user pool. The master user pool has a master user. The master user has all the permissions in the project.

## Project members

Project members are users who have access to the project. A project member can have different roles in the project.
Project members must be part of the default user pool. 

A project member can have the following roles:

- `admin`: An admin can create user pools, groups, and users. An admin can also delete user pools, groups, and users.
- `writer`: A user can create groups and users. A user can also delete groups and users.
- `viewer`: A viewer can only view the project.

A admin can create a user pool with the user pool master. The user pool master has all the permissions in the user pool.

## User pool 

A user pool is a collection of users. A user can be part of multiple user pools with different account.
A user pool has a name and a description. A user pool has a master user. The master user has all the permissions in the user pool.

## User pool members

User pool members are users who have access to the user pool. A user pool member can have different roles in the user pool.

A user pool member can have the following roles:

- `admin`: An admin can create groups and users. An admin can also delete groups and users.
- `writer`: A user can create groups and users. A user can also delete groups and users.
- `viewer`: A viewer can only view the user pool.

## Group

A group is a collection of users. A group has some permissions. The user of a group inherits these permissions.

## Group members

Group members are users who have access to the group. A group member can have different roles in the group.

## Role

A role is a permission that can be assigned to a user, a group. Example of roles are `document.read`, `document.write`, `document.delete`.

## Resource

A resource is an entity that can be accessed by a user. A resource can be a file, a folder, a project, a user pool, or a group.

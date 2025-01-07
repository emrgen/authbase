# Predefined scopes for accounts using groups

## Introduction

For many services, the user is the main entity. The user logs in and performs actions on the service. But in some cases, the user is part of a group. The group has some permissions and the user inherits these permissions.

For example, in a company, the user is part of a team. The team has some permissions and the user inherits these permissions. The user can also have some individual permissions.

In this document, we will discuss how to manage groups and permissions in the authbase service.

# Group

A group is a collection of users. The group has some permissions. The user inherits these permissions.
A group is created in the context of a user pool. The user pool is a collection of users

## User pool

A project can have multiple user pools. When a project is created a default user pool is created along with pool members. The user pool can have multiple groups.

## Create a group

To create a group, the user must have the `group.create` permission.

The user sends a request to create a group. The group service verifies the token validity and the user permission. If the token is valid and the user has the `group.create` permission, the group is created.

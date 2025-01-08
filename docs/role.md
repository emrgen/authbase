# To use or not to use the authbase role

## When not to use authbase role

The `authbase` role is not suitable for use in the following cases:

- When a small number of account holders need to access the application.
- Accounts has wide access to the resources resulting few scopes in the jwt token.
- The resource and user hierarchy is simple and does not require fine-grained access control.

## When not to use authbase role

The `authbase` role is not suitable for use in the following cases:

- When a number of accounts and resources are huge and need fine-grained access control. 
- When the user scopes will be too many to add to the jwt token.
- When the resource and user hierarchy is complex and requires fine-grained access control.


## Scenarios

Imagine a scenario where few users control a vast number of resources. 
We need to check if the user has access to a particular resource category.
In this case, the `authbase` role is suitable as the user scopes are limited and we dont care about the resource hierarchy.

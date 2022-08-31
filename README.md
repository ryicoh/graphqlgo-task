# Task management app using grahpql-go

```
{
  User(where: {userID: {_eq: "4624DBD4-F795-4FD9-9C02-3FB0539B6808"}}) {
    name
    tasks {
      taskID
      name
    }
  }
}

=>

{
  "data": {
    "User": {
      "name": "user1taro",
      "tasks": [
        {
          "name": "task1todo",
          "taskID": "ac1eba9d-8e63-4e8c-9b89-365f718d92e5"
        },
        {
          "name": "task2todo",
          "taskID": "e48a1f07-2de5-49cb-a31b-f69fd0466c80"
        }
      ]
    }
  }
}
```
# graphqlgo-task

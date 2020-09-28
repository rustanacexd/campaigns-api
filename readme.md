# campaigns-api

> A simple REST API built as sample project

* No third-party packages/dependencies

To be able to show the desired features of curl this REST API must match a few
requirements:

* [x] `GET /campaigns/` returns list of campaigns as JSON
* [x] `GET /campaigns/{id}/` returns details of specific campaign as JSON
* [x] `POST /campaigns/` accepts a new campaign to be added
* [x] `POST /campaigns/` returns status 415 if content is not `application/json`
* [x] `PUT /campaigns/{id}/` returns status 415 if content is not `application/json`
* [x] `PUT /campaigns/{id}/` update an existing campaign
* [x] `DELETE /campaigns/{id}/` delete an existing campaign



### Data Types

A campaign object should look like this:
```json
{
  "ID": 1,
  "Name": "name of the campaign",
  "Status": "some status",
  "Created": "2020-09-28T15:07:17.976388+08:00",
}
```

### Persistence

There is no persistence, a temporary in-mem story is fine.

# Budget Control Application

This is teh project in whit the BCA Backend App will be developed and deployed in.

## Technical stack

![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

## Main Features

- Create projects to manage
- Create a project budget by adding budget items to it
- Register the invoices, each invoice shall decrease the values from the budget

## Environment Variables

In order to configure this project, the following environment variables are required:

- PGDATABASE
- PGHOST
- PGPASSWORD
- PGPORT
- PGUSER
- SECRET

## API Routes

- **/** this is the home route, for the time being is a route that will check if the server is working

### Authentication Endpoints

- **/login** Using a POST request to this is the endpoint where the user will provide their credentials and log in
- **/api/v1/logout** Using a GET request to this endpoint, the user will be logged out from the application

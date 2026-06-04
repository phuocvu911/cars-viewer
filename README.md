# Run project from root 

```bash 
go run run.go 
```

It is a startup script for the JS server as well as the go program. If either fails, the project will shut down. 

# Task Split

Hoang does the `:8080/` root and Markus does `/car/{id}` page. 

Hoang do `/compare`.

Hoang go on with `/gallery` page with search option and advanced filter while Markus do `/recommendation`.

# Run the project 

## Clone the project

```bash 
git clone https://gitea.kood.tech/hoangphuocvu/viewer
cd viewer
```

## Run api server
```bash 
make api
```

## Run Go backend server
Open the new terminal. On the new terminal:
```bash 
make run
```

# Overview
...

# Extras

## Search Options And Advanced Filter
This feature shipped on top of Gallery page. The user can do free-word search, i.e `au` can return `Audi A4` car.

The user can filter cars view by clicking the drop down to choose `Categories`, `Brand`, `Year` and `Drivetrain`.

The filter bar send `GET` request to `/gallery` with data as query parameters. This is right choice for retrieving/filtering data — these are read-only operations that don't modify server state.

## Comparision Feature
Selecting cars and clicking ``Compare`` button sends `POST /compare` with selected IDs
## Recommendation Feature

## Store Analytics 
Some data analytics about the cars that we have in the store.

## Auto Refreshing Data

There is a `go routines` that running in the background to update the cars data every 10 minutes. The user can see the latest data when they refresh the page.

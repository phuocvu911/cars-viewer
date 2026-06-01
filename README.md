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

## Search options and advanced filter
This feature shipped on top of Gallery page. The user can do free-word search, i.e `au` can return `Audi A4` car.

The user can filter cars view by clicking the drop down to choose `Categories`, `Brand`, `Year` and `Drivetrain`.

## Comparision feature
Selecting cars and clicking ``Compare`` button sends `POST /compare` with selected IDs
## Recommendation feature

## Page analytics
Some data analytics about the cars that we have in the store.

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
Selecting cars and clicking ``Compare`` button sends `GET /compare` with selected IDs to retrieve detailed information about those cars for side-by-side comparison.


## Recommendation Feature

The recommendation feature is based on cookies saved to the client. The website recommends cars to the user based on the most clicked car brand as well as the most clicked chassis type. Clicking lots of sedan audis will give you also bmw sedans as a result. Then if you visit lots of Ford ads on the site, the website might also recommend e.g. Ford pickup trucks.

*The user needs visit accumulate 2 same brands or 2 same chassis types to start receiving recommendations, with the default settings.* 

The cookie flow is made in combination with the browser and the backend. The browser prompts for the cookies if it doesn't contain a right named cookie set as "true" or as "false". If the requests made to the backend contains invalid data, the cookies will be deleted and the consent will be prompted again. The website does not have any registry for the cookies given to the clients so verifying the cookie values is basically just that the contains text. (due to validating cookies is out of scope for this project)

## Store Analytics 
Some data analytics about the cars that we have in the store. 

## Auto Refreshing Data

There is a `go routines` that running in the background to update the cars data every 10 minutes. The user can see the latest data when they refresh the page.

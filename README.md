# TODO 

- Create a 404 page
- Logic for validating illegal paths 


# Requirements

- Go 1.26+

- Node 24+

- NPM 11+ (included with node 24)

# Run project from root 

## Clone the project

```bash 
git clone https://gitea.kood.tech/hoangphuocvu/viewer
cd viewer
```
## Start the Servers

API server:

```bash 
cd cars-api && make build && make run
```

On the **new terminal** (in the root directory `./viewer`):

Backend server:

```bash 
cd go-backend && go run main.go
```
AutoVault is ready to see at http://localhost:8080/
# Overview

**AutoVault** is a car catalogue web app where you can browse, search, filter, compare and get recommendations for cars. It is built from two services:

- **`cars-api`** — a small Node/Express server (`:3000`) that serves the raw car dataset (`data.json`) and the car images. It exposes read-only endpoints for models, manufacturers, categories and images.
- **`go-backend`** — a Go server (`:8080`) that renders the actual website. It fetches data from the API, *enriches* each car (joining its manufacturer and category details), caches the result, and serves the HTML pages. The frontend never talks to the API directly — image requests are reverse-proxied through `/api/images/`.

## Architecture

```
Browser ──▶ go-backend (:8080)  ──▶ cars-api (:3000)
            HTML pages, caching,      raw data.json
            recommendations,          + car images
            analytics
```

The Go backend keeps an in-memory **store** of all car models and categories. On startup it loads the store from the API, and a background goroutine refreshes it every 10 minutes so the data stays current without restarting the server. Individual car details are enriched concurrently (manufacturer + category fetched in parallel) and cached to keep page loads fast.

## Pages

- **Home** (`/`) — landing page.
- **Gallery** (`/gallery`) — the main view: a grid of all cars with free-text **search**, **filters** (category, brand, year, drivetrain) and a personalised **recommendations** section.
- **Car details** (`/car/{id}`) — full specifications for a single car.
- **Compare** (`/compare`) — side-by-side comparison of the cars you select in the gallery.
- **Stats** (`/stats`) — store-wide analytics about the cars in the catalogue.

Recommendations are driven by **cookies**: clicks are tracked per visitor and used to suggest cars matching the brands and chassis types they browse most. See [Recommendation Feature](#recommendation-feature) below for the full flow.

# Extras

## Search Options And Advanced Filter
This feature shipped on top of Gallery page. The user can do free-word search, i.e `au` can return `Audi A4` car.

The user can filter cars view by clicking the drop down to choose `Categories`, `Brand`, `Year` and `Drivetrain`.

The filter bar send `GET` request to `/gallery` with data as query parameters. This is right choice for retrieving/filtering data — these are read-only operations that don't modify server state.

## Comparision Feature
Selecting cars and clicking ``Compare`` button sends `GET /compare` with selected IDs to retrieve detailed information about those cars for side-by-side comparison.


## Recommendation Feature

This feature embedded into of Gallery page.For the first time user, clicking on some cars and then refreshing the gallery page will show the recommendation section.

The recommendation feature is based on cookies saved to the client. The website recommends cars to the user based on the most clicked car brand as well as the most clicked chassis type. Clicking lots of sedan audis will give you also bmw sedans as a result. Then if you visit lots of Ford ads on the site, the website might also recommend e.g. Ford pickup trucks.

*The user needs visit accumulate 2 same brands or 2 same chassis types to start receiving recommendations, with the default settings.* 

The cookie flow is made in combination with the browser and the backend. The browser prompts for the cookies if it doesn't contain a right named cookie set as "true" or as "false". If the requests made to the backend contains invalid data, the cookies will be deleted and the consent will be prompted again. The website does not have any registry for the cookies given to the clients so verifying the cookie values is basically just that the contains text. (due to validating cookies is out of scope for this project)

## Store Analytics 
Some data analytics about the cars that we have in the store. 

## Auto Refreshing Data

There is a `go routines` that running in the background to update the cars data every 10 minutes. The user can see the latest data when they refresh the page.

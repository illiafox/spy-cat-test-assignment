# Spy Cats API

A RESTful API for managing spy cats, their missions, and targets.


<details>

<b><summary>Test Assignment Description (click to expand)</summary></b>

> ## Overview
> This task involves building of a a simple yet comprehensive CRUD application. The goal is to create a simple yet comprehensive system that demonstrates your understanding in building RESTful APIs, interacting with SQL-like databases, integrating third-party services, and optionally creating user interfaces. The test assessment is expected to be done under 4 hours.
> ## Business Task
> Spy Cat Agency (SCA) asked you to create a management application, so that it simplifies their spying work processes. SCA needs a system to manage their cats, missions they undertake, and targets they are assigned to.
> From cats perspective, a mission consists of spying on targets and collecting data. One cat can only have one mission at a time, and a mission assumes a range of targets (minimum: 1, maximum: 3). While spying, cats should be able to share the collected data into the system by writing notes on a specific target. Cats will be updating their notes from time to time and eventually mark the target as complete. If the target is complete, notes should be frozen, i.e. cats should not be able to update them in any way. After completing all of the targets, the mission is marked as completed.
> From agency perspective, they regularly hire new spy cats and so should be able to add them to  and visualize in the system. SCA should be able to create new missions and then assign them to cats that are available. Targets are created in place along with a mission, meaning that there will be no page to see/create all/individual targets. A target can be added to or deleted from a mission.
> ## Backend Requirements
>   **- Spy Cats**
>   - Ability to create a spy cat in the system
>     - A cat is described as Name, Years of Experience, Breed, and Salary.
>     - Breed must be validated, see General
>     - Ability to remove spy cats from the system
>     - Ability to update spy cats’ information (Salary)
>     - Ability to list spy cats
>     - Ability to get a single spy cat
>- **Missions / Targets**
>  - Ability to create a mission in the system along with targets
>    - A mission contains information about Cat, Targets and Complete state
>    - Each target is unique to a mission, so the endpoint should accept an object describing targets
>    - A target is described as Name, Country, Notes and Complete state
>  - Ability to delete a mission
>    - A mission cannot be deleted if it is already assigned to a cat
>  - Ability to update mission
>    - Ability to mark it as completed
>  - Ability to update mission targets
>    - Ability to mark it as completed
>    - Ability to update Notes
>      - Notes cannot be updated if either the target or the mission is completed
>  - Ability to delete targets from an existing mission
>    - A target cannot be deleted if it is already completed
>  - Ability to add targets to an existing mission
>    - A target cannot be added if the mission is already completed
>  - Ability to assign a cat to a mission
>  - Ability to list missions
>  - Ability to get a single mission
>- **General**
>  - Framework
>    - Use any modern framework, e.g. Gin, Echo, Fiber, etc.
>  - Database
>    - You can use any SQL-like database
>  - Logger middleware
>    - Integrate a logging middleware to log HTTP requests and responses for debugging and monitoring
>  - Validations
>  - Make sure endpoints validate the request body and return an adequate status code if it’s not valid
>  - Validate cat’s breed with [TheCatAPI](https://api.thecatapi.com/v1/breeds)

</details>

## Table of Contents

- [Installation](#installation)
- [Postman Collection](#postman-collection)
- [API Endpoints](#api-endpoints)
    - [Cats](#cats)
    - [Missions](#missions)
    - [Targets](#targets)
- [Contributing](#contributing)

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/illiafox/spy-cat-test-assignment.git
   cd spy-cat-test-assignment
   ```

2. Start docker compose:

   ```sh
    docker compose up -d
   ```

The API will be available at `http://127.0.0.1:8080`.

Postgres URI: `postgresql://postgres:dev_pass@localhost:5432/spy_cats`.

## Postman Collection

1. Open in browser: https://documenter.getpostman.com/view/36386828/2sA3e5eoet
2. (works better) OR download a [collection.json](collection.json) and [import it](https://learning.postman.com/docs/getting-started/importing-and-exporting/importing-data/)


## API Endpoints

### Cats

- **List All Cats**
    - **GET** `/cats/`
    - Example request: `GET http://127.0.0.1:8080/cats/`

- **Retrieve Cat Info**
    - **GET** `/cats/:id`
    - Example request: `GET http://127.0.0.1:8080/cats/1`

- **Update Cat**
    - **PATCH** `/cats/:id`
    - Example request:
      ```sh
      PATCH http://127.0.0.1:8080/cats/1
      Content-Type: application/json
      {
        "name": "Fiona",
        "experience": 11,
        "salary": 6000
      }
      ```

- **Remove Cat**
    - **DELETE** `/cats/:id`
    - Example request: `DELETE http://127.0.0.1:8080/cats/1`

- **Add Cat**
    - **POST** `/cats/`
    - Example request:
      ```sh
      POST http://127.0.0.1:8080/cats/
      Content-Type: application/json
      {
        "name": "Minni",
        "breed": "Abyssinian",
        "experience": 10,
        "salary": 5000
      }
      ```

### Missions

- **List All Missions**
    - **GET** `/missions/`
    - Example request: `GET http://127.0.0.1:8080/missions/`

- **Retrieve Mission Info**
    - **GET** `/missions/:id`
    - Example request: `GET http://127.0.0.1:8080/missions/1`

- **Update Mission**
    - **PATCH** `/missions/:id`
    - Example request:
      ```sh
      PATCH http://127.0.0.1:8080/missions/1
      Content-Type: application/json
      {
        "assigned_cat_id": 1
      }
      ```

- **Complete Mission**
    - **POST** `/missions/:id/complete`
    - Example request: `POST http://127.0.0.1:8080/missions/1/complete`

- **Remove Mission**
    - **DELETE** `/missions/:id`
    - Example request: `DELETE http://127.0.0.1:8080/missions/1`

- **Add Mission**
    - **POST** `/missions/`
    - Example request:
      ```sh
      POST http://127.0.0.1:8080/missions/
      Content-Type: application/json
      {
        "targets": [
          {
            "name": "Mister",
            "country": "UA"
          },
          {
            "name": "Second Mister",
            "country": "US"
          }
        ]
      }
      ```

### Targets

- **List Mission Targets**
    - **GET** `/missions/:mission_id/targets`
    - Example request: `GET http://127.0.0.1:8080/missions/1/targets`

- **Add Targets To Mission**
    - **POST** `/missions/:mission_id/targets/add`
    - Example request:
      ```sh
      POST http://127.0.0.1:8080/missions/1/targets/add
      Content-Type: application/json
      {
        "targets": [
          {
            "name": "Another Mister",
            "country": "UK"
          },
          {
            "name": "Last Mister",
            "country": "GE"
          }
        ]
      }
      ```

- **Add Notes To Target**
    - **POST** `/missions/:mission_id/targets/:target_id/notes/add`
    - Example request:
      ```sh
      POST http://127.0.0.1:8080/missions/1/targets/1/notes/add
      Content-Type: application/json
      {
        "notes": [
          "He lives somewhere near Times Square, more investigation is required",
          "I wonder whether he loves drinking coffee"
        ]
      }
      ```

- **Retrieve Target Info**
    - **GET** `/missions/:mission_id/targets/:target_id/`
    - Example request: `GET http://127.0.0.1:8080/missions/6/targets/11/`

- **Complete Target**
    - **POST** `/missions/:mission_id/targets/:target_id/complete`
    - Example request: `POST http://127.0.0.1:8080/missions/6/targets/13/complete`

- **Remove Target**
    - **DELETE** `/missions/:mission_id/targets/:target_id`
    - Example request: `DELETE http://127.0.0.1:8080/missions/6/targets/11`


# Contributing
Please refer to [CONTRIBUTING.md](CONTRIBUTING.md) 
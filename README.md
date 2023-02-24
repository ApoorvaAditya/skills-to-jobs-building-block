# Skill to Jobs Building Block
Building block to handle the backend for the skills to occupations/jobs feature.

## Architecture
The service is based on clear hexagonal architecture. The hexagonal architecture divides the system into several loosely-coupled components, such as the application core and different adapters. We categorize the adapters in two categories - driver and driven.

### core
This is the core component of the service. It keeps the data model and the logic of the service. It communicates with the outer world via adapters.

### driver adapters
What the service provides - user interface, rest adapter, test agent etc.

### driven adapters
What the service depends on - database, mock database, integration with other services etc.

## Documentation
The functionality provided by this application is documented in the [Wiki](https://github.com/rokwire/skills-to-jobs-building-block/wiki).

The API documentation is available here: https://api.rokwire.illinois.edu/skills-to-jobs/api/doc/ui/index.html

## Set Up

### Prerequisites
MongoDB v4.2.2+

Go v1.19+

### Environment variables
The following Environment variables are supported. The service will not start unless those marked as Required are supplied.

Name|Format|Required|Description|Default
---|---|---|---|---
SKILLS_TO_JOBS_BASE_URL | < url > | yes | Base URL where this application is being hosted
SKILLS_TO_JOBS_PORT | < int > | no | Port to be used by this application | 80
SKILLS_TO_JOBS_MONGO_AUTH | <mongodb://USER:PASSWORD@HOST:PORT/DATABASE NAME> | yes | MongoDB authentication string. The user must have read/write privileges.
SKILLS_TO_JOBS_MONGO_DATABASE | < string > | yes | MongoDB database name
SKILLS_TO_JOBS_MONGO_TIMEOUT | < int > | no | MongoDB timeout in milliseconds | 500
SKILLS_TO_JOBS_CORE_BB_BASE_URL | < url > | yes | Core BB base URL

### Run Application

#### Run locally without Docker

1. Clone the repo (outside GOPATH)

2. Open the terminal and go to the root folder
  
3. Make the project  
```
$ make
...
▶ building executable(s)… 1.9.0 2020-08-13T10:00:00+0300
```

4. Run the executable
```
$ ./bin/application
```

#### Run locally as Docker container

1. Clone the repo (outside GOPATH)

2. Open the terminal and go to the root folder
  
3. Create Docker image  
```
docker build -t skills-to-jobs .
```
4. Run as Docker container
```
docker-compose up
```

#### Tools

##### Run tests
```
$ make tests
```

##### Run code coverage tests
```
$ make cover
```

##### Run golint
```
$ make lint
```

##### Run gofmt to check formatting on all source files
```
$ make checkfmt
```

##### Run gofmt to fix formatting on all source files
```
$ make fixfmt
```

##### Cleanup everything
```
$ make clean
```

##### Run help
```
$ make help
```
##### Generate Swagger docs
To run this command, you will need to install [swagger-cli](https://github.com/APIDevTools/swagger-cli)
```
$ make oapi-gen-docs
```

##### Generate models from Swagger docs
To run this command, you will need to install [oapi-codegen](https://github.com/deepmap/oapi-codegen)
```
$ make make oapi-gen-types
```

### Test Application APIs

Verify the service is running as calling the get version API.

#### Call get version API

curl -X GET -i https://api-dev.rokwire.illinois.edu/skills-to-jobs/api/version

Response
```
0.1.2
```

## Contributing
If you would like to contribute to this project, please be sure to read the [Contributing Guidelines](CONTRIBUTING.md), [Code of Conduct](CODE_OF_CONDUCT.md), and [Conventions](CONVENTIONS.md) before beginning.

### Secret Detection
This repository is configured with a [pre-commit](https://pre-commit.com/) hook that runs [Yelp's Detect Secrets](https://github.com/Yelp/detect-secrets). If you intend to contribute directly to this repository, you must install pre-commit on your local machine to ensure that no secrets are pushed accidentally.

```
# Install software 
$ git pull  # Pull in pre-commit configuration & baseline 
$ pip install pre-commit 
$ pre-commit install

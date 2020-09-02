# Forky

```
git clone git@github.com:flowup/forky.git myproject
go run gen.go -project myproject
```
or
```
go run gen.go -project myproject -git-origin git@github.com:flowup/myproject.git
```

## Setup infrastructure
Create GCP project myproject-dev

Create bucket myproject-tf-states

Create services account `terraform@myproject-dev.iam.gserviceaccount.com`


Add yourself role as service account token creator and service account admin

Modify `project_number` in `env/dev.tvars` by using `gcloud projects list | grep myproject`

```
gcloud config set project myproject-dev
cd infrastructure/terraform
terraform init
terraform plan -var-file=env/dev.tfvars -out dev.tfplan 
terraform apply dev.tfplan
cd ../..
gcloud builds submit --config infrastructure/firebase.cloudbuild.yaml . 
gcloud builds submit --config infrastructure/ko.cloudbuild.yaml .
```

In cloudbuild install GithubApp and create trigger

Included files: `services/user/**`

Cloud build config file: `services/user/build/ci/cloudbuild.yaml`

## Running for the first time
Install and setup Go environment: https://golang.org/doc/install
```
cd services/user
go run cmd/compose/main.go
```

## Add new RPC
1. Modify proto file
1. Regenerate SDKs inside `apis`

### Go service:
1. In service/internal/service/service.go implement new method

In case of new DB calls, implement method on Repository interface.


## Generate SDKs from proto files

Go to `apis` and follow README.md

## Repository Structure

```
.
├── apis - Proto files and generated sdk's
│   ├── 3rdparty - Third party proto files
│   │   ├── google - Common proto files like Empty message
│   │   └── validate - For validation request messages
│   ├── go-sdk - Generated proto message and gRPC services for Go
│   ├── proto - Actual proto files for project
│   ├── Makefile - Used for running gen commands
│   └── README.md
├── infrastructure
│   └── terraform - IaS
├── README.md
└── services
    └── user
        ├── build
        │   └── ci - Cloudbuild for automatic builds
        ├── cmd
        │   ├── compose - Run each service at once
        │   └── user - Run single service
        └── internal
            └── user
                ├── repository.go - DB Layer
                └── service.go - API Layer


```

## FAQ

**Q: I already had some resources enabled in my project (app engine, firebase, etc.), what now?**

A: Simply import them using `terraform import`, e.g. you can import the app_engine app by running `terraform import -var-file=env/dev.tfvars google_app_engine_application.app <your-project-id>`

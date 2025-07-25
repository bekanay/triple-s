# triple-s

A lightweight clone of Amazon S3, built in Go.  
It offers a RESTful API for managing buckets and objects — allowing file upload, retrieval, and deletion — all served via HTTP with XML-based responses.

## Objectives

- HTTP server design in Go
- RESTful API implementation
- Basic networking and routing
- Metadata and file storage using CSV

## Features

- Create, list, and delete buckets
- Upload, retrieve, and delete objects (files)
- Serve and store data in local file system
- Store metadata in CSV files
- Responds in **XML format** to match S3 specs
- Command-line options for port and storage path

## Endpoints Overview

### Bucket Management

| Method | Endpoint         | Description                  |
|--------|------------------|------------------------------|
| PUT    | `/{BucketName}`  | Create a new bucket          |
| GET    | `/`              | List all buckets             |
| DELETE | `/{BucketName}`  | Delete an existing bucket    |

### Object Management

| Method | Endpoint                   | Description              |
|--------|----------------------------|--------------------------|
| PUT    | `/{BucketName}/{ObjectKey}`| Upload or overwrite file |
| GET    | `/{BucketName}/{ObjectKey}`| Retrieve object content  |
| DELETE | `/{BucketName}/{ObjectKey}`| Delete an object         |

> All responses follow XML format similar to Amazon S3.

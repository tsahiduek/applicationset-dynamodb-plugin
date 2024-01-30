# applicationset-dynamodb-plugin

This project is experimental plugin to pull configurations from DynamoDB table for ArgoCD ApplicationSets

## Table of Contents

- [applicationset-dynamodb-plugin](#applicationset-dynamodb-plugin)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Prerequisites](#prerequisites)
  - [Getting Started](#getting-started)
    - [Building](#building)

## Introduction

This application is designed to serve as a web server that interacts with DynamoDB to fetch and return parameters from a specified table. It supports both command-line flags and HTTP POST requests to customize the behavior.

## Prerequisites

Make sure you have the following prerequisites installed:

- Go (version 1.16 or higher)
- AWS CLI (if using AWS services)

## Getting Started

Follow the steps below to get started with the application.

Deploy the plugin into your cluster
`kubectl apply -f applicationset-sample.yaml`

Make sure to assign the appropriate EKS Pod Identity permissions to the service account created by the sample. You can easily follow this [blog post](https://aws.amazon.com/blogs/aws/amazon-eks-pod-identity-simplifies-iam-permissions-for-applications-on-amazon-eks-clusters/).



### Building

```bash
make build
```

# Terraform provider for CloudAMQP

Setup your CloudAMQP cluster from Terraform

## Getting Started (As a Terraform User)

### Prerequisites
Golang, Dep

## Install

* Install golang: https://golang.org/dl/

* Install go dep: https://golang.github.io/dep/docs/installation.html

* Install terraform: https://learn.hashicorp.com/terraform/getting-started/install.html

* Create a CloudAMQP account if you haven't already:

    * Go to https://www.cloudamqp.com/

    * Click "Sign Up"

    * Sign in

    * Go to API access (https://customer.cloudamqp.com/apikeys) and create a key.
      (note that this is the API key for one of the two APIs CloudAMQP supports.
      See https://docs.cloudamqp.com/cloudamqp_api.html.  We will discuss the other
      later.)

### Install CloudAMQP Terraform Provider
```sh
cd $GOPATH/src/github.com
mkdir cloudamqp
cd cloudamqp
git clone https://github.com/cloudamqp/terraform-provider.git
cd terraform-provider
make depupdate
make init
```

Now the provider is installed in the terraform plugins folder and ready to be used.

### Example Usage: Deploying a First CloudAMQP RMQ server

(See the examples.tf file in the repo.  It has a bunny VPC example and a simple lemur example.)

```sh
cd terraform-provider  #This is the root of the repo where examples.tf lives.
terraform plan
```
When prompted paste in your CloudAMQP API key (created above).

This will give you output on stdout that tells you what would have been created:
* rmq_lemur

Next run--
```sh
terraform apply
```

Again, paste in your API key.  This should create an actual CloudAMQP instance.

## Example

```hcl
provider "cloudamqp" {}

resource "cloudamqp_instance" "rmq_url"{
  name = "rmq_url"
  plan = "lemur"
  nodes = 1
  region = "amazon-web-services::us-east-1"
  rmq_version = "3.6.16"
  vpc_subnet = "10.201.0.0/24"
  tags = ["production", "project-name"]
}

output "rmq_url" {
  value = "${cloudamqp_instance.rmq_url.url}"
}

resource "cloudamqp_notification" "recipient_01" {
  instance_id = "${cloudamqp_instance.rmq_url.id}"
  type = "email"
  value = "alarm@example.com"
}

resource "cloudamqp_alarm" "alarm_01" {
  instance_id = "${cloudamqp_instance.rmq_url.id}"
  type = "cpu"
  value_threshold = 90
  time_threshold = 600
  notifications = ["${cloudamqp_notification.recipient_01.id}"]
}
```

## Import

Import existing infrastructure into state and bring the resource under Terraform management. Information about the resource will be added to the terraform.state file. Then add manually the given information to the .tf file. Once this is done, run terraform plan to see that the resource is under Terraform management. From here it's possible to add more resources such as alarm.

### Instance:

Import cloudamqp instance and bring it under Terraform management. First declare an empty instance resource

```resource "cloudamqp_instance"."rmq_url" {}´´´

Generic form of terraform import command
```terraform import {resource_type}.{resource_name} {resource_id}´´´

Example of terraform import command
```terraform import cloudamqp_instance.rmq_url 80´´´

### Resources depending on an instance:

All resources depending on the instance resource also need the instance id when using Terraform import. This is becuse the API call made, need to know what instance the resource depends on. Resource id and instance id is seperated with ",".

Resource affected by this is:
- cloudamqp_notification
- cloudmaqp_alarm

Generic form of terraform import command
```terraform import {resource_type}.{resource_name} {resource_id},{instance_id}´´´

Example of terraform import command
```terraform import cloudamqp_notification.recipient_01 10,80´´´
```terraform import cloudamqp_alarm.alarm_01 65,80´´´
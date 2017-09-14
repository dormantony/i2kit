# i2kit
i2kit is an immutable infrastructure (i2) deployment tool. It transforms k8 deployments in minimal linux distributions using linuxkit, and uses cloud provider technology to support networking, load balancing and service discovery.

The selling point is to have the goodness of docker for local development, ci and distribution of content, but keeping the robustness and performance of classic cloud vendor technologies. i2kit does not require a central service, eliminating the complexity and the abstraction layer of a cluster management tool. Security is also enhanced for two reasons: KVN provides true isolation between containers, and linuxkit reduce the attacking surface by reducing the amount of installed OS dependencies to the minimum.

# Implementation Details

The first prototype focuses on AWS, using VPC for networking, ELBs for exposing k8 deployments (aká a set of pods), Route53 CNAMES for k8 services and deployment endpoints and EBS for persistency. In other words, a k8 deployment is transformed into a linuxkit AMI, an auto scalability group with the desired number of instances, and a ELB configured for the ports defined in the k8 deployment. Also, a Route53 CNAME for `deployment-name.i2kit.com` is created that resolves to the deployment ELB.
i2kit also supports the deployment of k8 services by creating a CNAME that resolves to the ELBs of the deployment matching the the k8 service selector. In order to find these ELBs, i2kit uses AWS tags, tagging every k8 deployment with its labels.

# Getting Started

Make sure you have the `linuxkit` tool installed:

```
go get -u github.com/linuxkit/linuxkit/src/cmd/linuxkit
```

and the `aws-cli` configured with your credentials.

Now, build the `i2kit` binary:

```
go build -o /usr/local/bin/i2kit
```

- `I2KIT_DOCKER_CONFIG`: base64 encoded value for docker credentials.
- `I2KIT_CREDENTIALS`: pointing to AWS file credentials.
- `I2KIT_REGION`: AWS region where the deployment takes place.
- `I2KIT_HOSTED_ZONE`: a AWS hosted zone to create CNAMEs (for example "i2kit.com.").
- `I2KIT_SECURITY_GROUP`: AWS security group associated to the EC2 instances (default: will be auto-generated soon).
- `I2KIT_SUBNET`: AWS subnet where the EC2 instances are created (default: will be auto-generated soon).
- `I2KIT_KEYPAIR`: AWS keypair associated to the EC2 instances (default: will be auto-generated soon).
- `I2KIT_INSTANCE_TYPE`: AWS instance type used to create EC2 instances (default: `t2.micro`)

Extra configuration:

- S3 bucket in AWS named `linuxkit` // This will be autogenerated soon
- AWS hosted zone to create `CNAMEs`.

Execute commands by running:

```
i2kit deploy -f i2kit.yml
i2kit destroy -f i2kit.yml
```

where `i2kit.yml` is the path to your i2kit Manifest file.

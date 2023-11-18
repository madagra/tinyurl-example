import pulumi_aws as aws
from pulumi_aws import acm, route53
from pulumi import ResourceOptions

from config import APP_NAME, APP_PORT, DOMAIN, HOSTED_ZONE


# Read back the default VPC and public subnets
default_vpc = aws.ec2.get_vpc(default=True)
default_vpc_subnets = aws.ec2.get_subnets(
    filters=[
        aws.ec2.GetSubnetsFilterArgs(
            name="vpc-id",
            values=[default_vpc.id],
        ),
    ],
)

# Create a SecurityGroup that permits HTTP ingress and unrestricted egress.
app_security_group = aws.ec2.SecurityGroup(
    f"{APP_NAME}-secgrp",
    vpc_id=default_vpc.id,
    description="Enable HTTP(S) access",
    ingress=[
        aws.ec2.SecurityGroupIngressArgs(
            protocol="tcp",
            from_port=80,
            to_port=80,
            cidr_blocks=["0.0.0.0/0"],
        ),
        aws.ec2.SecurityGroupIngressArgs(
            protocol="tcp",
            from_port=443,
            to_port=443,
            cidr_blocks=["0.0.0.0/0"],
        ),
        aws.ec2.SecurityGroupIngressArgs(
            protocol="tcp",
            from_port=APP_PORT,
            to_port=APP_PORT,
            cidr_blocks=["0.0.0.0/0"],
        ),
    ],
    egress=[
        aws.ec2.SecurityGroupEgressArgs(
            protocol="-1",
            from_port=0,
            to_port=0,
            cidr_blocks=["0.0.0.0/0"],
        )
    ],
)

# Request a certificate for our domain
cert = acm.Certificate(
    f"{APP_NAME}-cert",
    domain_name=DOMAIN,
    validation_method="DNS",
)

# Route53 record for ACM valdiation
cert_validation_record = route53.Record(
    f"{APP_NAME}-certValidationRecord",
    allow_overwrite=True,
    name=cert.domain_validation_options[0]["resource_record_name"],
    records=[cert.domain_validation_options[0]["resource_record_value"]],
    ttl=60,
    type=cert.domain_validation_options[0]["resource_record_type"],
    zone_id=HOSTED_ZONE,
)

# Validate the certificate by applying the DNS record
cert_validation = acm.CertificateValidation(
    f"{APP_NAME}-certValidation",
    certificate_arn=cert.arn,
    validation_record_fqdns=[cert_validation_record.fqdn],
)

# Create a load balancer to listen for HTTP(S) traffic
load_balancer = aws.lb.LoadBalancer(
    f"{APP_NAME}-lb",
    security_groups=[app_security_group.id],
    subnets=default_vpc_subnets.ids,
)

lb_target_group = aws.lb.TargetGroup(
    f"{APP_NAME}-tg",
    port=APP_PORT,
    protocol="HTTP",
    target_type="ip",
    vpc_id=default_vpc.id,
)

# HTTPs listener with certificate
https_list = aws.lb.Listener(
    f"{APP_NAME}-https-listener",
    load_balancer_arn=load_balancer.arn,
    port=443,
    certificate_arn=cert.arn,
    protocol="HTTPS",
    ssl_policy="ELBSecurityPolicy-2016-08",
    default_actions=[
        aws.lb.ListenerDefaultActionArgs(
            type="forward",
            target_group_arn=lb_target_group.arn,
        )
    ],
    opts=ResourceOptions(depends_on=[cert, cert_validation, cert_validation_record]),
)

# HTTP listener with redirection
http_list = aws.lb.Listener(
    f"{APP_NAME}-http-listener",
    load_balancer_arn=load_balancer.arn,
    port=80,
    protocol="HTTP",
    default_actions=[
        aws.lb.ListenerDefaultActionArgs(
            type="redirect",
            redirect=aws.lb.ListenerDefaultActionRedirectArgs(
                port="443",
                protocol="HTTPS",
                status_code="HTTP_301",
            ),
        )
    ],
)

# DNS record to point from the domain name to the load balancer
dns_record = route53.Record(
    f"{APP_NAME}-dnsRecord",
    type="A",
    name=DOMAIN,
    aliases=[
        {
            "name": load_balancer.dns_name,
            "zone_id": load_balancer.zone_id,
            "evaluate_target_health": True,
        },
    ],
    zone_id=HOSTED_ZONE,
)

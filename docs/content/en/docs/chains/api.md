---
title: Api
description: All about api chains.
weight: 50
---
{{% alert title="Warning" color="warning" %}}
The api chain has the potential to be susceptible to Server-Side Request Forgery (SSRF) attacks if not used carefully and securely. SSRF allows an attacker to manipulate the server into making unintended and unauthorized requests to internal or external resources, which can lead to potential security breaches and unauthorized access to sensitive information.

To mitigate the risks associated with SSRF attacks, it is strongly advised to use the VerifyURL hook diligently. The VerifyURL hook should be implemented to validate and ensure that the generated URLs are restricted to authorized and safe resources only. By doing so, unauthorized access to sensitive resources can be prevented, and the application's security can be significantly enhanced.

It is the responsibility of developers and administrators to ensure the secure usage of the API chain. We strongly recommend thorough testing, security reviews, and adherence to secure coding practices to protect against potential security threats, including SSRF and other vulnerabilities.

See an example below.
{{% /alert %}}

{{< ghcode src="https://raw.githubusercontent.com/hupe1980/golc/main/examples/chains/api/main.go" >}}

Output:
```text
The current temperature in Munich, Germany is 55.5°F.
```
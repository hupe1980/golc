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

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/model/chatmodel"
)

func main() {
	openai, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *chatmodel.OpenAIOptions) {
		o.Temperature = 0
	})
	if err != nil {
		log.Fatal(err)
	}

	api, err := chain.NewAPI(openai, apiDoc, func(o *chain.APIOptions) {
		o.VerifyURL = func(url string) bool {
			return strings.HasPrefix(url, "https://api.open-meteo.com/")
		}
	})

	answer, err := golc.SimpleCall(context.Background(), api, "What is the weather like right now in Munich, Germany in degrees Fahrenheit?")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(answer)
}

const apiDoc = `BASE URL: https://api.open-meteo.com/

API Documentation
The API endpoint /v1/forecast accepts a geographical coordinate, a list of weather variables and responds with a JSON hourly weather forecast for 7 days. Time always starts at 0:00 today and contains 168 hours. All URL parameters are listed below:

Parameter	Format	Required	Default	Description
latitude, longitude	Floating point	Yes		Geographical WGS84 coordinate of the location
hourly	String array	No		A list of weather variables which should be returned. Values can be comma separated, or multiple &hourly= parameter in the URL can be used.
daily	String array	No		A list of daily weather variable aggregations which should be returned. Values can be comma separated, or multiple &daily= parameter in the URL can be used. If daily weather variables are specified, parameter timezone is required.
current_weather	Bool	No	false	Include current weather conditions in the JSON output.
temperature_unit	String	No	celsius	If fahrenheit is set, all temperature values are converted to Fahrenheit.
windspeed_unit	String	No	kmh	Other wind speed speed units: ms, mph and kn
precipitation_unit	String	No	mm	Other precipitation amount units: inch
timeformat	String	No	iso8601	If format unixtime is selected, all time values are returned in UNIX epoch time in seconds. Please note that all timestamp are in GMT+0! For daily values with unix timestamps, please apply utc_offset_seconds again to get the correct date.
timezone	String	No	GMT	If timezone is set, all timestamps are returned as local-time and data is returned starting at 00:00 local-time. Any time zone name from the time zone database is supported. If auto is set as a time zone, the coordinates will be automatically resolved to the local time zone.
past_days	Integer (0-2)	No	0	If past_days is set, yesterday or the day before yesterday data are also returned.
start_date
end_date	String (yyyy-mm-dd)	No		The time interval to get weather data. A day must be specified as an ISO8601 date (e.g. 2022-06-30).
models	String array	No	auto	Manually select one or more weather models. Per default, the best suitable weather models will be combined.

Hourly Parameter Definition
The parameter &hourly= accepts the following values. Most weather variables are given as an instantaneous value for the indicated hour. Some variables like precipitation are calculated from the preceding hour as an average or sum.

Variable	Valid time	Unit	Description
temperature_2m	Instant	°C (°F)	Air temperature at 2 meters above ground
snowfall	Preceding hour sum	cm (inch)	Snowfall amount of the preceding hour in centimeters. For the water equivalent in millimeter, divide by 7. E.g. 7 cm snow = 10 mm precipitation water equivalent
rain	Preceding hour sum	mm (inch)	Rain from large scale weather systems of the preceding hour in millimeter
showers	Preceding hour sum	mm (inch)	Showers from convective precipitation in millimeters from the preceding hour
weathercode	Instant	WMO code	Weather condition as a numeric code. Follow WMO weather interpretation codes. See table below for details.
snow_depth	Instant	meters	Snow depth on the ground
freezinglevel_height	Instant	meters	Altitude above sea level of the 0°C level
visibility	Instant	meters	Viewing distance in meters. Influenced by low clouds, humidity and aerosols. Maximum visibility is approximately 24 km.`
```
Output:
```text
The current temperature in Munich, Germany is 55.5°F.
```
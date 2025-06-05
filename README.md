# My Go Application

This is a simple Go application that starts a http server on 8080, and allows you to hit a single endpoint
gathering the forecast based on latitude and longitude. This was a quick, hour + change project.

## Project Requirements
Write an HTTP server that serves the forecasted weather. Your server should expose
an endpoint that:
1.     Accepts latitude and longitude coordinates
2.     Returns the short forecast for that area for Today (“Partly Cloudy” etc)
3.     Returns a characterization of whether the temperature is “hot”, “cold”, or
“moderate” (use your discretion on mapping temperatures to each type)
4.     Use the National Weather Service API Web Service as a data source.
The purpose of this exercise is to provide a sample of your work that we can discuss
together in the Technical Interview.
•         We respect your time. Spend as long as you need, but we intend it to take around
an hour.
•         We do not expect a production-ready service, but you might want to comment on
your shortcuts.
•         The submitted project should build and have brief instructions so we can verify
that it works.
•         The Coding Project should be written in the language for the job you’re applying
for. (golang)


## Project Structure

```
weather app
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── handlers
│   │   └── handler.go   # HTTP request handlers
│   └── services
│       └── service.go   # Business logic services
├── pkg
│   └── utils
│       └── utils.go     # Utility functions
├── go.mod               # Module dependencies
└── go.sum               # Module checksums
```

## Getting Started

To run the application, follow these steps:

1. Clone the repository:
   ```
   git clone <repository-url>
   cd weather-app
   ```

2. Install the dependencies:
   ```
   make tidy
   ```

3. Run the application:
   ```
   make build
   make run
   ```

## API Endpoints

- **GET /api/v1/forecast**: The endpoint to get the weather forecast.
It requires a lat and lon to be provided.
example ```
curl http://localhost:8080/api/v1/forecast\?lon\=-97.089\&lat\=39.745
{"latitude":39.745,"longitude":-97.089,"name":"Tonight","short_forecast":"Mostly Cloudy","temp_vibe":"perfect"}%
```

## Contributing

Feel free to submit issues or pull requests for improvements and bug fixes. 

## License

This project is licensed under the MIT License.
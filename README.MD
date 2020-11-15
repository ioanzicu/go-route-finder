# Route Finder

A REST API that takes the input parameters for home address location and destination(s) in format of latitude and longitude. The result is a sorted list of routes in JSON format.

### Prerequisites:

- [Golang](https://golang.org/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop) (if you want to run the application in a docker container)
- [`jq`(optional)](https://stedolan.github.io/jq/) tool for formatted json output

## Clone the repository on your computer

```bash
git clone https://github.com/ioanzicu/route-search-go.git
```

## Install dependencies:

```bash
cd go-route-finder/src
go get
go build
```

## Run the application

Make sure you are in the `src` folder

```bash
go run main.go
```

## Run the Application in the Docker Container

- **Navigate to the `/src` folder**

```bash
cd src
```

- **Build the docker image**

```bash
docker build -t go-route-finder .
```

Instead of `go-route-finder`, you can use your custom name for the docker image.

For Linux environment it can be require to use `sudo`.

```bash
sudo docker build -t go-route-finder .
```

- **Run the docker container**

```bash
sudo docker run -p [host port]:[exposed container port] go-route-finder
```

Example:

```bash
sudo docker run -p 8000:8000 go-route-finder
```

To run the container in the detached mode use `-d` flag

```bash
sudo docker run -d -p 8000:8000 go-route-finder
```

You can change the exposed container port number, but if you do so, do not
forget to change it in the `main.go` file too.

# Consume API with Curl

[`jq`](https://stedolan.github.io/jq/) command for formatted json output

```bash
curl -v "localhost:3000/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219" | jq
```

- `src=13.388860,52.517037`- represents the home address of the customer.

  **Note**: the `src` parameter should contain the longitude and latitude in the float format. Only one `src` parameter is allowed.

- `dst=13.428555,52.523219` - represents the destination location or pickup point.

  There can be more than one `dst` parameter provided separated by `&` symbol. Example: `dst=13.428555,52.523219&dst=13.342435,54.546675`.

The output will be a json object with fields: `source` and `routes`.
`routes` field hold a sorted by duration and/or distance list of destinations.

Example Output:

```bash
{
   "source":"13.38886,52.517037",
   "routes":[
      {
         "destination":"13.397634,52.529407",
         "duration":251.5,
         "distance":1884.8
      },
      {
         "destination":"13.428555,52.523219",
         "duration":394.2,
         "distance":3841.7
      }
   ]
}
```

## Unit Tests

In order to run Unit Test navigate into the controller folder `cd src/controller` and
run the command `go test -v`, flag `-v` for verbose output.

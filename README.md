## How to run program ?
Firstly, you must fill up environment variables and create a ``.env`` file in parent directory.

``
cp .env.template .env
``

``GOOGLE_MAPS_API_KEY`` must be present in ``.env`` file otherwise program will not run.

You must be in parent directory, then run command:

``
go run cmd/*.go
``

## How to generate gRPC GoLang client/server code from Protobuf ?

If you are in parent directory run:
```
cd http/grpc/proto

protoc --go_out=../proto-gen --go_opt=paths=source_relative \
--go-grpc_out=../proto-gen --go-grpc_opt=paths=source_relative \
*.proto
```

## What functionalities are completed?
- GetAlbum - it will return the album content, must provide ``album_id`` for a request,
CSV file's name represents ``album_id``. For example, CSV file name is ``1.csv`` then the ``album_id=1``.
This API method will return the parsed CSV data with the generated story titles for each
entry of CSV file. It has ``filter`` and ``sort`` fields in request which will help to return the list of photos
based on the given ``filter`` and ``sort``. 

- ListAlbums - it return the list of albums with their associated photos. It has the same functionallities
as above GetAlbum.



## What can be improved ?
* Currently the unit and integration tests are missing from the project, due to the time constraint for
the purpose of this project, the tests are not done, in real projects, the use of
BDD test with the help of Cucumber framework would be useful to check the behaviour of
the business logic, as it allows to make a client gRPC connection to the server itself.
If we implement BDD test we should create Interface for the API methods for third party providers so we can
mock this API calls during gRPC calls.


* I could not find a FREE API provider which would expose the historical data to fetch weather.
This would have done the story title of the photo interesting. Such as, user can select the photos
which are taken in sunny location, so a set of Sunny photos will be choosen and printed.



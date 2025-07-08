# blockchain-explorerCRD

## **Blockchain Explorer** is an ETL Pipeline that allows You to host your own cointanerized blockchain data blocks.

***Instructions:***
- (Needs Google Api json certificate authentication in the root directory).
- Once you have the binary created (it is already here but feel free to ***go build*** it again) you can run it with ***./blockchain-explorerCRD*** please add the flag ***--migrate*** in case the database needs migration.
- Here are cURL samples to test it in your machine:
  - curl -X POST http://localhost:8080/blocks \
-H 'Content-Type: application/json' \
-d '{"start":"2022-12-02T13:30:10.000Z","end":"2022-12-02T15:30:10.000Z"}'
  - curl -X GET http://localhost:8080/blocks
  - curl -X GET http://localhost:8080/blocks?nonce=NONCE_OF_BLOCK
  - curl -X DELETE http://localhost:8080/blocks?nonce=NONCE_OF_BLOCK
- Docker setup:
  - Build the image by using the Dockerfile in the current directory: *docker build -t blockchain-api-image .*
  - Run the container using the image built above: *docker run -it --rm --name blockchain-api-container -p 8080:8080 blockchain-api-image*
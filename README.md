# blockchain-explorerCRD

## **Blockchain Custom Explorer** can be used as it is, and provides a way to build on top of it, interesting applications.

***Instructions:***
- Once you have the binary created (it is already here but feel free to ***go build*** it again) you can run it with ***./blockchain-explorerCRD*** please add the flag ***--migrate*** in case the database needs migration.
- Here are cURL samples to test it in your machine:
  - curl -X POST http://localhost:8080/blocks \
-H 'Content-Type: application/json' \
-d '{"start":"2022-12-02T13:30:10.000Z","end":"2022-12-02T15:30:10.000Z"}'
  - curl -X GET http://localhost:8080/blocks
  - curl -X GET http://localhost:8080/blocks?nonce=NONCE_OF_BLOCK
  - curl -X DELETE http://localhost:8080/blocks?nonce=NONCE_OF_BLOCK
# blockchain-explorerCRD

## **Blockchain Custom Explorer** can be used as it is, and provides a way to build on top of it, interesting applications.

*Instructions:*
- Once you have the binary created (it is already here but feel free to *go build* it again) you can run it with *./blockchain-explorerCRD* please add the flag *--migrate* in case the database needs migration.
- In this item there will be cURL samples to test it on your machine
  - curl -X POST http://localhost:8080/blocks
  - curl -X GET http://localhost:8080/blocks
  - curl -X GET http://localhost:8080/blocks?nonce=FIELD_NONCE_OF_BLOCK_IN_DB
  - curl -X DELETE http://localhost:8080/blocks?nonce=FIELD_NONCE_OF_BLOCK_IN_DB
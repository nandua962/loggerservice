
# Partner API Documentation

## Introduction

The Partner  service provides functionality for managing partners within a system. This service allows  to create, retrieve, update, and delete partners efficiently through API.

## Purpose

The Partner API service offers the following key functionalities:

- **Create partner**: Allows to create new partner .
- **Get all partners**: Retrieves all partners  based on certain criteria such as filters .
- **Update partner**: Enables to update the partner details of existing partner .
- **Delete partner**: Allows to delete a partner.
- **Get partner By Id**: Retrieves the partner details of particular partner id .
- **Get partner Oauth credential**: Retrieves partner's oauth credential data for a specific provider .
- **Get terms and conditions**: Retrieves terms and conditions of a partner .
- **Update terms and conditions**:Enables to Updated terms and conditions of a partner.
- **Delete Arist role**: Allows to delete artist role.
- **Delete partner Genre**: Allows to delete genre  .
- **Update partner stores**: Enables to update the partner store .
- **Get partner stores**:retrieves all partner store.
- **Update partner status**: Enables to update the partner status .
- **Get partner payment gateways**:Retrieves the partner payment gateways.


## How to Use

For More details about endpoints and its payload 
refer :  https://docs.google.com/document/d/1Sbx0mDyh9DK151P1fBqVV5rTpfI4ZrCigznSk959xzE/edit

### 1. Create partner

To create a new partner, send a POST request to the `/api/v1/partners` endpoint with the required headers and a JSON payload 

### 2. Get all partners

To get all partners , send a GET request to the `/api/v1/partners` endpoint with optional query parameters such as page, limit, sort, order, name, and status.

### 3. Update a partner
 To update partner details , send a PATCH request to the `/api/v1/partners/:partner_id` endpoint with a valid partner_id

### 4. Delete a partner
 To delete a partner  , send a DELETE request to the `/api/v1/partners/:partner_id` endpoint with a valid partner_id

### 5. Get partner By Id
 To get a partner details  by id  , send a GET request to the `/api/v1/partners/:partner_id` endpoint with a valid partner_id

### 6. Get partner Oauth credential
 To get a partner oauth credentials , send a GET request to the `/api/v1/partners/:partner_id/oauth-credentials` endpoint with a valid partner_id

### 7. Get terms and conditions
 To get a partner terms and conditions , send a GET request to the `/api/v1/partners/:partner_id/terms-and-conditions` endpoint with a valid partner_id

### 8. Update terms and conditions
 To update a partner terms and conditions , send a PATCH request to the `/api/v1/partners/:partner_id/terms-and-conditions` endpoint with a valid partner_id

### 9. Delete artist role
 To delete a artist role , send a DELETE request to the `/api/v1/partners/:partner_id/artist-role/role_id` endpoint with a valid partner_id and role_id


### 10. Delete Partner genre
 To delete a genre , send a DELETE request to the `/api/v1/partners/:partner_id/genres/genre_id` endpoint with a valid partner_id and genre_id

### 11. Update Partner stores
 To update a partner store , send a PATCH request to the `/api/v1/partners/:partner_id/stores` endpoint with a valid partner_id

### 12. Get Partner stores
 To Get a partner store , send a GET request to the `/api/v1/partners/:partner_id/stores` endpoint with a valid partner_id 

### 13. Update Partner status
 To update a partner store , send a PATCH request to the `/api/v1/partners/:partner_id/statuss` endpoint with a valid partner_id

### 14. Get Partner payment gateways
 To get  partner payment gateways , send a GET request to the `/api/v1/partners/:partner_id/payment-gateways` endpoint with a valid partner_id
 


 ## env variables

- **PARTNER_DEBUG**: Indicates whether debugging mode is enabled (`true` or `false`).
- **PARTNER_PORT**: Specifies the port number for the service.
- **PARTNER_DB_USER**: Specifies the username for the database connection.
- **PARTNER_DB_PORT**: Specifies the port number for the database connection.
- **PARTNER_DB_PASSWORD**: Specifies the password for the database connection.
- **PARTNER_DB_DATABASE**: Specifies the name of the database to connect to.
- **PARTNER_ACCEPTED_VERSIONS**: Indicates the accepted versions of something (API versions).
- **PARTNER_DB_SCHEMA**: Specifies the schema to be used in the database.
- **PARTNER_DB_HOST**: Specifies the host address for the database connection.
- **PARTNER_CACHE_EXPIRATION**: Specifies the expiration time for cached data (in some unit, possibly days or hours).
- **PARTNER_DB_SSLMODE**: Specifies whether SSL mode is enabled for the database connection.
- **LOCALISATION_SERVICE_URL**: Specifies the URL for a localization service.
- **ENDPOINT_URL**: Specifies the URL for an endpoint.
- **LOGGER_SERVICE_URL**: Specifies the URL for a logger service.
- **LOGGER_SECRET**: Specifies a secret key for the logger service.
- **ERROR_HELP_LINK**: Specifies a link for accessing error help documentation.
- **ACTIVITY_LOG_SERVICE_URL**: Specifies the URL for an activity log service.
- **CLIENT_CREDENTIAL_ENCRYPTION_KEY**:Encryption and decryption key for client credential data
- **UTILITY_SERVICE_URL**: specifies the URL for fetching  data from utility service
- **SUBSCRIPTION_SERVICE_URL**:specifies the URL for fetching subcription data from subscription service
- **MEMBER_SERVICE_URL**:specifies the URL for fetching member data from member service
- **PARTNER_API_URL**:specifies the URL of PARTNER service
- **OAUTH_SERVICE_URL**:specifies the URL for fetching oauth provider data from oauth service
- **STORE_SERVICE_URL**:specifies the URL for fetching store data from store service
- **PARTNER_PAYMENT_GATEWAY_ENCRYPTION_KEY** :specifies the payment gateway encryption and decryption  key


# How to run

To run the Partner service, execute the command `go run main.go`. This will start the service, allowing it to perform its functionalities within the system.

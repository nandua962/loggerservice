# Utility

The Utility service is a comprehensive system that provides essential functionalities for managing various aspects like currency, language, country, and more. This service ensures efficient handling and retrieval of relevant data through API endpoints.

## 1. Country Management

The Country API service provides comprehensive tools for managing country and state data.

- **Get Countries** : Retrieve a list of countries based on provided query parameters.
- **Get States of Country**: Retrieve a list of states for a specific country.
- **Check Country Exists**: Check the existence of countries based on ISO parameters.
- **Get All Country Codes**: Retrieve all country codes.
- **Check State Exists**: Check the existence of states based on ISO and country code parameters.

## 2. Currency Management

The Currency API service handles the retrieval and management of currency data.

- **Get All Currencies**: Retrieves a list of currencies based on provided parameters.
- **Get Currency by ID**: Retrieves currency details based on the provided currency ID.
- **Get Currency by ISO Code**: Retrieves currency details based on the provided ISO code.

## 3. Language Management

The Language API service enables the management and retrieval of language data.

- **Get Languages**: Allows users to etrieve a list of all available languages with options for filtering, sorting, and pagination.
- **Check Language Code Existence**: Enables users to check if a specific language code exists in the system.

## 4. Theme Management

The Theme API service allows for the management of themes within the system.

- **Get Theme by ID** : Allows users to retrieve a specific theme using its unique identifier.

## 5. Role Management

The Role API service provides tools for managing user roles.

- **Create Role**: Allows users to create new roles with specified attributes.
- **Get Role by ID**: Retrieves a specific role by its ID.
- **Get all Roles**: Retrieve all roles or filter roles based on criteria such as page, limit, sort, order, and name.
- **Update Role**: To update the attributes of existing roles.
- **Delete Role**: To remove roles from the system.

## 6. Payment Gateway Management

The Payment Gateway API service manages payment gateway information.
 
- **Get Payment Gateway by ID**: Retrieve detailed information about a specific payment gateway by its ID.
- **Get Payment Gateway by ID**: Retrieve all payment gateway details or filter payment gateway based on criteria such as page, limit, sort, order, and name.

## 7. Lookup Management

The Lookup API service handles the retrieval of lookup values.

- **Get Lookup by lookup type name**: To retrieve all lookup values or filter them based on a specific lookup type name
- **Get Lookup Value by ID**: To retrieve a specific lookup value by its ID.


## 8. Genre Management

The Genre API service provides tools for managing music genres.

- **Create Genre**: Allows users to create new genres with specified attributes.
- **Get all Genres**: Retrieves all genres or filters genres based on certain criteria such as page, limit, sort, order, and name.
- **Get Genre by ID**: Retrieve detailed information about a specific genre by genre ID.
- **Update Genre**: To update the attributes of existing genres.
- **Delete Genre**: To remove genres from the system.

## env variables

- **UTILITY_DEBUG**: Indicates whether debugging mode is enabled (`true` or `false`).
- **UTILITY_PORT**: Specifies the port number for the service.
- **UTILITY_DB_USER**: Specifies the username for the database connection.
- **UTILITY_DB_PORT**: Specifies the port number for the database connection.
- **UTILITY_DB_PASSWORD**: Specifies the password for the database connection.
- **UTILITY_DB_DATABASE**: Specifies the name of the database to connect to.
- **UTILITY_ACCEPTED_VERSIONS**: Indicates the accepted versions of something (API versions).
- **UTILITY_DB_SCHEMA**: Specifies the schema to be used in the database.
- **UTILITY_DB_HOST**: Specifies the host address for the database connection.
- **UTILITY_LOCALISATION_SERVICE_URL**: Specifies the URL for a localization service.
- **UTILITY_LOGGER_SERVICE_URL**: Specifies the URL for a logger service.
- **UTILITY_LOGGER_SECRET**: Specifies a secret key for the logger service.
- **UTILITY_ERROR_HELP_LINK**: Specifies a link for accessing error help documentation.
- **UTILITY_ACTIVITY_LOG_SERVICE_URL**: Specifies the URL for an activity log service.

# How to run

To run the Utility service, execute the command `go run main.go`. This will start the service, allowing it to perform its functionalities within the system.

# For further information:

https://docs.google.com/document/d/10TrYLkynTNG_-Axu3twyfA0OTjQ_a_ou/edit

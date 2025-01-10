# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/proto/auth/auth.proto](#api_proto_auth_auth-proto)
    - [AuthRequest](#auth-AuthRequest)
    - [AuthResponse](#auth-AuthResponse)
    - [VerifyTokenRequest](#auth-VerifyTokenRequest)
    - [VerifyTokenResponse](#auth-VerifyTokenResponse)
  
    - [AuthService](#auth-AuthService)
  
- [api/proto/graph/graph.proto](#api_proto_graph_graph-proto)
    - [ConnectionTypeRequest](#graph-ConnectionTypeRequest)
    - [ConnectionTypesList](#graph-ConnectionTypesList)
    - [Empty](#graph-Empty)
    - [EntitiesList](#graph-EntitiesList)
    - [EntityRequest](#graph-EntityRequest)
    - [PropertyTypeRequest](#graph-PropertyTypeRequest)
    - [PropertyTypesList](#graph-PropertyTypesList)
    - [SearchRequest](#graph-SearchRequest)
    - [UserData](#graph-UserData)
    - [UserRequest](#graph-UserRequest)
    - [UsersConnectionType](#graph-UsersConnectionType)
    - [UsersEntity](#graph-UsersEntity)
    - [UsersPropertyType](#graph-UsersPropertyType)
  
    - [GraphService](#graph-GraphService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_proto_auth_auth-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/proto/auth/auth.proto



<a name="auth-AuthRequest"></a>

### AuthRequest
AuthRequest represents the authentication request for both registration and login.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| email | [string](#string) |  | [REQUIRED] [MAX LEN 255] Email must be a valid email address format (e.g., &#34;user@example.com&#34;). Example: &#34;john.doe@company.com&#34; |
| password | [string](#string) |  | [REQUIRED] [MIN LEN 8] [MAX LEN 32] Password is hashed on the server before storing. Example: &#34;MySecurePass123!&#34; |






<a name="auth-AuthResponse"></a>

### AuthResponse
AuthResponse represents the server&#39;s response to successful authentication.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  | Unique identifier for the user. Format: UUID v4. Example: &#34;123e4567-e89b-12d3-a456-426614174000&#34; |
| email | [string](#string) |  | Email address associated with the authenticated user. Email must be a valid email address format (e.g., &#34;user@example.com&#34;). Example: &#34;john.doe@company.com&#34; |
| token | [string](#string) |  | JWT token for subsequent authenticated requests. Format: JWT string (header.payload.signature). Must be included in subsequent requests as &#34;authorization&#34; metadata. Valid for: 24 hours by default. Example: &#34;eyJhbGciOiJIUzI1NiIs...&#34; |






<a name="auth-VerifyTokenRequest"></a>

### VerifyTokenRequest
VerifyTokenRequest represents a token verification request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  | [REQUIRED] JWT token to verify. Must be a valid JWT token previously issued by the auth service. Example: &#34;eyJhbGciOiJIUzI1NiIs...&#34; |






<a name="auth-VerifyTokenResponse"></a>

### VerifyTokenResponse
VerifyTokenResponse contains user information if token is valid.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  | User ID associated with the token. Format: UUID v4. Example: &#34;123e4567-e89b-12d3-a456-426614174000&#34; |
| email | [string](#string) |  | Email address associated with the user. Format: valid email address. Example: &#34;john.doe@company.com&#34; |





 

 

 


<a name="auth-AuthService"></a>

### AuthService
AuthService provides authentication and authorization functionality.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Register | [AuthRequest](#auth-AuthRequest) | [AuthResponse](#auth-AuthResponse) | Register creates a new user account. Errors: (INVALID_ARGUMENT): If email format is invalid or password doesn&#39;t meet requirements (ALREADY_EXISTS): If the email is already registered (INTERNAL): For server-side errors |
| Login | [AuthRequest](#auth-AuthRequest) | [AuthResponse](#auth-AuthResponse) | Login authenticates an existing user. Errors: (INVALID_ARGUMENT): If email format is invalid (NOT_FOUND): If the email is not registered (UNAUTHENTICATED): If the password is incorrect (INTERNAL): For server-side errors |
| VerifyToken | [VerifyTokenRequest](#auth-VerifyTokenRequest) | [VerifyTokenResponse](#auth-VerifyTokenResponse) | VerifyToken validates a JWT token and returns associated user information. Errors: (INVALID_ARGUMENT): If token format is invalid (UNAUTHENTICATED): If token is expired or invalid (INTERNAL): For server-side errors |

 



<a name="api_proto_graph_graph-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/proto/graph/graph.proto



<a name="graph-ConnectionTypeRequest"></a>

### ConnectionTypeRequest
ConnectionTypeRequest represents a request to create a connection type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | [OPTIONAL] [MAX LEN 255] [FORMAT UUID v4] Unique identifier for the connection type, recieved by the FindConnectionTypes endpoint. You cannot create your own ID. If provided for CREATE operations, server will link the users version of the connection type to the shared connection type. Required for UPDATE operations. Example: &#34;123e4567-e89b-12d3-a456-426614174000&#34; |
| name | [string](#string) |  | [REQUIRED] [MAX LEN 255] Name for the users version of this connection type. Can be seen by other users. Example: &#34;Works at&#34; or &#34;Brother&#34; |
| definition | [string](#string) |  | [REQUIRED] [MAX LEN 4096] Description for the users version of this connection type. Should provide clear, comprehensive information about the connection type, can be seen by other users. Example: &#34;Represents a current employment relationship between a person and a company&#34; |






<a name="graph-ConnectionTypesList"></a>

### ConnectionTypesList
ConnectionTypesList represents a collection of connection types matching a search query.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| connection_types | [UsersConnectionType](#graph-UsersConnectionType) | repeated | List of connection types matching the search criteria. May be empty if no matches are found |






<a name="graph-Empty"></a>

### Empty
Empty message for requests/responses that don&#39;t need any data






<a name="graph-EntitiesList"></a>

### EntitiesList
EntitiesList represents a collection of entities matching a search query.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entities | [UsersEntity](#graph-UsersEntity) | repeated | List of entities matching the search criteria. May be empty if no matches are found |






<a name="graph-EntityRequest"></a>

### EntityRequest
EntityRequest represents a request to create or update an entity.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | [OPTIONAL] [MAX LEN 255] [FORMAT UUID v4] Unique identifier for the entity, recieved by the FindEntities endpoint. You cannot create your own ID. If provided for CREATE operations, server will link the users version of the entity to the shared entity. Required for UPDATE operations. Example: &#34;123e4567-e89b-12d3-a456-426614174000&#34; |
| name | [string](#string) |  | [REQUIRED] [MAX LEN 255] Name for the users version of this entity. Can be seen by other users. Example: &#34;John Doe&#34; or &#34;Company XYZ&#34;. |
| definition | [string](#string) |  | [REQUIRED] [MAX LEN 4096] Description for the users version of this entity. Should provide clear, comprehensive information about the entity, can be seen by other users. Example: &#34;Senior Software Engineer with 10 years of experience...&#34; |






<a name="graph-PropertyTypeRequest"></a>

### PropertyTypeRequest
PropertyTypeRequest represents a request to create a property type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | [OPTIONAL] [MAX LEN 255] [FORMAT UUID v4] Unique identifier for the property type, recieved by the FindPropertyTypes endpoint. You cannot create your own ID. If provided for CREATE operations, server will link the users version of the property type to the shared property type. Required for UPDATE operations. Example: &#34;123e4567-e89b-12d3-a456-426614174000&#34; |
| name | [string](#string) |  | [REQUIRED] [MAX LEN 255] Name for the users version of this property type Can be seen by other users Example: &#34;Salary&#34; or &#34;Starts on&#34; |
| definition | [string](#string) |  | [REQUIRED] [MAX LEN 4096] Description for the users version of this property type Should provide clear, comprehensive information about the property type, can be seen by other users Example: &#34;Annual gross salary in USD&#34; |
| value_type | [string](#string) |  | Data type for this property MUST be one of: &#34;string&#34;, &#34;int&#34;, &#34;float&#34;, &#34;boolean&#34; Example: &#34;float&#34; for salary |






<a name="graph-PropertyTypesList"></a>

### PropertyTypesList
PropertyTypesList represents a collection of property types matching a search query.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| property_types | [UsersPropertyType](#graph-UsersPropertyType) | repeated | List of property types matching the search criteria. May be empty if no matches are found |






<a name="graph-SearchRequest"></a>

### SearchRequest
SearchRequest represents a search query for finding entities, connection types, or property types.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | [REQUIRED] [MAX LEN 255] Name to search for (case-sensitive exact match). Example: &#34;Person&#34; or &#34;Vehicle&#34; |






<a name="graph-UserData"></a>

### UserData
UserData represents all graph data associated with a user.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entities | [UsersEntity](#graph-UsersEntity) | repeated | List of all entities created or linked by the user. May be empty for new users |
| connection_types | [UsersConnectionType](#graph-UsersConnectionType) | repeated | List of all connection types created or linked by the user. May be empty for new users |
| property_types | [UsersPropertyType](#graph-UsersPropertyType) | repeated | List of all property types created or linked by the user. May be empty for new users |






<a name="graph-UserRequest"></a>

### UserRequest
UserRequest represents a request to create a new user in the graph service. This endpoint can only be used by the auth service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | [REQUIRED] Unique identifier for the user Example: &#34;123e4567-e89b-12d3-a456-426614174000&#34; |






<a name="graph-UsersConnectionType"></a>

### UsersConnectionType
UsersConnectionType represents a user&#39;s version of a connection type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | This user&#39;s name for the connection type Example: &#34;Works at&#34; |
| definition | [string](#string) |  | This user&#39;s definition of the connection type Example: &#34;Represents a current employment relationship...&#34; |
| user_id | [string](#string) |  | ID of the user who created this version of the entity |
| connection_type_id | [string](#string) |  | ID of the underlying shared connection type |






<a name="graph-UsersEntity"></a>

### UsersEntity
UsersEntity represents a user&#39;s version of an entity.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | This user&#39;s name for the entity Example: &#34;John Doe&#34; |
| definition | [string](#string) |  | This user&#39;s definition of the entity Example: &#34;Senior Software Engineer in our team...&#34; |
| user_id | [string](#string) |  | ID of the user who created this version of the entity |
| entity_id | [string](#string) |  | ID of the underlying shared entity Format: UUID v4 |






<a name="graph-UsersPropertyType"></a>

### UsersPropertyType
UsersPropertyType represents a user&#39;s version of a property type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | This user&#39;s name for the property type Example: &#34;Salary&#34; |
| definition | [string](#string) |  | This user&#39;s definition of the property type Example: &#34;Annual gross salary in USD&#34; |
| property_type_id | [string](#string) |  | ID of the underlying shared property type |
| user_id | [string](#string) |  | ID of the user who created this version of the entity |
| value_type | [string](#string) |  | Data type for this property One of: &#34;string&#34;, &#34;int&#34;, &#34;float&#34;, &#34;boolean&#34; |





 

 

 


<a name="graph-GraphService"></a>

### GraphService
GraphService provides operations for managing graph-based knowledge representation.
All operations require authentication via JWT token in the &#34;authorization&#34; metadata.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateUser | [UserRequest](#graph-UserRequest) | [Empty](#graph-Empty) | CreateUser initializes a new user in the graph service. This endpoint is only accessible by the auth service. Errors: (INVALID_ARGUMENT): If user_id format is invalid (ALREADY_EXISTS): If user already exists (INTERNAL): For server-side errors |
| GetUserData | [Empty](#graph-Empty) | [UserData](#graph-UserData) | GetUserData retrieves all entities, connection types, and property types associated with the authenticated user. Errors: (UNAUTHENTICATED): If authentication is missing or invalid (INTERNAL): For server-side errors |
| CreateEntity | [EntityRequest](#graph-EntityRequest) | [UsersEntity](#graph-UsersEntity) | CreateEntity creates a new entity or links to an existing one if ID is provided. Errors: (INVALID_ARGUMENT): If name or definition exceed length limits (NOT_FOUND): If entity ID is provided but entity doesn&#39;t exist (UNAUTHENTICATED): If authentication is missing or invalid (INTERNAL): For server-side errors |
| UpdateEntity | [EntityRequest](#graph-EntityRequest) | [Empty](#graph-Empty) | UpdateEntity modifies an existing entity. Errors: (INVALID_ARGUMENT): If name or definition exceed length limits (NOT_FOUND): If entity doesn&#39;t exist (UNAUTHENTICATED): If authentication is missing or invalid (INTERNAL): For server-side errors |
| FindEntities | [SearchRequest](#graph-SearchRequest) | [EntitiesList](#graph-EntitiesList) | FindEntities searches for entities by exact name match. Errors: (INVALID_ARGUMENT): If name is empty or too long (UNAUTHENTICATED): If authentication is missing or invalid (INTERNAL): For server-side errors |
| CreateConnectionType | [ConnectionTypeRequest](#graph-ConnectionTypeRequest) | [UsersConnectionType](#graph-UsersConnectionType) | CreateConnectionType creates a new connection type or links to an existing one. Errors: (INVALID_ARGUMENT): If name or definition exceed length limits (NOT_FOUND): If entity ID is provided but entity doesn&#39;t exist (UNAUTHENTICATED): If authentication is missing or invalid (INTERNAL): For server-side errors |
| FindConnectionTypes | [SearchRequest](#graph-SearchRequest) | [ConnectionTypesList](#graph-ConnectionTypesList) | FindConnectionTypes searches for connection types by exact name match. Errors: (INVALID_ARGUMENT): If name is empty or too long (UNAUTHENTICATED): If authentication is missing or invalid (INTERNAL): For server-side errors |
| CreatePropertyType | [PropertyTypeRequest](#graph-PropertyTypeRequest) | [UsersPropertyType](#graph-UsersPropertyType) | CreatePropertyType creates a new property type or links to an existing one. Errors: (INVALID_ARGUMENT): If name or definition exceed length limits (NOT_FOUND): If entity ID is provided but entity doesn&#39;t exist (UNAUTHENTICATED): If authentication is missing or invalid (INTERNAL): For server-side errors Example request: |
| FindPropertyTypes | [SearchRequest](#graph-SearchRequest) | [PropertyTypesList](#graph-PropertyTypesList) | FindPropertyTypes searches for property types by exact name match. Errors: (INVALID_ARGUMENT): If name is empty or too long (UNAUTHENTICATED): If authentication is missing or invalid (INTERNAL): For server-side errors |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |


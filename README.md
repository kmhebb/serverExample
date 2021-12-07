# A microservice api implementation

This API is a microservice architecture for delivering a remote procedure call interface to a thin client front end. While this implementation is structured as a standalone server to be deployed on a service like heroku or linode, the server architecture is handler based and therefore deployable in a variety of environments, from kubernetes, heroku, lambda, or even hosted on local metal. Using bespoke handler options and type functions, the handler can support rendering responses in virtually any format - responses here are primarily application/json, but there are some hardwired html responses for some specific direct user interactions.

It provides an interface that can be easily monitored and will take itself offline if any required dependencies are found to be unavailable. It implements a logging package that can interface with many reporting platforms out there. It implements a customized error interface to deliver actionable error information to the client. It implements integration with sendgrid and slack for user communications. It can be configured to use caching for ephemeral data.

It is a fully authenticated and context driven implementation using a customizable bearer token authentication.

This implementation has a tailored data interface for a utility billing provider. While the specific database implementation here is transaction based postgres, the service interface is agnostic and transactionable.

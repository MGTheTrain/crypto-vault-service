
# Blob Storage Schema Design

## Decision:

Blobs will be stored in a Blob Storage using the schema /{ID}/{Name}, where:

- **ID:** A unique identifier for the entity (e.g., user ID, product ID).
- **Name:** The name or description of the specific file (e.g., profile_picture.png, invoice.pdf).

## Conclusion:

The `/{ID}/{Name}` schema provides a logical, scalable, and efficient way to organize blobs. It ensures clear grouping of related blobs, simplifies retrieval, and leverages a Blob Storage's hierarchical structure for better performance and management.
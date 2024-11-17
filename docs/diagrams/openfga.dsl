model
  schema 1.1
  
type user

type blob
  relations
    # Permissions related to file management
    define can_manage_all_blobs: admin
    define can_manage_own_blobs: owner
    define can_download_blobs_with_given_permission: grantee
    define can_view_blobs_with_given_permission: grantee

    # Cryptographic actions
    define create_own_cryptographic_keys: owner
    define encrypt_decrypt_own_files: owner
    define generate_signature_for_own_files: owner
    define verify_file_signature: owner or grantee  # Public key verification is possible for grantee
    
    # Access control for owners and grantees
    define can_grant_access_to_download_owned_blobs: owner
    define can_grant_access_to_view_owned_blobs: owner

    # Ownership and user roles
    define owner: [user]
    define grantee: [user]  # A user who has been granted permissions for an owner's blob
    define admin: [user]    # Admin can manage all blobs, including cryptographic actions

    # Additional clarifications
    # - Admin has full control over all blobs
    # - Owner controls access to their own blob, including granting permissions
    # - Grantee has permission to download or view blobs if granted by the owner

    # In order to visualize the relationship based model visit https://play.fga.dev/sandbox/?store=github and paste in this content

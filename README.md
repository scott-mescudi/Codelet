# Code snippets app - NexGPT

# **Features:**

# TODO
- fortify auth


- **Distinct Categories for Each Programming Language**:
    
    The platform will organize code snippets into different categories based on programming languages, allowing users to easily find and manage snippets relevant to a specific language.
    
- **Tagging System for Snippets**:
    
    Each code snippet can be assigned multiple tags such as "frontend," "backend," "database," or other relevant categories. This helps improve searchability and enables users to filter snippets based on their needs.
    
- **User Management System**:
    
    Users will have personalized accounts, enabling them to save, manage, and organize their own code snippets. This system will also support authentication and secure access.
    
- **Public and Private Code Snippets**:
    
    Users can choose whether to make their code snippets publicly available for others to view or keep them private for personal use. This ensures flexibility in sharing and privacy.
    
- **Explore Page**:
    
    A dedicated "Explore" section will showcase publicly shared code snippets, enabling users to browse through popular, trending, or recently added snippets from the community.
    
- **Comprehensive Metadata for Snippets**:
    
    Each snippet will include the following essential details:
    
    - **Code**: The actual code content of the snippet.
    - **Tags**: Relevant keywords or categories for easy organization and filtering.
    - **Title**: A descriptive title summarizing the purpose of the snippet.
    - **Description**: A brief explanation of what the code does, its purpose, or how it should be used.
    - **Date Added**: The date the snippet was initially created or uploaded.
    - **Last Updated**: The date the snippet was last modified.
    - Size: Each snippet must be max 3kb
- **Snippet Editing**:
    
    Users will have the ability to edit their existing snippets to update or refine the content, metadata, or tags as needed.
    

# **Frontend:**

- nextjs with typescript
- tailwind

# **Backend:**

- go + fiber for gateway
- go + grpc for services
- supabase
- docker
- memcached

# Misc:

- ZSTD compression for snippets
- ratelimit

# User DB schema

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  role VARCHAR(20) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

# Snippets DB schema

```sql
CREATE TABLE snippets (
  id SERIAL PRIMARY KEY,
  userid INT NOT NULL,
  language VARCHAR(50),
  title VARCHAR(255) NOT NULL,
  code TEXT,
  description TEXT,
  tags VARCHAR(50)[],
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
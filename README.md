# Codelet - An Open-Source Code Snippet App

![Home Screen](./docs/image2.png)
![Dashboard](./docs/image.png)

## All Your Code Snippets, Organized

Stop wasting time searching. Access all your code snippets instantly, organized in one place for easy use.

## Features

- **User System**: Sign up and manage your snippets with a personal account.
- **Syntax Highlighting**: Enjoy beautifully formatted code with syntax highlighting.
- **Categories**: Organize snippets by category for easy retrieval.
- **Dashboard**: A clean and intuitive dashboard to manage all your saved snippets.

## Tech Stack

- **Frontend**: Next.js, Tailwind CSS
- **Backend**: Golang
- **Database**: PostgreSQL
- **Containerization**: Docker
- **Authentication**: JWT-based authentication
- **Code Highlighting**: Prism.js

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/codelet.git
   ```
2. Navigate to the project folder:
   ```sh
   cd codelet
   ```
3. Start the Backend:
   ```sh
   docker compose up -d
   ```
4. Start the Frontend:
   ```sh
   cd web && npm install && npm run dev
   ```
5. Access the application at:
   ```
   http://localhost:3000
   ```


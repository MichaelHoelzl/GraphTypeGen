# 📦 GraphTypeGen

**GraphTypeGen** is a command-line tool that generates a fully-typed **TypeScript library** from your **GraphQL schema** files. It streamlines the process of creating type-safe query and mutation functions, making your codebase more maintainable, readable, and less error-prone.

---

## 🚀 Features

- 🔒 **Type-safe**: Automatically generates TypeScript interfaces from your GraphQL types.  
- ⚡ **Fast & Efficient**: Quickly produces ready-to-use TypeScript functions for queries and mutations.  
- 🧩 **Customizable**: Add custom headers and specify client configurations (e.g., Apollo Client).  
- 🚀 **Error Handling**: Optional built-in error handling for safer API calls.  
- 📝 **Clean Output**: Organized, readable, and production-ready TypeScript code.  

---

## 🛠️ How It Works

1. **Provide a GraphQL Schema:** Pass your `.graphql` or `.gql` schema file to the tool.  
2. **Configure Output:** Specify the output path, client name (e.g., `apolloClient`), and optional error handling preferences.  
3. **Generate Code:** The tool creates TypeScript interfaces for your GraphQL types and functions for your queries and mutations.  

---

## 📝 Example Usage

```bash
graph-type-gen \
  -schema="./schema.graphql" \
  -output="./generated/graphql.ts" \
  -header="import { gql } from '@apollo/client';" \
  -client="apolloClient" \
  -error="true"
```

#### ✅ This generates:
- **TypeScript interfaces** matching your GraphQL types.  
- **Query and Mutation functions** using your provided client.  
- **Automatic error handling** (if `-error` flag is provided).  

---

## 📂 Example Generated Code

```typescript
export interface User {
  id: string;
  name: string;
  email: string;
}

async function getUser(id: string): Promise<[User | null, any | null]> {
  try {
    const response = await apolloClient.query({
      query: gql(`
        query getUser($id: String) {
          getUser(id: $id) {
            id
            name
            email
          }
        }
      `),
      variables: { id },
    });
    return [response.data.getUser, null];
  } catch (error) {
    return [null, error];
  }
}
```

---

## 🧷 Command-Line Options

| Flag         | Description                                       | Required |
|--------------|---------------------------------------------------|----------|
| `-schema`    | Path to the GraphQL schema file                   | ✅       |
| `-output`    | Path to the generated TypeScript file             | ✅       |
| `-header`    | Header code to insert (e.g., imports)             | ✅       |
| `-client`    | Name of the GraphQL client (e.g., `apolloClient`) | ✅       |
| `-error`     | Return errors in function calls (optional)        | ❌       |

---

## 🧪 Development

1. Clone the repository:  
   ```bash
   git clone https://github.com/MichaelHoelzl/GraphTypeGen.git
   cd GraphTypeGen
   ```

2. Build and run the tool:  
   ```bash
   go build -o graph-type-gen
   ./graph-type-gen -schema="./schema.graphql" -output="./generated.ts" -header="import { gql } from '@apollo/client';" -client="apolloClient"
   ```

---

## 📄 License

This project is licensed under the **MIT License** – feel free to use and contribute!  


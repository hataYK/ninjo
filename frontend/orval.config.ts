import { defineConfig } from "orval";

export default defineConfig({
  ninjo: {
    input: {
      target: "../docs/openapi/openapi.yaml",
    },
    output: {
      target: "./src/api/generated/ninjo.ts",
      client: "react-query",
      mode: "single",
      override: {
        mutator: {
          path: "./src/lib/fetch.ts",
          name: "customFetch",
        },
        query: {
          useQuery: true,
          useMutation: true,
        },
      },
    },
  },
});

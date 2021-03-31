import { ApolloClient, InMemoryCache } from "@apollo/client";
import { createUploadLink } from "apollo-upload-client";
import fetch from "isomorphic-unfetch";
import { GRAPHQL_API_URL } from "@src/constants/api";


const apolloClient = new ApolloClient({
  ssrMode: true,
  link: createUploadLink({
    uri: `${GRAPHQL_API_URL}/query`,
    credentials: "include",
    fetch,
    headers: {
      "X-Requested-By": "eitan-web-client", // for CSRF validation
    },
  }),
  cache: new InMemoryCache(),
});

export default apolloClient;

import React from 'react';
import { AppProps } from 'next/app';
import Head from 'next/head';
import 'tailwindcss/tailwind.css';
import { ApolloProvider } from "@apollo/client"
import Header from '../components/Header';
import apolloClient from "@src/services/apolloClient";

function MyApp({ Component, pageProps }: AppProps) {
  return (
      <ApolloProvider client={apolloClient}>
        <Head>
          <meta name="robots" content="noindex,nofollow,noarchive" />
          <link
            rel="apple-touch-icon"
            sizes="180x180"
            href="/apple-touch-icon.png"
          />
          <link
            rel="icon"
            type="image/png"
            sizes="32x32"
            href="/favicon-32x32.png"
          />
          <link
            rel="icon"
            type="image/png"
            sizes="16x16"
            href="/favicon-16x16.png"
          />
          <link rel="manifest" href="/site.webmanifest" />
        </Head>
        <Header />
        <Component {...pageProps} />
      </ApolloProvider>
  );
}

export default MyApp;

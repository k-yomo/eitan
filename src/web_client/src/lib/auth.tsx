import React from "react";
import { NextPage } from "next";
import { useRouter } from "next/router";
import { routes } from "@src/constants/routes";
import { AccountInfoFragment, ErrorCode, useCurrentAccountQuery } from "@src/generated/graphql"

export interface CurrentAccountProps {
  currentAccount: AccountInfoFragment,
}

function isBrowser() {
  return typeof window !== "undefined";
}

// Require the user to be authenticated in order to render the component.
// If the user isn't authenticated, redirect to the login page.
export function withAuth(WrappedComponent: NextPage<any>) {
  const ComponentWithAuth: NextPage = (props) => {
    if (!isBrowser()) {
      return <></>;
    }
    const router = useRouter();
    const { data, loading, error } = useCurrentAccountQuery();

    if (loading) {
      return <div>Loading...</div>;
    }

    if (error) {
      for (const e of error.graphQLErrors) {
        if (e.extensions!.code === ErrorCode.Unauthenticated) {
          router.push(routes.loginWithOriginalPath(router.asPath));
          return <></>;
        }
        // TODO: Show error page
        return <>Internal server error</>
      }
    }

    const propsWithCurrentUser = {
      ...props,
      currentAccount: data?.currentAccount,
    };
    return <WrappedComponent {...propsWithCurrentUser} />;
  };

  return ComponentWithAuth;
}

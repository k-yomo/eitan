import React from 'react';
import { NextPage } from 'next';
import { useRouter } from 'next/router';
import { routes } from '@src/constants/routes';
import {
  ErrorCode,
  useCurrentUserProfileQuery,
  UserProfileInfoFragment,
} from '@src/generated/graphql';

export interface CurrentUserProfileProps {
  currentUserProfile: UserProfileInfoFragment;
}

function isBrowser() {
  return typeof window !== 'undefined';
}

// Require the user to be authenticated in order to render the component.
// If the user isn't authenticated, redirect to the login page.
export function withAuth(WrappedComponent: NextPage<any>) {
  const ComponentWithAuth: NextPage = (props) => {
    if (!isBrowser()) {
      return <></>;
    }
    const router = useRouter();
    const { data, loading, error } = useCurrentUserProfileQuery();

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
      }
      return <div>Internal server error</div>;
    }

    const propsWithCurrentUser = {
      ...props,
      currentUserProfile: data?.currentUserProfile,
    };
    return <WrappedComponent {...propsWithCurrentUser} />;
  };

  return ComponentWithAuth;
}

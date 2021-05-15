import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> &
  { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> &
  { [SubKey in K]: Maybe<T[SubKey]> };
const defaultOptions = {};
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Time: any;
};

export enum ErrorCode {
  Unauthenticated = 'UNAUTHENTICATED',
  Unauthorized = 'UNAUTHORIZED',
  AlreadyExist = 'ALREADY_EXIST',
  InvalidArgument = 'INVALID_ARGUMENT',
  NotFound = 'NOT_FOUND',
  Internal = 'INTERNAL',
}

export type Node = {
  id: Scalars['ID'];
};

export type Query = {
  node?: Maybe<Node>;
  nodes: Array<Maybe<Node>>;
  currentUserProfile: UserProfile;
};

export type QueryNodeArgs = {
  id: Scalars['ID'];
};

export type QueryNodesArgs = {
  ids: Array<Scalars['ID']>;
};

export enum Role {
  User = 'USER',
}

export type UserProfile = Node & {
  id: Scalars['ID'];
  email: Scalars['String'];
  displayName: Scalars['String'];
  screenImgUrl?: Maybe<Scalars['String']>;
};

export type UserProfilePublic = Node & {
  id: Scalars['ID'];
  displayName: Scalars['String'];
  screenImgUrl?: Maybe<Scalars['String']>;
};

export type UserProfileInfoFragment = Pick<
  UserProfile,
  'id' | 'email' | 'displayName' | 'screenImgUrl'
>;

export type CurrentUserProfileQueryVariables = Exact<{ [key: string]: never }>;

export type CurrentUserProfileQuery = {
  currentUserProfile: UserProfileInfoFragment;
};

export const UserProfileInfoFragmentDoc = gql`
  fragment userProfileInfo on UserProfile {
    id
    email
    displayName
    screenImgUrl
  }
`;
export const CurrentUserProfileDocument = gql`
  query currentUserProfile {
    currentUserProfile {
      ...userProfileInfo
    }
  }
  ${UserProfileInfoFragmentDoc}
`;

/**
 * __useCurrentUserProfileQuery__
 *
 * To run a query within a React component, call `useCurrentUserProfileQuery` and pass it any options that fit your needs.
 * When your component renders, `useCurrentUserProfileQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useCurrentUserProfileQuery({
 *   variables: {
 *   },
 * });
 */
export function useCurrentUserProfileQuery(
  baseOptions?: Apollo.QueryHookOptions<
    CurrentUserProfileQuery,
    CurrentUserProfileQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<
    CurrentUserProfileQuery,
    CurrentUserProfileQueryVariables
  >(CurrentUserProfileDocument, options);
}
export function useCurrentUserProfileLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    CurrentUserProfileQuery,
    CurrentUserProfileQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<
    CurrentUserProfileQuery,
    CurrentUserProfileQueryVariables
  >(CurrentUserProfileDocument, options);
}
export type CurrentUserProfileQueryHookResult = ReturnType<
  typeof useCurrentUserProfileQuery
>;
export type CurrentUserProfileLazyQueryHookResult = ReturnType<
  typeof useCurrentUserProfileLazyQuery
>;
export type CurrentUserProfileQueryResult = Apollo.QueryResult<
  CurrentUserProfileQuery,
  CurrentUserProfileQueryVariables
>;

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
const defaultOptions =  {}
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
  Time: any;
};


export type Account = Node & {
  id: Scalars['ID'];
  email: Scalars['String'];
  displayName: Scalars['String'];
  screenImgUrl?: Maybe<Scalars['String']>;
};

export type AccountPublic = Node & {
  id: Scalars['ID'];
  displayName: Scalars['String'];
  screenImgUrl?: Maybe<Scalars['String']>;
};

export enum ErrorCode {
  Unauthenticated = 'UNAUTHENTICATED',
  Unauthorized = 'UNAUTHORIZED',
  AlreadyExist = 'ALREADY_EXIST',
  InvalidArgument = 'INVALID_ARGUMENT',
  NotFound = 'NOT_FOUND',
  Internal = 'INTERNAL'
}

export type Node = {
  id: Scalars['ID'];
};

export type Query = {
  node?: Maybe<Node>;
  nodes: Array<Maybe<Node>>;
  currentAccount: Account;
};


export type QueryNodeArgs = {
  id: Scalars['ID'];
};


export type QueryNodesArgs = {
  ids: Array<Scalars['ID']>;
};

export enum Role {
  User = 'USER'
}


export type AccountInfoFragment = Pick<Account, 'id' | 'email' | 'displayName' | 'screenImgUrl'>;

export type CurrentAccountQueryVariables = Exact<{ [key: string]: never; }>;


export type CurrentAccountQuery = { currentAccount: AccountInfoFragment };

export const AccountInfoFragmentDoc = gql`
    fragment accountInfo on Account {
  id
  email
  displayName
  screenImgUrl
}
    `;
export const CurrentAccountDocument = gql`
    query currentAccount {
  currentAccount {
    ...accountInfo
  }
}
    ${AccountInfoFragmentDoc}`;

/**
 * __useCurrentAccountQuery__
 *
 * To run a query within a React component, call `useCurrentAccountQuery` and pass it any options that fit your needs.
 * When your component renders, `useCurrentAccountQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useCurrentAccountQuery({
 *   variables: {
 *   },
 * });
 */
export function useCurrentAccountQuery(baseOptions?: Apollo.QueryHookOptions<CurrentAccountQuery, CurrentAccountQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<CurrentAccountQuery, CurrentAccountQueryVariables>(CurrentAccountDocument, options);
      }
export function useCurrentAccountLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<CurrentAccountQuery, CurrentAccountQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<CurrentAccountQuery, CurrentAccountQueryVariables>(CurrentAccountDocument, options);
        }
export type CurrentAccountQueryHookResult = ReturnType<typeof useCurrentAccountQuery>;
export type CurrentAccountLazyQueryHookResult = ReturnType<typeof useCurrentAccountLazyQuery>;
export type CurrentAccountQueryResult = Apollo.QueryResult<CurrentAccountQuery, CurrentAccountQueryVariables>;
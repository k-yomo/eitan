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

export type CurrentUserProfile = {
  id: Scalars['ID'];
  email: Scalars['String'];
  displayName: Scalars['String'];
  screenImgUrl?: Maybe<Scalars['String']>;
};

export enum ErrorCode {
  Unauthenticated = 'UNAUTHENTICATED',
  Unauthorized = 'UNAUTHORIZED',
  AlreadyExist = 'ALREADY_EXIST',
  InvalidArgument = 'INVALID_ARGUMENT',
  NotFound = 'NOT_FOUND',
  Internal = 'INTERNAL',
}

export type FourChoicesQuiz = {
  id: Scalars['ID'];
  quizType: QuizType;
  question: Scalars['String'];
  choices: Array<QuizChoice>;
};

export type FourChoicesQuizAnswer = {
  quiz: Quiz;
  answeredPlayerID: Scalars['ID'];
  correctChoiceID: Scalars['ID'];
};

export type Mutation = {
  updatePlayerId: Player;
  cancelWaitingMatch: Scalars['Boolean'];
  updateCurrentPlayerQuizRoomStatus?: Maybe<Scalars['Boolean']>;
  createAnswer?: Maybe<Scalars['Boolean']>;
};

export type MutationUpdatePlayerIdArgs = {
  input: UpdatePlayerIdInput;
};

export type MutationUpdateCurrentPlayerQuizRoomStatusArgs = {
  roomId: Scalars['ID'];
};

export type MutationCreateAnswerArgs = {
  roomId: Scalars['ID'];
  quizId: Scalars['ID'];
  answer: Scalars['String'];
};

export type Node = {
  id: Scalars['ID'];
};

export type Player = Node & {
  id: Scalars['ID'];
  userId: Scalars['ID'];
  userProfile: UserProfile;
};

export type Query = {
  node?: Maybe<Node>;
  nodes: Array<Maybe<Node>>;
  currentUserProfile: CurrentUserProfile;
  currentPlayer: Player;
};

export type QueryNodeArgs = {
  id: Scalars['ID'];
};

export type QueryNodesArgs = {
  ids: Array<Scalars['ID']>;
};

export type Quiz = FourChoicesQuiz;

export type QuizAnswer = FourChoicesQuizAnswer;

export type QuizChoice = Node & {
  id: Scalars['ID'];
  choice: Scalars['String'];
};

export type QuizRoom = Node & {
  id: Scalars['ID'];
  players: Array<Player>;
};

export enum QuizType {
  FourChoices = 'FourChoices',
}

export enum Role {
  User = 'USER',
}

export type Subscription = {
  randomMatchRoomDecided: QuizRoom;
  quizPosted: Quiz;
  quizAnswered: QuizAnswer;
};

export type SubscriptionQuizPostedArgs = {
  roomId: Scalars['ID'];
};

export type SubscriptionQuizAnsweredArgs = {
  roomId: Scalars['ID'];
};

export type UpdatePlayerIdInput = {
  playerId: Scalars['String'];
};

export type UserProfile = Node & {
  id: Scalars['ID'];
  displayName: Scalars['String'];
  screenImgUrl?: Maybe<Scalars['String']>;
};

export type CurrentUserProfileQueryVariables = Exact<{ [key: string]: never }>;

export type CurrentUserProfileQuery = {
  currentUserProfile: Pick<
    CurrentUserProfile,
    'id' | 'email' | 'displayName' | 'screenImgUrl'
  >;
};

export type IndexPageCurrentPlayerQueryVariables = Exact<{
  [key: string]: never;
}>;

export type IndexPageCurrentPlayerQuery = { currentPlayer: Pick<Player, 'id'> };

export type WaitingRoomPageRandomMatchRoomDecidedSubscriptionVariables = Exact<{
  [key: string]: never;
}>;

export type WaitingRoomPageRandomMatchRoomDecidedSubscription = {
  randomMatchRoomDecided: Pick<QuizRoom, 'id'> & {
    players: Array<
      Pick<Player, 'id'> & {
        userProfile: Pick<UserProfile, 'id' | 'displayName' | 'screenImgUrl'>;
      }
    >;
  };
};

export const CurrentUserProfileDocument = gql`
  query currentUserProfile {
    currentUserProfile {
      id
      email
      displayName
      screenImgUrl
    }
  }
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
export const IndexPageCurrentPlayerDocument = gql`
  query indexPageCurrentPlayer {
    currentPlayer {
      id
    }
  }
`;

/**
 * __useIndexPageCurrentPlayerQuery__
 *
 * To run a query within a React component, call `useIndexPageCurrentPlayerQuery` and pass it any options that fit your needs.
 * When your component renders, `useIndexPageCurrentPlayerQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useIndexPageCurrentPlayerQuery({
 *   variables: {
 *   },
 * });
 */
export function useIndexPageCurrentPlayerQuery(
  baseOptions?: Apollo.QueryHookOptions<
    IndexPageCurrentPlayerQuery,
    IndexPageCurrentPlayerQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<
    IndexPageCurrentPlayerQuery,
    IndexPageCurrentPlayerQueryVariables
  >(IndexPageCurrentPlayerDocument, options);
}
export function useIndexPageCurrentPlayerLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    IndexPageCurrentPlayerQuery,
    IndexPageCurrentPlayerQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<
    IndexPageCurrentPlayerQuery,
    IndexPageCurrentPlayerQueryVariables
  >(IndexPageCurrentPlayerDocument, options);
}
export type IndexPageCurrentPlayerQueryHookResult = ReturnType<
  typeof useIndexPageCurrentPlayerQuery
>;
export type IndexPageCurrentPlayerLazyQueryHookResult = ReturnType<
  typeof useIndexPageCurrentPlayerLazyQuery
>;
export type IndexPageCurrentPlayerQueryResult = Apollo.QueryResult<
  IndexPageCurrentPlayerQuery,
  IndexPageCurrentPlayerQueryVariables
>;
export const WaitingRoomPageRandomMatchRoomDecidedDocument = gql`
  subscription waitingRoomPageRandomMatchRoomDecided {
    randomMatchRoomDecided {
      ... on QuizRoom {
        id
        players {
          id
          userProfile {
            id
            displayName
            screenImgUrl
          }
        }
      }
    }
  }
`;

/**
 * __useWaitingRoomPageRandomMatchRoomDecidedSubscription__
 *
 * To run a query within a React component, call `useWaitingRoomPageRandomMatchRoomDecidedSubscription` and pass it any options that fit your needs.
 * When your component renders, `useWaitingRoomPageRandomMatchRoomDecidedSubscription` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the subscription, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useWaitingRoomPageRandomMatchRoomDecidedSubscription({
 *   variables: {
 *   },
 * });
 */
export function useWaitingRoomPageRandomMatchRoomDecidedSubscription(
  baseOptions?: Apollo.SubscriptionHookOptions<
    WaitingRoomPageRandomMatchRoomDecidedSubscription,
    WaitingRoomPageRandomMatchRoomDecidedSubscriptionVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useSubscription<
    WaitingRoomPageRandomMatchRoomDecidedSubscription,
    WaitingRoomPageRandomMatchRoomDecidedSubscriptionVariables
  >(WaitingRoomPageRandomMatchRoomDecidedDocument, options);
}
export type WaitingRoomPageRandomMatchRoomDecidedSubscriptionHookResult =
  ReturnType<typeof useWaitingRoomPageRandomMatchRoomDecidedSubscription>;
export type WaitingRoomPageRandomMatchRoomDecidedSubscriptionResult =
  Apollo.SubscriptionResult<WaitingRoomPageRandomMatchRoomDecidedSubscription>;

import React from 'react';
import gql from 'graphql-tag';
import { withAuth, WithAuthProps } from '@src/lib/auth';
import { NextPage } from 'next';
import { useWaitingRoomPageRandomMatchRoomDecidedSubscription } from '@src/generated/graphql';

const waitingRoomPageRandomMatchRoomDecided = gql`
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

const WaitingRoom: NextPage<WithAuthProps> = ({
  currentUserProfile,
}: WithAuthProps) => {
  const {
    data,
    loading,
    error,
  } = useWaitingRoomPageRandomMatchRoomDecidedSubscription();

  if (error) {
    // TODO: Show error page
    return <div>Internal server error</div>;
  }

  return (
    <div className="grid place-items-center mx-2 md:my-20 my-10">
      {loading && <div>Matching...</div>}
      {data && (
        <div>
          <div>room id: {data.randomMatchRoomDecided.id}</div>
          {data.randomMatchRoomDecided.players.map((p) => (
            <>
              <div>id: {p.id}</div>
              <div>Name: {p.userProfile.displayName}</div>
            </>
          ))}
        </div>
      )}
    </div>
  );
};

export default withAuth(WaitingRoom);

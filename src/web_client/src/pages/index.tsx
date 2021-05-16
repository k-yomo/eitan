import React from 'react';
import gql from 'graphql-tag';
import {
  WithAuthProps,
  withCurrentUser,
  WithCurrentUserProps,
} from '@src/lib/auth';
import { NextPage } from 'next';
import { useIndexPageCurrentPlayerQuery } from '@src/generated/graphql';
import Link from 'next/link';
import { routes } from '@src/constants/routes';

const indexPageCurrentPlayer = gql`
  query indexPageCurrentPlayer {
    currentPlayer {
      id
    }
  }
`;

const IndexPage: NextPage<WithCurrentUserProps> = ({
  currentUserProfile,
}: WithCurrentUserProps) => {
  return (
    <div className="grid place-items-center mx-2 md:my-20 my-10">
      {currentUserProfile && (
        <ProfileBox currentUserProfile={currentUserProfile} />
      )}
    </div>
  );
};

const ProfileBox = ({ currentUserProfile }: WithAuthProps) => {
  const { data, loading, error } = useIndexPageCurrentPlayerQuery();

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    // TODO: Show error page
    return <div>Internal server error</div>;
  }
  return (
    <>
      <div
        className="w-11/12 p-12 sm:w-8/12 md:w-6/12 lg:w-5/12 2xl:w-4/12
            px-6 py-10 sm:px-10 sm:py-6
            bg-white rounded-lg shadow-md lg:shadow-lg"
      >
        <div className="flex items-center">
          <div>
            <img
              className="inline-block h-9 w-9"
              src={
                currentUserProfile.screenImgUrl ||
                'https://placehold.jp/24/cccccc/cccccc/100x100.png?text=-' // TODO: set appropriate image
              }
              alt={currentUserProfile.displayName}
            />
          </div>
          <div className="ml-3">
            <p className="text-sm font-medium text-gray-700 group-hover:text-gray-900">
              {currentUserProfile.displayName}
            </p>
            <p className="text-xs font-medium text-gray-500 group-hover:text-gray-700">
              @{data?.currentPlayer.id}
            </p>
          </div>
        </div>
        <Link href={routes.waitingRoom()}>
          <a>
            <button className="flex flex-row items-center justify-center w-full space-x-2 my-6 p-3 rounded-sm text-md text-white bg-blue-600 hover:bg-blue-700">
              <div>Random Match</div>
            </button>
          </a>
        </Link>
      </div>
    </>
  );
};

export default withCurrentUser(IndexPage);
